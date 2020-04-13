package keeper

import (
	"fmt"
	"github/irismod/token/exported"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	"github/irismod/token/internal/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	// The supplyKeeper to reduce the supply of the network
	supplyKeeper types.SupplyKeeper

	feeCollectorName string

	// params subspace
	paramSpace params.Subspace
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,
	supplyKeeper types.SupplyKeeper, feeCollectorName string) Keeper {
	// ensure token module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	keeper := Keeper{
		storeKey:         key,
		cdc:              cdc,
		paramSpace:       paramSpace.WithKeyTable(types.ParamKeyTable()),
		supplyKeeper:     supplyKeeper,
		feeCollectorName: feeCollectorName,
	}

	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// IssueToken issues a new token
func (k Keeper) IssueToken(ctx sdk.Context, msg types.MsgIssueToken) error {
	symbol := strings.ToLower(msg.Symbol)
	name := strings.TrimSpace(msg.Name)
	minUnitAlias := strings.ToLower(strings.TrimSpace(msg.MinUnit))

	token := types.NewToken(
		symbol, name, minUnitAlias, msg.Scale, msg.InitialSupply,
		msg.MaxSupply, msg.Mintable, msg.Owner,
	)

	if err := k.AddToken(ctx, token); err != nil {
		return err
	}

	initialSupply := sdk.NewCoin(
		token.MinUnit,
		sdk.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Scale)),
	)

	mintCoins := sdk.NewCoins(initialSupply)

	// Add coins into owner's account
	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		return err
	}

	// sent coins to owner's account
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, token.Owner, mintCoins,
	); err != nil {
		return err
	}

	return nil
}

// EditToken edits the specified token
func (k Keeper) EditToken(ctx sdk.Context, msg types.MsgEditToken) error {
	// get the destination token
	tokenI, err := k.GetToken(ctx, msg.Symbol)
	if err != nil {
		return err
	}

	token := tokenI.(types.Token)

	if !msg.Owner.Equals(token.Owner) {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "the address %d is not the owner of the token %s", msg.Owner, msg.Symbol)
	}

	if msg.MaxSupply > 0 {
		issuedAmt := k.getTokenSupply(ctx, token.MinUnit)
		issuedMainUnitAmt := uint64(issuedAmt.Quo(sdk.NewIntWithDecimal(1, int(token.Scale))).Int64())
		if msg.MaxSupply < issuedMainUnitAmt {
			return sdkerrors.Wrapf(types.ErrInvalidMaxSupply, "max supply must not be less than %d", issuedMainUnitAmt)
		}

		token.MaxSupply = msg.MaxSupply
	}

	if msg.Name != types.DoNotModify {
		token.Name = msg.Name
	}

	if msg.Mintable != types.Nil {
		token.Mintable = msg.Mintable.ToBool()
	}

	if err := k.setToken(ctx, token); err != nil {
		return err
	}

	return nil
}

// TransferTokenOwner transfers the owner of the specified token to a new one
func (k Keeper) TransferTokenOwner(ctx sdk.Context, msg types.MsgTransferTokenOwner) error {
	tokenI, err := k.GetToken(ctx, msg.Symbol)
	if err != nil {
		return err
	}

	token := tokenI.(types.Token)

	if !msg.SrcOwner.Equals(token.Owner) {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "the address %s is not the owner of the token %s", msg.SrcOwner.String(), msg.Symbol)
	}

	token.Owner = msg.DstOwner
	// update token information
	if err := k.setToken(ctx, token); err != nil {
		return err
	}

	// reset all index for query-token
	if err := k.resetStoreKeyForQueryToken(ctx, msg, token); err != nil {
		return err
	}

	return nil
}

// MintToken mints specified amount token to a specified owner
func (k Keeper) MintToken(ctx sdk.Context, msg types.MsgMintToken) error {
	tokenI, err := k.GetToken(ctx, msg.Symbol)
	if err != nil {
		return err
	}

	token := tokenI.(types.Token)

	if !msg.Owner.Equals(token.Owner) {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "the address %s is not the owner of the token %s", msg.Owner.String(), msg.Symbol)
	}

	if !token.Mintable {
		return sdkerrors.Wrapf(types.ErrNotMintable, "the token %s is set to be non-mintable", msg.Symbol)
	}

	issuedAmt := k.getTokenSupply(ctx, token.MinUnit)
	mintableMaxAmt := sdk.NewIntWithDecimal(int64(token.MaxSupply), int(token.Scale)).Sub(issuedAmt)
	mintableMaxMainUnitAmt := uint64(mintableMaxAmt.Quo(sdk.NewIntWithDecimal(1, int(token.Scale))).Int64())

	if msg.Amount > mintableMaxMainUnitAmt {
		return sdkerrors.Wrapf(types.ErrInvalidMaxSupply, "The amount of minting tokens plus the total amount of issued tokens has exceeded the maximum supply, only accepts amount (0, %d]", mintableMaxMainUnitAmt)
	}

	mintCoin := sdk.NewCoin(token.MinUnit, sdk.NewIntWithDecimal(int64(msg.Amount), int(token.Scale)))
	mintCoins := sdk.NewCoins(mintCoin)

	// mint coins
	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		return err
	}

	mintAcc := msg.To
	if mintAcc.Empty() {
		mintAcc = token.Owner
	}

	// sent coins to owner's account
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mintAcc, mintCoins); err != nil {
		return err
	}

	return nil
}

// GetTokens returns all existing tokens
func (k Keeper) GetTokens(ctx sdk.Context, owner sdk.AccAddress) (tokens []exported.TokenI) {
	store := ctx.KVStore(k.storeKey)

	var it sdk.Iterator
	if owner == nil {
		it = sdk.KVStorePrefixIterator(store, types.PrefixTokenForSymbol)
		defer it.Close()

		for ; it.Valid(); it.Next() {
			var token types.Token
			k.cdc.MustUnmarshalBinaryLengthPrefixed(it.Value(), &token)

			tokens = append(tokens, token)
		}
		return
	}

	it = sdk.KVStorePrefixIterator(store, types.KeyTokens(owner, ""))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var symbol string
		k.cdc.MustUnmarshalBinaryLengthPrefixed(it.Value(), &symbol)

		token, err := k.GetToken(ctx, symbol)
		if err != nil {
			continue
		}
		tokens = append(tokens, token)
	}
	return
}

// GetToken returns the token of the specified symbol or minUint
func (k Keeper) GetToken(ctx sdk.Context, denom string) (token exported.TokenI, err error) {
	store := ctx.KVStore(k.storeKey)

	if token, err := k.getToken(ctx, denom); err == nil {
		return token, nil
	}

	bz := store.Get(types.KeyMinUint(denom))

	var symbol string
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &symbol)
	return k.getToken(ctx, symbol)
}

// AddToken saves a new token
func (k Keeper) AddToken(ctx sdk.Context, token types.Token) error {
	if k.HasToken(ctx, token.Symbol) {
		return sdkerrors.Wrapf(types.ErrSymbolAlreadyExists, "symbol already exists: %s", token.Symbol)
	}

	if k.HasToken(ctx, token.MinUnit) {
		return sdkerrors.Wrapf(types.ErrMinUnitAlreadyExists, "min-unit already exists: %s", token.MinUnit)
	}

	// set token
	if err := k.setToken(ctx, token); err != nil {
		return err
	}

	// Set token to be prefixed with owner
	if err := k.setWithOwner(ctx, token.Owner, token.Symbol); err != nil {
		return err
	}

	// Set token to be prefixed with min_unit
	if err := k.setWithMinUnit(ctx, token.MinUnit, token.Symbol); err != nil {
		return err
	}

	return nil
}

// HasToken asserts a token exists
func (k Keeper) HasToken(ctx sdk.Context, param string) bool {
	store := ctx.KVStore(k.storeKey)
	existed := store.Has(types.KeySymbol(param))
	if existed {
		return existed
	}

	return store.Has(types.KeyMinUint(param))
}

// GetParamSet returns asset params from the global param store
func (k Keeper) GetParamSet(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

// SetParamSet set asset params from the global param store
func (k Keeper) SetParamSet(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) setWithOwner(ctx sdk.Context, owner sdk.AccAddress, symbol string) error {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(symbol)

	store.Set(types.KeyTokens(owner, symbol), bz)
	return nil
}

func (k Keeper) setWithMinUnit(ctx sdk.Context, minUnit, symbol string) error {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(symbol)

	store.Set(types.KeyMinUint(minUnit), bz)
	return nil
}

func (k Keeper) setToken(ctx sdk.Context, token types.Token) error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(token)

	store.Set(types.KeySymbol(token.Symbol), bz)
	return nil
}

func (k Keeper) getToken(ctx sdk.Context, symbol string) (token types.Token, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeySymbol(symbol))
	if bz == nil {
		return token, sdkerrors.Wrap(types.ErrTokenNotExists, fmt.Sprintf("token %s does not exist", symbol))
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &token)
	return token, nil
}

// reset all index by DstOwner of token for query-token command
func (k Keeper) resetStoreKeyForQueryToken(ctx sdk.Context, msg types.MsgTransferTokenOwner, token types.Token) error {
	store := ctx.KVStore(k.storeKey)

	// delete the old key
	store.Delete(types.KeyTokens(msg.SrcOwner, token.Symbol))

	// add the new key
	return k.setWithOwner(ctx, msg.DstOwner, token.Symbol)
}

// getTokenSupply query issued tokens supply from the total supply
func (k Keeper) getTokenSupply(ctx sdk.Context, denom string) sdk.Int {
	return k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)
}

package keeper

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	"github/irismod/asset/internal/types"
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
	// ensure asset module account is set
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

	token := types.NewFungibleToken(
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
	token, err := k.getToken(ctx, msg.Symbol)
	if err != nil {
		return err
	}

	if !msg.Owner.Equals(token.Owner) {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "the address %d is not the owner of the token %s", msg.Owner, msg.Symbol)
	}

	if msg.MaxSupply > 0 {
		issuedAmt := k.getTokenSupply(ctx, token.MinUnit)
		issuedMainUnitAmt := uint64(issuedAmt.Quo(sdk.NewIntWithDecimal(1, int(token.Scale))).Int64())
		if msg.MaxSupply < issuedMainUnitAmt {
			return sdkerrors.Wrapf(types.ErrInvalidAssetMaxSupply, "max supply must not be less than %d", issuedMainUnitAmt)
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
	// get the destination token
	token, err := k.getToken(ctx, msg.Symbol)
	if err != nil {
		return err
	}

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
	token, err := k.getToken(ctx, msg.Symbol)
	if err != nil {
		return err
	}

	if !msg.Owner.Equals(token.Owner) {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "the address %s is not the owner of the token %s", msg.Owner.String(), msg.Symbol)
	}

	if !token.Mintable {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "the token %s is set to be non-mintable", msg.Symbol)
	}

	issuedAmt := k.getTokenSupply(ctx, token.MinUnit)
	mintableMaxAmt := sdk.NewIntWithDecimal(int64(token.MaxSupply), int(token.Scale)).Sub(issuedAmt)
	mintableMaxMainUnitAmt := uint64(mintableMaxAmt.Quo(sdk.NewIntWithDecimal(1, int(token.Scale))).Int64())

	if msg.Amount > mintableMaxMainUnitAmt {
		return sdkerrors.Wrapf(types.ErrInvalidAssetMaxSupply, "The amount of minting tokens plus the total amount of issued tokens has exceeded the maximum supply, only accepts amount (0, %d]", mintableMaxMainUnitAmt)
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

// IterateTokens iterates through all existing tokens
func (k Keeper) IterateTokens(ctx sdk.Context, op func(token types.FungibleToken) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.PrefixToken)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var token types.FungibleToken
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &token)

		if stop := op(token); stop {
			break
		}
	}
}

// GetAllTokens returns all existing tokens
func (k Keeper) GetAllTokens(ctx sdk.Context) (tokens []types.FungibleToken) {
	k.IterateTokens(ctx, func(token types.FungibleToken) (stop bool) {
		tokens = append(tokens, token)
		return false
	})

	return
}

// AddToken saves a new token
func (k Keeper) AddToken(ctx sdk.Context, token types.FungibleToken) error {
	if k.HasToken(ctx, token.Symbol) {
		return sdkerrors.Wrapf(types.ErrAssetAlreadyExists, "token already exists: %s", token.Symbol)
	}

	// set token
	if err := k.setToken(ctx, token); err != nil {
		return err
	}

	// Set token to be prefixed with owner
	if err := k.setOwnerToken(ctx, token.Owner, token); err != nil {
		return err
	}

	return nil
}

// HasToken asserts a token exists
func (k Keeper) HasToken(ctx sdk.Context, symbol string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyToken(symbol))
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

func (k Keeper) iterateTokensWithOwner(ctx sdk.Context, owner sdk.AccAddress, op func(token types.FungibleToken) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyTokens(owner, ""))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var symbol string
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &symbol)

		token, err := k.getToken(ctx, symbol)
		if err != nil {
			continue
		}

		if stop := op(token); stop {
			break
		}
	}
}

func (k Keeper) setOwnerToken(ctx sdk.Context, owner sdk.AccAddress, token types.FungibleToken) error {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(token.Symbol)

	store.Set(types.KeyTokens(owner, token.Symbol), bz)
	return nil
}

func (k Keeper) setToken(ctx sdk.Context, token types.FungibleToken) error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(token)

	store.Set(types.KeyToken(token.Symbol), bz)
	return nil
}

// reset all index by DstOwner of token for query-token command
func (k Keeper) resetStoreKeyForQueryToken(ctx sdk.Context, msg types.MsgTransferTokenOwner, token types.FungibleToken) error {
	store := ctx.KVStore(k.storeKey)

	// delete the old key
	store.Delete(types.KeyTokens(msg.SrcOwner, token.Symbol))

	// add the new key
	return k.setOwnerToken(ctx, msg.DstOwner, token)
}

func (k Keeper) getToken(ctx sdk.Context, symbol string) (token types.FungibleToken, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyToken(symbol))
	if bz == nil {
		return token, sdkerrors.Wrap(types.ErrAssetNotExists, fmt.Sprintf("token %s does not exist", symbol))
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &token)
	return token, nil
}

// GetToken wraps getToken for export
func (k Keeper) GetToken(ctx sdk.Context, symbol string) (types.FungibleToken, error) {
	return k.getToken(ctx, symbol)
}

// getTokenSupply query issued tokens supply from the total supply
func (k Keeper) getTokenSupply(ctx sdk.Context, denom string) sdk.Int {
	return k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)
}

// addCollectedFees implements an alias call to the underlying supply keeper's
// addCollectedFees to be used in BeginBlocker.
func (k Keeper) addCollectedFees(ctx sdk.Context, fees sdk.Coins) error {
	return k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}

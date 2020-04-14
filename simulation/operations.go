package simulation

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/irismod/token/exported"
	"github.com/irismod/token/keeper"
	"github.com/irismod/token/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Simulation operation weights constants
const (
	OpWeightMsgIssueToken         = "op_weight_msg_issue_token"
	OpWeightMsgEditToken          = "op_weight_msg_edit_token"
	OpWeightMsgMintToken          = "op_weight_msg_mint_token"
	OpWeightMsgTransferTokenOwner = "op_weight_msg_transfer_token_owner"
)

var (
	nativeToken = types.GetNativeToken()
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simulation.AppParams,
	cdc *codec.Codec,
	k keeper.Keeper,
	ak auth.AccountKeeper) simulation.WeightedOperations {

	var weightIssue, weightEdit, weightMint, weightTransfer int
	appParams.GetOrGenerate(cdc, OpWeightMsgIssueToken, &weightIssue, nil,
		func(_ *rand.Rand) {
			weightIssue = 100
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgEditToken, &weightEdit, nil,
		func(_ *rand.Rand) {
			weightEdit = 50
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgMintToken, &weightMint, nil,
		func(_ *rand.Rand) {
			weightMint = 50
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgTransferTokenOwner, &weightTransfer, nil,
		func(_ *rand.Rand) {
			weightTransfer = 50
		},
	)

	return simulation.WeightedOperations{
		//simulation.NewWeightedOperation(
		//	weightIssue,
		//	SimulateIssueToken(k, ak),
		//),
		simulation.NewWeightedOperation(
			weightEdit,
			SimulateEditToken(k, ak),
		),
		simulation.NewWeightedOperation(
			weightMint,
			SimulateMintToken(k, ak),
		),
		simulation.NewWeightedOperation(
			weightTransfer,
			SimulateTransferTokenOwner(k, ak),
		),
	}
}

// SimulateIssueToken tests and runs a single msg issue a new token
func SimulateIssueToken(k keeper.Keeper, ak auth.AccountKeeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		token, maxFees := genToken(ctx, r, ak, k, accs)
		msg := types.NewMsgIssueToken(token.Symbol, token.MinUnit, token.Name, token.Scale, token.InitialSupply, token.MaxSupply, token.Mintable, token.Owner)

		simAccount, found := simulation.FindAccount(accs, token.Owner)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("account %s not found", token.Owner)
		}
		account := ak.GetAccount(ctx, msg.Owner)
		fees, err := simulation.RandomFees(r, ctx, maxFees)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if _, _, err = app.Deliver(tx); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, "simulate issue token"), nil, nil
	}
}

// SimulateEditToken tests and runs a single msg edit a existed token
func SimulateEditToken(k keeper.Keeper, ak auth.AccountKeeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		token, _ := selectOneToken(ctx, k, ak, false)
		simAccount, found := simulation.FindAccount(accs, token.GetOwner())
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("account %s not found", token.GetOwner())
		}

		msg := types.NewMsgEditToken(token.GetName(), token.GetSymbol(), token.GetMaxSupply(), types.True, token.GetOwner())

		account := ak.GetAccount(ctx, msg.Owner)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if _, _, err = app.Deliver(tx); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, "simulate edit token"), nil, nil
	}
}

// SimulateMintToken tests and runs a single msg mint a existed token
func SimulateMintToken(k keeper.Keeper, ak auth.AccountKeeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		token, maxFee := selectOneToken(ctx, k, ak, true)
		ownerAccount, found := simulation.FindAccount(accs, token.GetOwner())
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("account %s not found", token.GetOwner())
		}

		simToAccount, _ := simulation.RandomAcc(r, accs)

		msg := types.NewMsgMintToken(token.GetSymbol(), token.GetOwner(), simToAccount.Address, 100)

		account := ak.GetAccount(ctx, msg.Owner)
		fees, err := simulation.RandomFees(r, ctx, maxFee)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			ownerAccount.PrivKey,
		)

		if _, _, err = app.Deliver(tx); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, "simulate mint token"), nil, nil
	}
}

// SimulateTransferTokenOwner tests and runs a single msg transfer to others
func SimulateTransferTokenOwner(k keeper.Keeper, ak auth.AccountKeeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		token, _ := selectOneToken(ctx, k, ak, false)
		simAccount, found := simulation.FindAccount(accs, token.GetOwner())
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("account %s not found", token.GetOwner())
		}

		var simToAccount, _ = simulation.RandomAcc(r, accs)
		for simToAccount.Address.Equals(token.GetOwner()) {
			simToAccount, _ = simulation.RandomAcc(r, accs)
		}

		msg := types.NewMsgTransferTokenOwner(token.GetOwner(), simToAccount.Address, token.GetSymbol())

		account := ak.GetAccount(ctx, msg.SrcOwner)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if _, _, err = app.Deliver(tx); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, "simulate transfer token"), nil, nil
	}
}

func selectOneToken(ctx sdk.Context,
	k keeper.Keeper,
	ak auth.AccountKeeper,
	mint bool) (token exported.TokenI, maxFees sdk.Coins) {
	tokens := k.GetTokens(ctx, nil)
	if len(tokens) == 0 {
		panic("No token available")
	}

	for _, t := range tokens {
		if t.GetSymbol() == types.GetNativeToken().Symbol {
			continue
		}
		if !mint {
			return t, nil
		}

		mintFee := k.GetTokenMintFee(ctx, t.GetSymbol())
		account := ak.GetAccount(ctx, t.GetOwner())
		spendable := account.SpendableCoins(ctx.BlockTime())
		spendableStake := spendable.AmountOf(nativeToken.MinUnit)
		if spendableStake.IsZero() || spendableStake.LT(mintFee.Amount) {
			continue
		}
		maxFees = sdk.NewCoins(sdk.NewCoin(nativeToken.MinUnit, spendableStake).Sub(mintFee))
		token = t
		return
	}
	panic("No token mintable")
}

func randStringBetween(r *rand.Rand, min, max int) string {
	strLen := simulation.RandIntBetween(r, min, max)
	randStr := simulation.RandStringOfLength(r, strLen)
	return randStr
}

func genToken(ctx sdk.Context,
	r *rand.Rand,
	ak auth.AccountKeeper,
	k keeper.Keeper,
	accs []simulation.Account,
) (types.Token, sdk.Coins) {

	var token types.Token
	token = randToken(r, accs)

	for k.HasToken(ctx, token.Symbol) {
		token = randToken(r, accs)
	}

	issueFee := k.GetTokenIssueFee(ctx, token.Symbol)

	account, maxFees := filterAccount(ctx, r, ak, accs, issueFee)
	token.Owner = account

	return token, maxFees
}

func filterAccount(ctx sdk.Context,
	r *rand.Rand,
	ak auth.AccountKeeper,
	accs []simulation.Account, fee sdk.Coin) (owner sdk.AccAddress, maxFees sdk.Coins) {
loop:
	simAccount, _ := simulation.RandomAcc(r, accs)
	account := ak.GetAccount(ctx, simAccount.Address)
	spendable := account.SpendableCoins(ctx.BlockTime())
	spendableStake := spendable.AmountOf(nativeToken.MinUnit)
	if spendableStake.IsZero() || spendableStake.LT(fee.Amount) {
		goto loop
	}
	owner = account.GetAddress()
	maxFees = sdk.NewCoins(sdk.NewCoin(nativeToken.MinUnit, spendableStake).Sub(fee))
	return
}

func randToken(r *rand.Rand,
	accs []simulation.Account,
) types.Token {

	symbol := randStringBetween(r, types.MinimumSymbolLen, types.MaximumSymbolLen)
	minUint := randStringBetween(r, types.MinimumMinUnitLen, types.MaximumMinUnitLen)
	name := randStringBetween(r, 1, types.MaximumNameLen)
	scale := simulation.RandIntBetween(r, 1, int(types.MaximumScale))
	initialSupply := r.Int63n(int64(types.MaximumInitSupply))
	maxSupply := r.Int63n(int64(types.MaximumMaxSupply-types.MaximumInitSupply)) + initialSupply
	simAccount, _ := simulation.RandomAcc(r, accs)

	return types.Token{
		Symbol:        strings.ToLower(symbol),
		Name:          name,
		Scale:         uint8(scale),
		MinUnit:       strings.ToLower(minUint),
		InitialSupply: uint64(initialSupply),
		MaxSupply:     uint64(maxSupply),
		Mintable:      true,
		Owner:         simAccount.Address,
	}
}

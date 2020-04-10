package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github/irismod/asset/internal/keeper"
	"github/irismod/asset/internal/types"

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
		simulation.NewWeightedOperation(
			weightIssue,
			SimulateIssueToken(k, ak),
		),
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

		token, maxFees := randomToken(ctx, r, ak, k, accs)
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

		_, _, err := SimulateIssueToken(k, ak)(r, app, ctx, accs, chainID)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		token := selectOneToken(ctx, k, false)
		simAccount, found := simulation.FindAccount(accs, token.Owner)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("account %s not found", token.Owner)
		}

		msg := types.NewMsgEditToken(token.Name, token.Symbol, token.MaxSupply, types.True, token.Owner)

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

		_, _, err := SimulateIssueToken(k, ak)(r, app, ctx, accs, chainID)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		token := selectOneToken(ctx, k, true)
		ownerAccount, found := simulation.FindAccount(accs, token.Owner)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("account %s not found", token.Owner)
		}

		simToAccount, _ := simulation.RandomAcc(r, accs)

		msg := types.NewMsgMintToken(token.Symbol, token.Owner, simToAccount.Address, 100)

		//mintFee := k.GetTokenMintFee(ctx, token.Symbol)

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

		_, _, err := SimulateIssueToken(k, ak)(r, app, ctx, accs, chainID)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		token := selectOneToken(ctx, k, false)
		simAccount, found := simulation.FindAccount(accs, token.Owner)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("account %s not found", token.Owner)
		}

		simToAccount, _ := simulation.RandomAcc(r, accs)
		msg := types.NewMsgTransferTokenOwner(token.Owner, simToAccount.Address, token.Symbol)

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

func selectOneToken(ctx sdk.Context, k keeper.Keeper, mintable bool) types.FungibleToken {
	tokens := k.GetTokens(ctx, nil)
	if len(tokens) == 0 {
		panic("No token available")
	}

	if !mintable {
	loop:
		idx := rand.Intn(len(tokens))
		token := tokens[idx]
		if token.Symbol == types.GetNativeToken().Symbol {
			goto loop
		}
		return token
	}

	for _, token := range tokens {
		if token.Mintable {
			return token
		}
	}

	for _, token := range tokens {
		if token.Symbol == types.GetNativeToken().Symbol {
			continue
		}
		if !mintable {
			return token
		}
		if token.Mintable {
			return token
		}
	}
	panic("No token mintable")
}

func randStringBetween(r *rand.Rand, min, max int) string {
	strLen := simulation.RandIntBetween(r, min, max)
	randStr := simulation.RandStringOfLength(r, strLen)
	return randStr
}

func randomToken(ctx sdk.Context,
	r *rand.Rand,
	ak auth.AccountKeeper,
	k keeper.Keeper,
	accs []simulation.Account,
) (types.FungibleToken, sdk.Coins) {

	symbol := randStringBetween(r, types.MinimumAssetSymbolLen, types.MaximumAssetSymbolLen)
	minUint := randStringBetween(r, types.MinimumAssetMinUnitLen, types.MaximumAssetMinUnitLen)
	name := randStringBetween(r, 1, types.MaximumAssetNameLen)
	scale := simulation.RandIntBetween(r, 1, int(types.MaximumAssetDecimal))
	initialSupply := r.Int63n(int64(types.MaximumAssetInitSupply))
	maxSupply := r.Int63n(int64(types.MaximumAssetMaxSupply-types.MaximumAssetInitSupply)) + initialSupply

	issueFee := k.GetTokenIssueFee(ctx, symbol)

	exit := make(chan int, 1)
	var owner sdk.AccAddress
	var maxFees sdk.Coins

	go func() {
	loop:
		simAccount, _ := simulation.RandomAcc(r, accs)
		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := account.SpendableCoins(ctx.BlockTime())
		spendableStake := spendable.AmountOf(nativeToken.MinUnit)
		if spendableStake.IsZero() || spendableStake.LT(issueFee.Amount) {
			goto loop
		}
		owner = account.GetAddress()
		maxFees = sdk.NewCoins(sdk.NewCoin(nativeToken.MinUnit, spendableStake).Sub(issueFee))
		exit <- 1
	}()

	select {
	case <-exit:
	case <-time.After(30 * time.Second):
		panic("no Spendable coins")
	}

	return types.FungibleToken{
		Symbol:        symbol,
		Name:          name,
		Scale:         uint8(scale),
		MinUnit:       minUint,
		InitialSupply: uint64(initialSupply),
		MaxSupply:     uint64(maxSupply),
		Mintable:      true,
		Owner:         owner,
	}, maxFees
}

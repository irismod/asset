package token_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	"github/irismod/token"
	simapp "github/irismod/token/app"
	"github/irismod/token/types"
)

const (
	isCheck = false
)

var (
	nativeToken = types.GetNativeToken()
	denom       = nativeToken.Symbol
	owner       = sdk.AccAddress([]byte("tokenTest"))
	initAmt     = sdk.NewIntWithDecimal(100000000, int(6))
	initCoin    = sdk.Coins{sdk.NewCoin(denom, initAmt)}
)

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

type HandlerSuite struct {
	suite.Suite

	cdc    *codec.Codec
	ctx    sdk.Context
	keeper token.Keeper
	sk     supply.Keeper
	bk     bank.Keeper
}

func (suite *HandlerSuite) SetupTest() {
	app := simapp.Setup(isCheck)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(isCheck, abci.Header{})
	suite.keeper = app.TokenKeeper
	suite.bk = app.BankKeeper
	suite.sk = app.SupplyKeeper

	// set params
	suite.keeper.SetParamSet(suite.ctx, types.DefaultParams())
	// init tokens to addr
	err := suite.sk.MintCoins(suite.ctx, types.ModuleName, initCoin)
	suite.NoError(err)
	err = suite.sk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, owner, initCoin)
	suite.NoError(err)
}

func (suite *HandlerSuite) TestIssueToken() {
	h := token.NewHandler(suite.keeper)

	balance := suite.bk.GetCoins(suite.ctx, owner)
	nativeTokenAmt1 := balance.AmountOf(denom)

	msg := types.NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 21000000, 21000000, false, owner)

	_, err := h(suite.ctx, msg)
	suite.NoError(err)

	balance = suite.bk.GetCoins(suite.ctx, owner)
	nativeTokenAmt2 := balance.AmountOf(denom)

	fee := suite.keeper.GetTokenIssueFee(suite.ctx, msg.Symbol)

	suite.Equal(nativeTokenAmt1.Sub(fee.Amount), nativeTokenAmt2)

	mintTokenAmt := sdk.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Scale))
	suite.Equal(balance.AmountOf(msg.MinUnit), mintTokenAmt)
}

func (suite *HandlerSuite) TestMintToken() {
	msg := types.NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 1000, 2000, true, owner)

	err := suite.keeper.IssueToken(suite.ctx, msg)
	suite.NoError(err)

	suite.True(suite.keeper.HasToken(suite.ctx, msg.Symbol))

	balance := suite.bk.GetCoins(suite.ctx, owner)
	beginBtcAmt := balance.AmountOf(msg.MinUnit)
	suite.Equal(sdk.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Scale)), beginBtcAmt)

	beginNativeAmt := balance.AmountOf(denom)

	h := token.NewHandler(suite.keeper)

	msgMintToken := types.NewMsgMintToken(msg.Symbol, owner, nil, 1000)
	_, err = h(suite.ctx, msgMintToken)
	suite.NoError(err)

	balance = suite.bk.GetCoins(suite.ctx, owner)
	endBtcAmt := balance.AmountOf(msg.MinUnit)

	mintBtcAmt := sdk.NewIntWithDecimal(int64(msgMintToken.Amount), int(msg.Scale))
	suite.Equal(beginBtcAmt.Add(mintBtcAmt), endBtcAmt)

	fee := suite.keeper.GetTokenMintFee(suite.ctx, msg.Symbol)
	endNativeAmt := balance.AmountOf(denom)
	suite.Equal(beginNativeAmt.Sub(fee.Amount), endNativeAmt)
}

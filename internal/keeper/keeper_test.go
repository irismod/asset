package keeper_test

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	simapp "github/irismod/token/app"
	"github/irismod/token/internal/keeper"
	"github/irismod/token/internal/types"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	isCheck = false
)

var (
	denom    = types.GetNativeToken().Symbol
	owner    = sdk.AccAddress([]byte("tokenTest"))
	initAmt  = sdk.NewIntWithDecimal(100000000, int(6))
	initCoin = sdk.Coins{sdk.NewCoin(denom, initAmt)}
)

type KeeperSuite struct {
	suite.Suite

	cdc    *codec.Codec
	ctx    sdk.Context
	keeper keeper.Keeper
	sk     supply.Keeper
	bk     bank.Keeper
}

func (suite *KeeperSuite) SetupTest() {

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

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperSuite))
}

func (suite *KeeperSuite) TestIssueToken() {
	msg := types.NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 21000000, 21000000, false, owner)

	err := suite.keeper.IssueToken(suite.ctx, msg)
	require.NoError(suite.T(), err)

	suite.True(suite.keeper.HasToken(suite.ctx, msg.Symbol))

	token, err := suite.keeper.GetToken(suite.ctx, msg.Symbol)
	require.NoError(suite.T(), err)

	suite.Equal(msg.MinUnit, token.GetMinUnit())
	suite.Equal(msg.Owner, token.GetOwner())

	ftJson, _ := json.Marshal(msg)
	tokenJson, _ := json.Marshal(token)
	suite.Equal(ftJson, tokenJson)
}

func (suite *KeeperSuite) TestEditToken() {

	suite.TestIssueToken()

	mintable := types.True
	msgEditToken := types.NewMsgEditToken("Bitcoin Token", "btc", 22000000, mintable, owner)
	err := suite.keeper.EditToken(suite.ctx, msgEditToken)
	require.NoError(suite.T(), err)

	token2, err := suite.keeper.GetToken(suite.ctx, msgEditToken.Symbol)
	require.NoError(suite.T(), err)

	expToken := types.NewToken("btc", "Bitcoin Token", "satoshi", 18, 21000000, 22000000, mintable.ToBool(), owner)

	expJson, _ := json.Marshal(expToken)
	actJson, _ := json.Marshal(token2)
	suite.Equal(expJson, actJson)

}

func (suite *KeeperSuite) TestMintToken() {

	msg := types.NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 1000, 2000, true, owner)

	err := suite.keeper.IssueToken(suite.ctx, msg)
	require.NoError(suite.T(), err)

	suite.True(suite.keeper.HasToken(suite.ctx, msg.Symbol))

	balance := suite.bk.GetCoins(suite.ctx, owner)
	amt := balance.AmountOf(msg.MinUnit)
	suite.Equal("1000000000000000000000", amt.String())

	msgMintToken := types.NewMsgMintToken(msg.Symbol, owner, nil, 1000)
	err = suite.keeper.MintToken(suite.ctx, msgMintToken)
	require.NoError(suite.T(), err)

	balance = suite.bk.GetCoins(suite.ctx, owner)
	amt = balance.AmountOf(msg.MinUnit)
	suite.Equal("2000000000000000000000", amt.String())
}

func (suite *KeeperSuite) TestTransferToken() {

	suite.TestIssueToken()

	dstOwner := sdk.AccAddress([]byte("TokenDstOwner"))
	msg := types.MsgTransferTokenOwner{
		SrcOwner: owner,
		DstOwner: dstOwner,
		Symbol:   "btc",
	}
	err := suite.keeper.TransferTokenOwner(suite.ctx, msg)
	require.NoError(suite.T(), err)

	token, err := suite.keeper.GetToken(suite.ctx, "btc")
	require.NoError(suite.T(), err)
	suite.Equal(dstOwner, token.GetOwner())
}

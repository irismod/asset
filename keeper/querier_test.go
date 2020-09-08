package keeper_test

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/irismod/token/keeper"
	"github.com/irismod/token/types"
)

func (suite *KeeperTestSuite) TestQueryToken() {
	ctx := suite.ctx
	querier := keeper.NewQuerier(suite.keeper, suite.cdc)

	params := types.QueryTokenParams{
		Denom: types.GetNativeToken().Symbol,
	}
	bz := suite.cdc.MustMarshalJSON(params)
	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryToken),
		Data: bz,
	}

	data, err := querier(ctx, []string{types.QueryToken}, query)
	suite.Nil(err)

	data2 := codec.MustMarshalJSONIndent(suite.cdc, types.GetNativeToken())
	suite.Equal(data2, data)

	//query by mint_unit
	params = types.QueryTokenParams{
		Denom: types.GetNativeToken().MinUnit,
	}

	bz = suite.cdc.MustMarshalJSON(params)
	query = abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryToken),
		Data: bz,
	}

	data, err = querier(ctx, []string{types.QueryToken}, query)
	suite.Nil(err)

	data2 = codec.MustMarshalJSONIndent(suite.cdc, types.GetNativeToken())
	suite.Equal(data2, data)
}

func (suite *KeeperTestSuite) TestQueryTokens() {
	ctx := suite.ctx
	querier := keeper.NewQuerier(suite.keeper, suite.cdc)

	params := types.QueryTokensParams{
		Owner: nil,
	}
	bz := suite.cdc.MustMarshalJSON(params)
	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryTokens),
		Data: bz,
	}

	data, err := querier(ctx, []string{types.QueryTokens}, query)
	suite.Nil(err)

	data2 := codec.MustMarshalJSONIndent(suite.cdc, []types.TokenI{types.GetNativeToken()})
	suite.Equal(data2, data)
}

func (suite *KeeperTestSuite) TestQueryFees() {
	ctx := suite.ctx
	querier := keeper.NewQuerier(suite.keeper, suite.cdc)

	params := types.QueryTokenFeesParams{
		Symbol: "btc",
	}
	bz := suite.cdc.MustMarshalJSON(params)
	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryFees),
		Data: bz,
	}

	data, err := querier(ctx, []string{types.QueryFees}, query)
	suite.Nil(err)

	var fee types.QueryFeesResponse
	suite.cdc.MustUnmarshalJSON(data, &fee)
	suite.Equal(false, fee.Exist)
	suite.Equal(fmt.Sprintf("60000%s", types.GetNativeToken().MinUnit), fee.IssueFee.String())
	suite.Equal(fmt.Sprintf("6000%s", types.GetNativeToken().MinUnit), fee.MintFee.String())
}

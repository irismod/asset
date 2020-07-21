package keeper_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	simapp "github.com/irismod/token/app"
	"github.com/irismod/token/keeper"
	"github.com/irismod/token/types"
)

func TestQueryToken(t *testing.T) {
	app := simapp.Setup(isCheck)
	ctx := app.BaseApp.NewContext(isCheck, abci.Header{})
	querier := keeper.NewQuerier(app.TokenKeeper)

	params := types.QueryTokenParams{
		Denom: types.GetNativeToken().Symbol,
	}
	bz := app.Codec().MustMarshalJSON(params)
	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryToken),
		Data: bz,
	}

	data, err := querier(ctx, []string{types.QueryToken}, query)
	require.Nil(t, err)

	data2 := codec.MustMarshalJSONIndent(app.Codec(), types.GetNativeToken())
	require.EqualValues(t, data2, data)

	//query by mint_unit
	params = types.QueryTokenParams{
		Denom: types.GetNativeToken().MinUnit,
	}

	bz = app.Codec().MustMarshalJSON(params)
	query = abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryToken),
		Data: bz,
	}

	data, err = querier(ctx, []string{types.QueryToken}, query)
	require.Nil(t, err)

	data2 = codec.MustMarshalJSONIndent(app.Codec(), types.GetNativeToken())
	require.EqualValues(t, data2, data)

}

func TestQueryTokens(t *testing.T) {
	app := simapp.Setup(isCheck)
	ctx := app.BaseApp.NewContext(isCheck, abci.Header{})
	querier := keeper.NewQuerier(app.TokenKeeper)

	params := types.QueryTokensParams{
		Owner: nil,
	}
	bz := app.Codec().MustMarshalJSON(params)
	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryTokens),
		Data: bz,
	}

	data, err := querier(ctx, []string{types.QueryTokens}, query)
	require.Nil(t, err)

	data2 := codec.MustMarshalJSONIndent(app.Codec(), []types.TokenI{types.GetNativeToken()})
	require.EqualValues(t, data2, data)
}

func TestQueryFees(t *testing.T) {
	app := simapp.Setup(isCheck)
	ctx := app.BaseApp.NewContext(isCheck, abci.Header{})
	querier := keeper.NewQuerier(app.TokenKeeper)

	params := types.QueryTokenFeesParams{
		Symbol: "btc",
	}
	bz := app.Codec().MustMarshalJSON(params)
	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryFees),
		Data: bz,
	}

	data, err := querier(ctx, []string{types.QueryFees}, query)
	require.Nil(t, err)

	var fee types.QueryFeesResponse
	app.Codec().MustUnmarshalJSON(data, &fee)
	require.Equal(t, false, fee.Exist)
	require.Equal(t, fmt.Sprintf("60000%s", types.GetNativeToken().MinUnit), fee.IssueFee.String())
	require.Equal(t, fmt.Sprintf("6000%s", types.GetNativeToken().MinUnit), fee.MintFee.String())
}

package keeper_test

import (
	"fmt"

	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github/irismod/token"
	simapp "github/irismod/token/app"
	"github/irismod/token/exported"
	"github/irismod/token/types"
)

func TestQueryToken(t *testing.T) {
	app := simapp.Setup(isCheck)
	ctx := app.BaseApp.NewContext(isCheck, abci.Header{})
	querier := token.NewQuerier(app.TokenKeeper)

	params := token.QueryTokenParams{
		Denom: types.GetNativeToken().Symbol,
	}
	bz := app.Codec().MustMarshalJSON(params)
	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", token.ModuleName, token.QueryToken),
		Data: bz,
	}

	data, err := querier(ctx, []string{token.QueryToken}, query)
	require.Nil(t, err)

	data2 := app.Codec().MustMarshalJSON(token.GetNativeToken())
	require.EqualValues(t, data2, data)

	//query by mint_unit
	params = token.QueryTokenParams{
		Denom: types.GetNativeToken().MinUnit,
	}

	bz = app.Codec().MustMarshalJSON(params)
	query = abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", token.ModuleName, token.QueryToken),
		Data: bz,
	}

	data, err = querier(ctx, []string{token.QueryToken}, query)
	require.Nil(t, err)

	data2 = app.Codec().MustMarshalJSON(token.GetNativeToken())
	require.EqualValues(t, data2, data)

}

func TestQueryTokens(t *testing.T) {
	app := simapp.Setup(isCheck)
	ctx := app.BaseApp.NewContext(isCheck, abci.Header{})
	querier := token.NewQuerier(app.TokenKeeper)

	params := token.QueryTokensParams{
		Owner: nil,
	}
	bz := app.Codec().MustMarshalJSON(params)
	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", token.ModuleName, token.QueryTokens),
		Data: bz,
	}

	data, err := querier(ctx, []string{token.QueryTokens}, query)
	require.Nil(t, err)

	data2 := app.Codec().MustMarshalJSON([]exported.TokenI{token.GetNativeToken()})
	require.EqualValues(t, data2, data)
}

func TestQueryFees(t *testing.T) {
	app := simapp.Setup(isCheck)
	ctx := app.BaseApp.NewContext(isCheck, abci.Header{})
	querier := token.NewQuerier(app.TokenKeeper)

	params := token.QueryTokenFeesParams{
		Symbol: "btc",
	}
	bz := app.Codec().MustMarshalJSON(params)
	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", token.ModuleName, token.QueryFees),
		Data: bz,
	}

	data, err := querier(ctx, []string{token.QueryFees}, query)
	require.Nil(t, err)

	var fee token.TokenFees
	app.Codec().MustUnmarshalJSON(data, &fee)
	require.Equal(t, false, fee.Exist)
	require.Equal(t, fmt.Sprintf("60000%s", types.GetNativeToken().MinUnit), fee.IssueFee.String())
	require.Equal(t, fmt.Sprintf("6000%s", types.GetNativeToken().MinUnit), fee.MintFee.String())
}

package token_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github/irismod/token"
	simapp "github/irismod/token/app"
	"github/irismod/token/types"
)

func TestExportGenesis(t *testing.T) {
	app := simapp.Setup(false)

	ctx := app.BaseApp.NewContext(false, abci.Header{})

	// export genesis
	genesisState := token.ExportGenesis(ctx, app.TokenKeeper)

	require.Equal(t, types.DefaultParams(), genesisState.Params)
	for _, token := range genesisState.Tokens {
		require.Equal(t, token, types.GetNativeToken())
	}
}

func TestInitGenesis(t *testing.T) {
	app := simapp.Setup(false)

	ctx := app.BaseApp.NewContext(false, abci.Header{})

	// add token
	addr := sdk.AccAddress([]byte("addr1"))
	ft := types.NewToken("btc", "Bitcoin Network", "satoshi", 1, 1, 1, true, addr)

	genesis := types.GenesisState{
		Params: types.DefaultParams(),
		Tokens: types.Tokens{ft},
	}

	// initialize genesis
	token.InitGenesis(ctx, app.TokenKeeper, genesis)

	// query all tokens
	var tokens = app.TokenKeeper.GetTokens(ctx, nil)
	require.Equal(t, len(tokens), 2)
	require.Equal(t, tokens[0], ft)
}

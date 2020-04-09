package asset_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github/irismod/asset"
	simapp "github/irismod/asset/app"
	"github/irismod/asset/internal/types"
)

func TestExportGenesis(t *testing.T) {
	app := simapp.Setup(false)

	ctx := app.BaseApp.NewContext(false, abci.Header{})

	// add token
	addr := sdk.AccAddress([]byte("addr1"))
	ft := types.NewFungibleToken("btc", "Bitcoin Network", "satoshi", 1, 1, 1, true, addr)

	genesis := types.GenesisState{
		Params: types.DefaultParams(),
		Tokens: types.Tokens{ft},
	}

	// initialize genesis
	asset.InitGenesis(ctx, app.AssetKeeper, genesis)

	// query all tokens
	var tokens types.Tokens
	app.AssetKeeper.IterateTokens(ctx, func(token types.FungibleToken) (stop bool) {
		tokens = append(tokens, token)
		return false
	})

	require.Equal(t, len(tokens), 1)

	// export genesis
	genesisState := asset.ExportGenesis(ctx, app.AssetKeeper)

	require.Equal(t, types.DefaultParams(), genesisState.Params)
	for _, token := range genesisState.Tokens {
		require.Equal(t, token, ft)
	}
}

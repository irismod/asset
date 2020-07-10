package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irismod/token/keeper"
	"github.com/irismod/token/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	k.SetParamSet(ctx, data.Params)

	//init tokens
	for _, token := range data.Tokens {
		if err := k.AddToken(ctx, token); err != nil {
			panic(err.Error())
		}
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	var tokens types.Tokens
	for _, token := range k.GetTokens(ctx, nil) {
		tokens = append(tokens, token.(types.Token))
	}
	return types.GenesisState{
		Params: k.GetParamSet(ctx),
		Tokens: tokens,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() types.GenesisState {
	return types.GenesisState{
		Params: types.DefaultParams(),
		Tokens: []types.Token{types.GetNativeToken()},
	}
}

// ValidateGenesis validates the provided token genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data types.GenesisState) error {
	if err := types.ValidateParams(data.Params); err != nil {
		return err
	}

	// validate token
	for _, token := range data.Tokens {
		if err := types.ValidateToken(token); err != nil {
			return err
		}
	}
	return nil
}

package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github/irismod/token/internal/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	if err := ValidateGenesis(data); err != nil && err != types.ErrNilOwner {
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
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var tokens Tokens
	for _, token := range k.GetTokens(ctx, nil) {
		tokens = append(tokens, token.(FungibleToken))
	}
	return GenesisState{
		Params: k.GetParamSet(ctx),
		Tokens: tokens,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
		Tokens: []FungibleToken{types.GetNativeToken()},
	}
}

// ValidateGenesis validates the provided asset genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := ValidateParams(data.Params); err != nil {
		return err
	}

	// validate tokens
	if err := data.Tokens.Validate(); err != nil {
		return err
	}

	return nil
}

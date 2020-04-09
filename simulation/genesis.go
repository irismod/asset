package simulation

// DONTCOVER

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github/irismod/asset/internal/types"
)

// Simulation parameter constants
const (
	AssetTaxRate      = "asset_tax_rate"
	IssueTokenBaseFee = "issue_token_base_fee"
	MintTokenFeeRatio = "mint_token_fee_ratio"
)

// RandomDec randomized sdk.RandomDec
func RandomDec(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(r.Int63())
}

// RandomInt randomized sdk.Int
func RandomInt(r *rand.Rand) sdk.Int {
	return sdk.NewInt(r.Int63())
}

// RandomizedGenState generates a random GenesisState for bank
func RandomizedGenState(simState *module.SimulationState) {

	var assetTaxRate sdk.Dec
	var issueTokenBaseFee sdk.Int
	var mintTokenFeeRatio sdk.Dec

	simState.AppParams.GetOrGenerate(
		simState.Cdc, AssetTaxRate, &assetTaxRate, simState.Rand,
		func(r *rand.Rand) { assetTaxRate = sdk.NewDecWithPrec(int64(r.Intn(5)), 1) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, IssueTokenBaseFee, &issueTokenBaseFee, simState.Rand,
		func(r *rand.Rand) {
			issueTokenBaseFee = sdk.NewInt(int64(10))
		},
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, MintTokenFeeRatio, &mintTokenFeeRatio, simState.Rand,
		func(r *rand.Rand) { mintTokenFeeRatio = sdk.NewDecWithPrec(int64(r.Intn(5)), 1) },
	)

	assetGenesis := types.NewGenesisState(
		types.NewParams(assetTaxRate, sdk.NewCoin(sdk.DefaultBondDenom, issueTokenBaseFee), mintTokenFeeRatio),
		types.Tokens{},
	)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(assetGenesis)

	fmt.Printf("Selected randomly generated asset parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, assetGenesis))

}

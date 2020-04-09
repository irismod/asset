package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github/irismod/asset/internal/types"
)

const (
	keyAssetTaxRate      = "AssetTaxRate"
	keyIssueTokenBaseFee = "IssueTokenBaseFee"
	keyMintTokenFeeRatio = "MintTokenFeeRatio"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simulation.ParamChange {
	return []simulation.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyAssetTaxRate,
			func(r *rand.Rand) string {
				return fmt.Sprintf("%s", RandomDec(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyIssueTokenBaseFee,
			func(r *rand.Rand) string {
				return fmt.Sprintf("%s", RandomInt(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyMintTokenFeeRatio,
			func(r *rand.Rand) string {
				return fmt.Sprintf("%s", RandomDec(r))
			},
		),
	}
}

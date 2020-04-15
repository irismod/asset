package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/irismod/token/types"
)

const (
	keyTokenTaxRate      = "TokenTaxRate"
	keyIssueTokenBaseFee = "IssueTokenBaseFee"
	keyMintTokenFeeRatio = "MintTokenFeeRatio"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simulation.ParamChange {
	return []simulation.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyTokenTaxRate,
			func(r *rand.Rand) string {
				return RandomDec(r).String()
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyIssueTokenBaseFee,
			func(r *rand.Rand) string {
				return RandomInt(r).String()
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyMintTokenFeeRatio,
			func(r *rand.Rand) string {
				return RandomDec(r).String()
			},
		),
	}
}

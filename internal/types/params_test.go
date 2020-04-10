package types

import (
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestValidateParams(t *testing.T) {
	defaultToken := GetNativeToken()
	tests := []struct {
		testCase string
		Params
		expectPass bool
	}{
		{"Minimum value",
			Params{
				AssetTaxRate:      sdk.ZeroDec(),
				MintTokenFeeRatio: sdk.ZeroDec(),
				IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.ZeroInt()),
			},
			true,
		},
		{"Maximum value",
			Params{
				AssetTaxRate:      sdk.NewDec(1),
				MintTokenFeeRatio: sdk.NewDec(1),
				IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(math.MaxInt64)),
			},
			true,
		},
		{"AssetTaxRate less than the maximum",
			Params{
				AssetTaxRate:      sdk.NewDecWithPrec(-1, 1),
				MintTokenFeeRatio: sdk.NewDec(0),
				IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(1)),
			},
			false,
		},
		{"MintTokenFeeRatio less than the maximum",
			Params{
				AssetTaxRate:      sdk.NewDec(0),
				MintTokenFeeRatio: sdk.NewDecWithPrec(-1, 1),
				IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(1)),
			},
			false,
		},
		{"AssetTaxRate greater than the maximum",
			Params{
				AssetTaxRate:      sdk.NewDecWithPrec(11, 1),
				MintTokenFeeRatio: sdk.NewDec(1),
				IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(1)),
			},
			false,
		},
		{"MintTokenFeeRatio greater than the maximum",
			Params{
				AssetTaxRate:      sdk.NewDec(1),
				MintTokenFeeRatio: sdk.NewDecWithPrec(11, 1),
				IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(1)),
			},
			false,
		},
		{"IssueTokenBaseFee is negative",
			Params{
				AssetTaxRate:      sdk.NewDec(1),
				MintTokenFeeRatio: sdk.NewDec(1),
				IssueTokenBaseFee: sdk.Coin{Denom: defaultToken.Symbol, Amount: sdk.NewInt(-1)},
			},
			false,
		},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, ValidateParams(tc.Params), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, ValidateParams(tc.Params), "test: %v", tc.testCase)
		}
	}
}

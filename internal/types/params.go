package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// parameter keys
var (
	KeyAssetTaxRate      = []byte("AssetTaxRate")
	KeyIssueTokenBaseFee = []byte("IssueTokenBaseFee")
	KeyMintTokenFeeRatio = []byte("MintTokenFeeRatio")
)

// asset params
type Params struct {
	AssetTaxRate      sdk.Dec  `json:"asset_tax_rate"`       // e.g., 40%
	IssueTokenBaseFee sdk.Coin `json:"issue_token_base_fee"` // e.g., 300000*10^18iris-atto
	MintTokenFeeRatio sdk.Dec  `json:"mint_token_fee_ratio"` // e.g., 10%
} // issuance fee = IssueTokenBaseFee / (ln(len(symbol))/ln3)^4

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAssetTaxRate, &p.AssetTaxRate, validateTaxRate),
		paramtypes.NewParamSetPair(KeyIssueTokenBaseFee, &p.IssueTokenBaseFee, validateIssueTokenBaseFee),
		paramtypes.NewParamSetPair(KeyMintTokenFeeRatio, &p.MintTokenFeeRatio, validateMintTokenFeeRatio),
	}
}

// NewParams asset params constructor
func NewParams(assetTaxRate sdk.Dec, issueTokenBaseFee sdk.Coin,
	mintTokenFeeRatio sdk.Dec,
) Params {
	return Params{
		AssetTaxRate:      assetTaxRate,
		IssueTokenBaseFee: issueTokenBaseFee,
		MintTokenFeeRatio: mintTokenFeeRatio,
	}
}

// ParamTypeTable returns the TypeTable for coinswap module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// default asset module params
func DefaultParams() Params {
	defaultToken := GetNativeToken()
	return Params{
		AssetTaxRate:      sdk.NewDecWithPrec(4, 1), // 0.4 (40%)
		IssueTokenBaseFee: sdk.NewCoin(defaultToken.Symbol, sdk.NewIntWithDecimal(60000, int(defaultToken.Scale))),
		MintTokenFeeRatio: sdk.NewDecWithPrec(1, 1), // 0.1 (10%)
	}
}

func ValidateParams(p Params) error {
	if err := validateTaxRate(p.AssetTaxRate); err != nil {
		return err
	}
	if err := validateMintTokenFeeRatio(p.MintTokenFeeRatio); err != nil {
		return err
	}
	if err := validateIssueTokenBaseFee(p.IssueTokenBaseFee); err != nil {
		return err
	}

	return nil
}

func validateTaxRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.NewDec(1)) || v.LT(sdk.ZeroDec()) {
		return fmt.Errorf("asset tax rate [%s] should be between [0, 1]", v.String())
	}
	return nil
}

func validateMintTokenFeeRatio(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.GT(sdk.NewDec(1)) || v.LT(sdk.ZeroDec()) {
		return fmt.Errorf("fee ratio for minting tokens [%s] should be between [0, 1]", v.String())
	}
	return nil
}

func validateIssueTokenBaseFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("base fee for issuing token should not be negative")
	}
	return nil
}

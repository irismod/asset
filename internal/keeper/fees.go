//nolint
package keeper

import (
	"math"
	"strconv"

	"github/irismod/asset/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// fee factor formula: (ln(len({name}))/ln{base})^{exp}
const (
	FeeFactorBase = 3
	FeeFactorExp  = 4
)

// DeductIssueTokenFee performs fee handling for issuing token
func (k Keeper) DeductIssueTokenFee(ctx sdk.Context, owner sdk.AccAddress, symbol string) error {
	// get the required issuance fee
	fee := k.GetTokenIssueFee(ctx, symbol)
	return feeHandler(ctx, k, owner, fee)
}

// DeductMintTokenFee performs fee handling for minting token
func (k Keeper) DeductMintTokenFee(ctx sdk.Context, owner sdk.AccAddress, symbol string) error {
	// get the required minting fee
	fee := k.GetTokenMintFee(ctx, symbol)
	return feeHandler(ctx, k, owner, fee)
}

// GetTokenIssueFee returns the token issurance fee
func (k Keeper) GetTokenIssueFee(ctx sdk.Context, symbol string) sdk.Coin {
	// get params
	params := k.GetParamSet(ctx)
	issueTokenBaseFee := params.IssueTokenBaseFee

	// compute the fee
	fee := calcFeeByBase(symbol, issueTokenBaseFee.Amount)

	return sdk.NewCoin(sdk.DefaultBondDenom, convertFeeToInt(fee))
}

// GetTokenMintFee returns the token minting fee
func (k Keeper) GetTokenMintFee(ctx sdk.Context, symbol string) sdk.Coin {
	// get params
	params := k.GetParamSet(ctx)
	mintTokenFeeRatio := params.MintTokenFeeRatio

	// compute the issurance and minting fees
	issueFee := k.GetTokenIssueFee(ctx, symbol)
	mintFee := sdk.NewDecFromInt(issueFee.Amount).Mul(mintTokenFeeRatio)

	return sdk.NewCoin(sdk.DefaultBondDenom, convertFeeToInt(mintFee))
}

// feeHandler handles the fee of asset
func feeHandler(ctx sdk.Context, k Keeper, feeAcc sdk.AccAddress, fee sdk.Coin) error {
	params := k.GetParamSet(ctx)
	assetTaxRate := params.AssetTaxRate

	// compute community tax and burned coin
	communityTaxCoin := sdk.NewCoin(fee.Denom, sdk.NewDecFromInt(fee.Amount).Mul(assetTaxRate).TruncateInt())
	burnedCoins := sdk.NewCoins(fee.Sub(communityTaxCoin))

	// send all fees to module account
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(
		ctx, feeAcc, types.ModuleName, sdk.NewCoins(fee),
	); err != nil {
		return err
	}

	// send community tax to collectedFees
	if err := k.addCollectedFees(ctx, sdk.NewCoins(communityTaxCoin)); err != nil {
		return err
	}

	// burn burnedCoin
	return k.supplyKeeper.BurnCoins(ctx, types.ModuleName, burnedCoins)
}

// calcFeeByBase computes the actual fee according to the given base fee
func calcFeeByBase(name string, baseFee sdk.Int) sdk.Dec {
	feeFactor := calcFeeFactor(name)
	actualFee := sdk.NewDecFromInt(baseFee).Quo(feeFactor)

	return actualFee
}

// calcFeeFactor computes the fee factor of the given name
// Note: make sure that the name size is examined before invoking the function
func calcFeeFactor(name string) sdk.Dec {
	nameLen := len(name)
	if nameLen == 0 {
		panic("the length of name must be greater than 0")
	}

	denominator := math.Log(FeeFactorBase)
	numerator := math.Log(float64(nameLen))

	feeFactor := math.Pow(numerator/denominator, FeeFactorExp)
	feeFactorDec, err := sdk.NewDecFromStr(strconv.FormatFloat(feeFactor, 'f', 2, 64))
	if err != nil {
		panic("invalid string")
	}

	return feeFactorDec
}

// convertFeeToInt converts the given fee to Int.
// if greater than 1, rounds it; returns 1 otherwise
func convertFeeToInt(fee sdk.Dec) sdk.Int {
	power := sdk.TokensToConsensusPower(fee.TruncateInt())
	if power > 1 {
		return sdk.TokensFromConsensusPower(power)
	} else {
		return sdk.TokensFromConsensusPower(1)
	}
}

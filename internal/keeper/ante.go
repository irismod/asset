package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github/irismod/token/internal/types"
)

type ValidateTokenFeeDecorator struct {
	k  Keeper
	ak types.AccountKeeper
}

func NewValidateTokenFeeDecorator(k Keeper, ak types.AccountKeeper) ValidateTokenFeeDecorator {
	return ValidateTokenFeeDecorator{
		k:  k,
		ak: ak,
	}
}

// AnteHandle returns an AnteHandler that checks if the balance of
// the fee payer is sufficient for asset related fee
func (dtf ValidateTokenFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {

	// new ctx
	newCtx = sdk.Context{}
	// total fee
	feeMap := make(map[string]sdk.Coin)
	for _, msg := range tx.GetMsgs() {
		// only check consecutive msgs which are routed to asset from the beginning
		if msg.Route() != types.ModuleName {
			break
		}

		switch msg := msg.(type) {
		case types.MsgIssueToken:
			fee := dtf.k.GetTokenIssueFee(ctx, msg.Symbol)
			if fe, ok := feeMap[msg.Owner.String()]; ok {
				feeMap[msg.Owner.String()] = fe.Add(fee)
			} else {
				feeMap[msg.Owner.String()] = fee
			}
		case types.MsgMintToken:
			fee := dtf.k.GetTokenMintFee(ctx, msg.Symbol)
			if fe, ok := feeMap[msg.Owner.String()]; ok {
				feeMap[msg.Owner.String()] = fe.Add(fee)
			} else {
				feeMap[msg.Owner.String()] = fee
			}
		}
	}

	for addr, fee := range feeMap {
		owner, _ := sdk.AccAddressFromBech32(addr)
		account := dtf.ak.GetAccount(ctx, owner)
		balance := account.GetCoins()
		if balance.IsAllLT(sdk.NewCoins(fee)) {
			return ctx, sdkerrors.Wrapf(
				sdkerrors.ErrInsufficientFunds, "insufficient coins for asset fee; %s < %s", balance, fee)
		}
	}
	// continue
	return next(ctx, tx, simulate)
}

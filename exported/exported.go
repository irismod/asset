package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenI interface {
	GetSymbol() string
	GetName() string
	GetScale() uint8
	GetMinUnit() string
	GetInitialSupply() uint64
	GetMaxSupply() uint64
	GetMintable() bool
	GetOwner() sdk.AccAddress

	ToMainCoin(coin sdk.Coin) (sdk.DecCoin, error)
	ToMinCoin(coin sdk.DecCoin) (sdk.Coin, error)
}

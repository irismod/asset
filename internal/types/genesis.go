package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params"` // asset params
	Tokens Tokens `json:"tokens"` // issued tokens
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, tokens Tokens) GenesisState {
	return GenesisState{
		Params: params,
		Tokens: tokens,
	}
}

//SetNativeToken reset the system's default native token
func SetNativeToken(symbol,
	name,
	minUnit string,
	decimal uint8,
	initialSupply,
	maxSupply uint64,
	mintable bool,
	owner sdk.AccAddress) {
	nativeToken = NewToken(symbol, name, minUnit, decimal, initialSupply, maxSupply, mintable, owner)
}

func GetNativeToken() Token {
	return nativeToken
}

var nativeToken = Token{
	Symbol:        sdk.DefaultBondDenom,
	Name:          "Network staking token ",
	Scale:         0,
	MinUnit:       sdk.DefaultBondDenom,
	InitialSupply: 2000000000,
	MaxSupply:     10000000000,
	Mintable:      true,
}

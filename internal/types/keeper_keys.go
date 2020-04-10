package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	PrefixTokenForSymbol  = []byte{0x1} // symbol prefix for the token
	PrefixTokenForMinUint = []byte{0x2} // min_unit prefix for the token
	PrefixTokens          = []byte{0x3} // prefix for the tokens
)

// KeySymbol returns the key of the token with the specified symbol
func KeySymbol(symbol string) []byte {
	return append(PrefixTokenForSymbol, []byte(symbol)...)
}

// KeyMinUint returns the key of the token with the specified min_unit
func KeyMinUint(minUnit string) []byte {
	return append(PrefixTokenForMinUint, []byte(minUnit)...)
}

// KeyTokens returns the key of the specified owner and symbol. Intended for querying all tokens of an owner
func KeyTokens(owner sdk.AccAddress, symbol string) []byte {
	return append(append(PrefixTokens, owner.Bytes()...), []byte(symbol)...)
}

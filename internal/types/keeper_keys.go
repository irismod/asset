package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	PrefixToken  = []byte("token:")       // prefix for the token store
	PrefixTokens = []byte("ownerTokens:") // prefix for the tokens store
)

// KeyToken returns the key of the token with the specified symbol
func KeyToken(symbol string) []byte {
	return append(PrefixToken, []byte(symbol)...)
}

// KeyTokens returns the key of the specified owner and symbol. Intended for querying all tokens of an owner
func KeyTokens(owner sdk.AccAddress, symbol string) []byte {
	return append(append(PrefixTokens, owner.Bytes()...), []byte(symbol)...)
}

package simulation

// DONTCOVER

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	kv "github.com/tendermint/tendermint/libs/kv"

	"github/irismod/asset/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding asset type
func DecodeStore(cdc *codec.Codec, kvA, kvB kv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.PrefixTokenForSymbol):
		var tokenA, tokenB types.FungibleToken
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &tokenA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &tokenB)
		return fmt.Sprintf("%v\n%v", tokenA, tokenB)
	case bytes.Equal(kvA.Key[:1], types.PrefixTokens):
		var symbolA, symbolB string
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &symbolA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &symbolB)
		return fmt.Sprintf("%v\n%v", symbolA, symbolB)
	default:
		panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
	}
}

package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterCodec registers the necessary x/bank interfaces and concrete types
// on the provided Amino codec. These types are used for Amino JSON serialization.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*TokenI)(nil), nil)

	cdc.RegisterConcrete(&Token{}, "irismod/token/Token", nil)

	cdc.RegisterConcrete(&MsgIssueToken{}, "irismod/token/MsgIssueToken", nil)
	cdc.RegisterConcrete(&MsgEditToken{}, "irismod/token/MsgEditToken", nil)
	cdc.RegisterConcrete(&MsgMintToken{}, "irismod/token/MsgMintToken", nil)
	cdc.RegisterConcrete(&MsgTransferTokenOwner{}, "irismod/token/MsgTransferTokenOwner", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgIssueToken{},
		&MsgEditToken{},
		&MsgMintToken{},
		&MsgTransferTokenOwner{},
	)
	registry.RegisterInterface(
		"irismod.token.TokenI",
		(*TokenI)(nil),
		&Token{},
	)
}

var (
	amino = codec.New()

	// ModuleCdc references the global irismod/token module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewHybridCodec(amino, types.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(amino)
	cryptocodec.RegisterCrypto(amino)
}

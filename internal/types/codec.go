package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIssueToken{}, "irismod/asset/MsgIssueToken", nil)
	cdc.RegisterConcrete(MsgEditToken{}, "irismod/asset/MsgEditToken", nil)
	cdc.RegisterConcrete(MsgMintToken{}, "irismod/asset/MsgMintToken", nil)
	cdc.RegisterConcrete(MsgTransferTokenOwner{}, "irismod/asset/MsgTransferTokenOwner", nil)

	cdc.RegisterConcrete(FungibleToken{}, "irismod/asset/FungibleToken", nil)

	cdc.RegisterConcrete(&Params{}, "irismod/asset/Params", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	codec.RegisterCrypto(ModuleCdc)
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

//nolint
package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// asset module sentinel errors
var (
	ErrNilAssetOwner            = sdkerrors.Register(ModuleName, 1, "the owner of the asset must be specified")
	ErrInvalidAssetName         = sdkerrors.Register(ModuleName, 2, "invalid token name")
	ErrInvalidAssetMinUnitAlias = sdkerrors.Register(ModuleName, 3, "invalid token min_unit")
	ErrInvalidAssetSymbol       = sdkerrors.Register(ModuleName, 4, "must be standard denom")
	ErrInvalidAssetInitSupply   = sdkerrors.Register(ModuleName, 5, "invalid token initial supply")
	ErrInvalidAssetMaxSupply    = sdkerrors.Register(ModuleName, 6, "invalid token max supply")
	ErrInvalidAssetDecimal      = sdkerrors.Register(ModuleName, 7, "invalid token decimal")
	ErrAssetAlreadyExists       = sdkerrors.Register(ModuleName, 8, "insufficient funds")
	ErrAssetNotExists           = sdkerrors.Register(ModuleName, 9, "token does not exist")
	ErrInvalidAddress           = sdkerrors.Register(ModuleName, 10, "the owner of the token must be specified")
	ErrInvalidToAddress         = sdkerrors.Register(ModuleName, 11, "the new owner must not be same as the original owner")
	ErrInvalidOwner             = sdkerrors.Register(ModuleName, 12, "invalid token owner")
	ErrAssetNotMintable         = sdkerrors.Register(ModuleName, 13, "the token is set to be non-mintable")
)

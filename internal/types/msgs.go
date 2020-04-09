package types

import (
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// MsgRoute identifies transaction types
	MsgRoute = "asset"

	TypeMsgIssueToken         = "issue_token"
	TypeMsgEditToken          = "edit_token"
	TypeMsgMintToken          = "mint_token"
	TypeMsgTransferTokenOwner = "transfer_token_owner"

	// constant used to indicate that some field should not be updated
	DoNotModify = "[do-not-modify]"

	MaximumAssetMaxSupply  = uint64(1000000000000) // maximal limitation for asset max supply，1000 billion
	MaximumAssetInitSupply = uint64(100000000000)  // maximal limitation for asset initial supply，100 billion
	MaximumAssetDecimal    = uint8(18)             // maximal limitation for asset decimal
	MinimumAssetSymbolLen  = 3                     // minimal limitation for the length of the asset's symbol / canonical_symbol
	MaximumAssetSymbolLen  = 8                     // maximal limitation for the length of the asset's symbol / canonical_symbol
	MaximumAssetNameLen    = 32                    // maximal limitation for the length of the asset's name
	MinimumAssetMinUnitLen = 3                     // minimal limitation for the length of the asset's min_unit
	MaximumAssetMinUnitLen = 10                    // maximal limitation for the length of the asset's min_unit
)

var (
	IsAlphaNumeric   = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString // only accepts alphanumeric characters
	IsBeginWithAlpha = regexp.MustCompile(`^[a-zA-Z].*`).MatchString
)

var _, _, _, _ sdk.Msg = &MsgIssueToken{}, &MsgEditToken{}, &MsgMintToken{}, &MsgTransferTokenOwner{}

// MsgIssueToken
type MsgIssueToken struct {
	Symbol        string         `json:"symbol"`
	Name          string         `json:"name"`
	Scale         uint8          `json:"scale"`
	MinUnit       string         `json:"min_unit"`
	InitialSupply uint64         `json:"initial_supply"`
	MaxSupply     uint64         `json:"max_supply"`
	Mintable      bool           `json:"mintable"`
	Owner         sdk.AccAddress `json:"owner"`
}

// NewMsgIssueToken - construct asset issue msg.
func NewMsgIssueToken(symbol string, minUnit string, name string, scale uint8, initialSupply, maxSupply uint64, mintable bool, owner sdk.AccAddress) MsgIssueToken {
	return MsgIssueToken{
		Symbol:        symbol,
		Name:          name,
		Scale:         scale,
		MinUnit:       minUnit,
		InitialSupply: initialSupply,
		MaxSupply:     maxSupply,
		Mintable:      mintable,
		Owner:         owner,
	}
}

// Implements Msg.
func (msg MsgIssueToken) Route() string { return MsgRoute }
func (msg MsgIssueToken) Type() string  { return TypeMsgIssueToken }

func ValidateMsgIssueToken(msg MsgIssueToken) error {
	msg.Name = strings.TrimSpace(msg.Name)
	msg.MinUnit = strings.TrimSpace(msg.MinUnit)

	if msg.MaxSupply == 0 {
		if msg.Mintable {
			msg.MaxSupply = MaximumAssetMaxSupply
		} else {
			msg.MaxSupply = msg.InitialSupply
		}
	}

	if msg.Owner.Empty() {
		return ErrNilAssetOwner
	}

	nameLen := len(msg.Name)
	if nameLen == 0 || nameLen > MaximumAssetNameLen {
		return ErrInvalidAssetName
	}

	if err := CheckSymbol(msg.Symbol); err != nil {
		return err
	}

	minUnitAliasLen := len(msg.MinUnit)
	if minUnitAliasLen > 0 && (minUnitAliasLen < MinimumAssetMinUnitLen || minUnitAliasLen > MaximumAssetMinUnitLen || !IsAlphaNumeric(msg.MinUnit) || !IsBeginWithAlpha(msg.MinUnit)) {
		return ErrInvalidAssetMinUnitAlias
	}

	if msg.InitialSupply > MaximumAssetInitSupply {
		return ErrInvalidAssetInitSupply
	}

	if msg.MaxSupply < msg.InitialSupply || msg.MaxSupply > MaximumAssetMaxSupply {
		return ErrInvalidAssetMaxSupply
	}

	if msg.Scale > MaximumAssetDecimal {
		return ErrInvalidAssetDecimal
	}

	return nil
}

// Implements Msg.
func (msg MsgIssueToken) ValidateBasic() error {
	return ValidateMsgIssueToken(msg)
}

// Implements Msg.
func (msg MsgIssueToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgIssueToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgTransferTokenOwner for transferring the token owner
type MsgTransferTokenOwner struct {
	SrcOwner sdk.AccAddress `json:"src_owner"` // the current owner address of the token
	DstOwner sdk.AccAddress `json:"dst_owner"` // the new owner
	Symbol   string         `json:"symbol"`    // the token symbol
}

func NewMsgTransferTokenOwner(srcOwner, dstOwner sdk.AccAddress, symbol string) MsgTransferTokenOwner {
	symbol = strings.TrimSpace(symbol)

	return MsgTransferTokenOwner{
		SrcOwner: srcOwner,
		DstOwner: dstOwner,
		Symbol:   symbol,
	}
}

// GetSignBytes implements Msg
func (msg MsgTransferTokenOwner) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgTransferTokenOwner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.SrcOwner}
}

func (msg MsgTransferTokenOwner) ValidateBasic() error {
	// check the SrcOwner
	if len(msg.SrcOwner) == 0 {
		return ErrInvalidAddress
	}

	// check if the `DstOwner` is empty
	if len(msg.DstOwner) == 0 {
		return ErrInvalidAddress
	}

	// check if the `DstOwner` is same as the original owner
	if msg.SrcOwner.Equals(msg.DstOwner) {
		return ErrInvalidToAddress
	}

	// check the symbol
	if err := CheckSymbol(msg.Symbol); err != nil {
		return err
	}

	return nil
}

// Route implements Msg
func (msg MsgTransferTokenOwner) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgTransferTokenOwner) Type() string { return TypeMsgTransferTokenOwner }

// MsgEditToken for editing a specified token
type MsgEditToken struct {
	Symbol    string         `json:"symbol"` //  symbol of token
	Owner     sdk.AccAddress `json:"owner"`  //  owner of token
	MaxSupply uint64         `json:"max_supply"`
	Mintable  Bool           `json:"mintable"` //  mintable of token
	Name      string         `json:"name"`
}

// NewMsgEditToken creates a MsgEditToken
func NewMsgEditToken(name, symbol string, maxSupply uint64, mintable Bool, owner sdk.AccAddress) MsgEditToken {
	name = strings.TrimSpace(name)

	return MsgEditToken{
		Name:      name,
		Symbol:    symbol,
		MaxSupply: maxSupply,
		Mintable:  mintable,
		Owner:     owner,
	}
}

// Route implements Msg
func (msg MsgEditToken) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgEditToken) Type() string { return TypeMsgEditToken }

// ValidateBasic implements Msg
func (msg MsgEditToken) ValidateBasic() error {
	// check owner
	if msg.Owner.Empty() {
		return ErrNilAssetOwner
	}

	nameLen := len(msg.Name)
	if DoNotModify != msg.Name && nameLen > MaximumAssetNameLen {
		return ErrInvalidAssetName
	}

	// check max_supply for fast failed
	if msg.MaxSupply > MaximumAssetMaxSupply {
		return ErrInvalidAssetMaxSupply
	}

	// check symbol
	if err := CheckSymbol(msg.Symbol); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements Msg
func (msg MsgEditToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgEditToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgMintToken for minting the token to a specified address
type MsgMintToken struct {
	Symbol string         `json:"symbol"` // the symbol of the token
	Owner  sdk.AccAddress `json:"owner"`  // the current owner address of the token
	To     sdk.AccAddress `json:"to"`     // address of minting token to
	Amount uint64         `json:"amount"` // amount of minting token
}

// NewMsgMintToken creates a MsgMintToken
func NewMsgMintToken(symbol string, owner, to sdk.AccAddress, amount uint64) MsgMintToken {
	symbol = strings.TrimSpace(symbol)

	return MsgMintToken{
		Symbol: symbol,
		Owner:  owner,
		To:     to,
		Amount: amount,
	}
}

// Route implements Msg
func (msg MsgMintToken) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgMintToken) Type() string { return TypeMsgMintToken }

// GetSignBytes implements Msg
func (msg MsgMintToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgMintToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// ValidateBasic implements Msg
func (msg MsgMintToken) ValidateBasic() error {
	// check the owner
	if len(msg.Owner) == 0 {
		return ErrInvalidAddress
	}

	if msg.Amount == 0 || msg.Amount > MaximumAssetMaxSupply {
		return ErrInvalidAssetMaxSupply
	}

	return CheckSymbol(msg.Symbol)
}

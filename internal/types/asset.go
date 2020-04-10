package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FungibleToken defines a struct for the fungible token
type FungibleToken struct {
	Symbol        string         `json:"symbol" yaml:"symbol"`
	Name          string         `json:"name" yaml:"name"`
	Scale         uint8          `json:"scale" yaml:"scale"`
	MinUnit       string         `json:"min_unit" yaml:"min_unit"`
	InitialSupply uint64         `json:"initial_supply" yaml:"initial_supply"`
	MaxSupply     uint64         `json:"max_supply" yaml:"max_supply"`
	Mintable      bool           `json:"mintable" yaml:"mintable"`
	Owner         sdk.AccAddress `json:"owner" yaml:"owner"`
}

// NewFungibleToken constructs a new FungibleToken instance
func NewFungibleToken(
	symbol,
	name,
	minUnit string,
	scale uint8,
	initialSupply,
	maxSupply uint64,
	mintable bool,
	owner sdk.AccAddress,
) FungibleToken {
	return FungibleToken{
		Symbol:        symbol,
		Name:          name,
		MinUnit:       minUnit,
		Scale:         scale,
		InitialSupply: initialSupply,
		MaxSupply:     maxSupply,
		Mintable:      mintable,
		Owner:         owner,
	}
}

// GetSymbol implements exported.TokenI
func (ft FungibleToken) GetSymbol() string {
	return ft.Symbol
}

// GetName implements exported.TokenI
func (ft FungibleToken) GetName() string {
	return ft.Name
}

// GetScale implements exported.TokenI
func (ft FungibleToken) GetScale() uint8 {
	return ft.Scale
}

// GetMinUnit implements exported.TokenI
func (ft FungibleToken) GetMinUnit() string {
	return ft.MinUnit
}

// GetInitialSupply implements exported.TokenI
func (ft FungibleToken) GetInitialSupply() uint64 {
	return ft.InitialSupply
}

// GetMaxSupply implements exported.TokenI
func (ft FungibleToken) GetMaxSupply() uint64 {
	return ft.MaxSupply
}

// GetMintable implements exported.TokenI
func (ft FungibleToken) GetMintable() bool {
	return ft.Mintable
}

// GetOwner implements exported.TokenI
func (ft FungibleToken) GetOwner() sdk.AccAddress {
	return ft.Owner
}

// String implements fmt.Stringer
func (ft FungibleToken) String() string {
	return fmt.Sprintf(`FungibleToken:
  Name:              %s
  Symbol:            %s
  Scale:             %d
  MinUnit:           %s
  Initial Supply:    %d
  Max Supply:        %d
  Mintable:          %v
  Owner:             %s`,
		ft.Name, ft.Symbol, ft.Scale, ft.MinUnit,
		ft.InitialSupply, ft.MaxSupply, ft.Mintable, ft.Owner,
	)
}

// Tokens is a set of tokens
type Tokens []FungibleToken

// String implements Stringer
func (tokens Tokens) String() string {
	if len(tokens) == 0 {
		return "[]"
	}

	out := ""
	for _, token := range tokens {
		out += fmt.Sprintf("%s \n", token.String())
	}

	return out[:len(out)-1]
}

func (tokens Tokens) Validate() error {
	if len(tokens) == 0 {
		return nil
	}

	for _, token := range tokens {
		msg := NewMsgIssueToken(
			token.Symbol, token.MinUnit, token.Name, token.Scale,
			token.InitialSupply, token.MaxSupply, token.Mintable, token.Owner,
		)
		if err := ValidateMsgIssueToken(msg); err != nil {
			return err
		}
	}

	return nil
}

// CheckSymbol checks if the given symbol is valid
func CheckSymbol(symbol string) error {
	if len(symbol) < MinimumAssetSymbolLen || len(symbol) > MaximumAssetSymbolLen {
		return ErrInvalidAssetSymbol
	}

	if !IsBeginWithAlpha(symbol) || !IsAlphaNumeric(symbol) {
		return ErrInvalidAssetSymbol
	}

	return nil
}

type Bool string

const (
	False Bool = "false"
	True  Bool = "true"
	Nil   Bool = ""
)

func (b Bool) ToBool() bool {
	v := string(b)
	if len(v) == 0 {
		return false
	}
	result, _ := strconv.ParseBool(v)
	return result
}

func (b Bool) String() string {
	return string(b)
}

// Marshal needed for protobuf compatibility
func (b Bool) Marshal() ([]byte, error) {
	return []byte(b), nil
}

// Unmarshal needed for protobuf compatibility
func (b *Bool) Unmarshal(data []byte) error {
	*b = Bool(data[:])
	return nil
}

// Marshals to JSON using string
func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (b *Bool) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*b = Bool(s)
	return nil
}
func ParseBool(v string) (Bool, error) {
	if len(v) == 0 {
		return Nil, nil
	}
	result, err := strconv.ParseBool(v)
	if err != nil {
		return Nil, err
	}
	if result {
		return True, nil
	}
	return False, nil
}

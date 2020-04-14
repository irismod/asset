package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Token defines a struct for the fungible token
type Token struct {
	Symbol        string         `json:"symbol" yaml:"symbol"`
	Name          string         `json:"name" yaml:"name"`
	Scale         uint8          `json:"scale" yaml:"scale"`
	MinUnit       string         `json:"min_unit" yaml:"min_unit"`
	InitialSupply uint64         `json:"initial_supply" yaml:"initial_supply"`
	MaxSupply     uint64         `json:"max_supply" yaml:"max_supply"`
	Mintable      bool           `json:"mintable" yaml:"mintable"`
	Owner         sdk.AccAddress `json:"owner" yaml:"owner"`
}

// NewToken constructs a new Token instance
func NewToken(
	symbol,
	name,
	minUnit string,
	scale uint8,
	initialSupply,
	maxSupply uint64,
	mintable bool,
	owner sdk.AccAddress,
) Token {
	symbol = strings.ToLower(strings.TrimSpace(symbol))
	minUnit = strings.ToLower(strings.TrimSpace(minUnit))
	name = strings.TrimSpace(name)

	if maxSupply == 0 {
		if mintable {
			maxSupply = MaximumMaxSupply
		} else {
			maxSupply = initialSupply
		}
	}

	return Token{
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
func (t Token) GetSymbol() string {
	return t.Symbol
}

// GetName implements exported.TokenI
func (t Token) GetName() string {
	return t.Name
}

// GetScale implements exported.TokenI
func (t Token) GetScale() uint8 {
	return t.Scale
}

// GetMinUnit implements exported.TokenI
func (t Token) GetMinUnit() string {
	return t.MinUnit
}

// GetInitialSupply implements exported.TokenI
func (t Token) GetInitialSupply() uint64 {
	return t.InitialSupply
}

// GetMaxSupply implements exported.TokenI
func (t Token) GetMaxSupply() uint64 {
	return t.MaxSupply
}

// GetMintable implements exported.TokenI
func (t Token) GetMintable() bool {
	return t.Mintable
}

// GetOwner implements exported.TokenI
func (t Token) GetOwner() sdk.AccAddress {
	return t.Owner
}

//ToMainCoin return the main denom coin from args
func (t Token) ToMainCoin(coin sdk.Coin) (sdk.DecCoin, error) {
	if t.Symbol != coin.Denom && t.MinUnit != coin.Denom {
		return sdk.NewDecCoinFromDec(coin.Denom, sdk.ZeroDec()), sdkerrors.Wrapf(ErrTokenNotExists, "token not match")
	}

	if t.Symbol == coin.Denom {
		return sdk.NewDecCoin(coin.Denom, coin.Amount), nil
	}

	// dest amount = src amount / 10^(scale)
	precision := sdk.NewDecFromIntWithPrec(sdk.NewInt(1), int64(t.Scale))
	amount := sdk.NewDecFromInt(coin.Amount)

	amt := amount.Quo(precision)
	return sdk.NewDecCoinFromDec(t.Symbol, amt), nil
}

//ToMinCoin return the min denom coin from args
func (t Token) ToMinCoin(coin sdk.DecCoin) (newCoin sdk.Coin, err error) {
	if t.Symbol != coin.Denom && t.MinUnit != coin.Denom {
		return sdk.NewCoin(coin.Denom, sdk.ZeroInt()), sdkerrors.Wrapf(ErrTokenNotExists, "token not match")
	}

	if t.MinUnit == coin.Denom {
		return sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()), nil
	}

	// dest amount = src amount * 10^(dest scale)
	precision := sdk.NewDecFromIntWithPrec(sdk.NewInt(1), int64(t.Scale))
	amount := coin.Amount

	amt := amount.Mul(precision)
	return sdk.NewCoin(t.MinUnit, amt.TruncateInt()), nil
}

// String implements fmt.Stringer
func (t Token) String() string {
	return fmt.Sprintf(`Token:
  Name:              %s
  Symbol:            %s
  Scale:             %d
  MinUnit:           %s
  Initial Supply:    %d
  Max Supply:        %d
  Mintable:          %v
  Owner:             %s`,
		t.Name, t.Symbol, t.Scale, t.MinUnit,
		t.InitialSupply, t.MaxSupply, t.Mintable, t.Owner,
	)
}

// Tokens is a set of tokens
type Tokens []Token

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

func ValidateToken(token Token) error {
	if token.Owner.Empty() {
		return ErrNilOwner
	}

	nameLen := len(strings.TrimSpace(token.Name))
	if nameLen == 0 || nameLen > MaximumNameLen {
		return sdkerrors.Wrapf(ErrInvalidName, "invalid token name %s, only accepts length (0, %d]", token.Name, MaximumNameLen)
	}

	if err := CheckSymbol(token.Symbol); err != nil {
		return err
	}

	minUnitLen := len(strings.TrimSpace(token.MinUnit))
	if minUnitLen < MinimumMinUnitLen || minUnitLen > MaximumMinUnitLen || !IsAlphaNumeric(token.MinUnit) || !IsBeginWithAlpha(token.MinUnit) {
		return sdkerrors.Wrapf(ErrInvalidMinUnit, "invalid token min_unit %s, only accepts alphanumeric characters, and begin with an english letter, length [%d, %d]", token.MinUnit, MinimumMinUnitLen, MaximumMinUnitLen)
	}

	if token.InitialSupply > MaximumInitSupply {
		return sdkerrors.Wrapf(ErrInvalidInitSupply, "invalid token initial supply %d, only accepts value [0, %d]", token.InitialSupply, MaximumInitSupply)
	}

	if token.MaxSupply < token.InitialSupply || token.MaxSupply > MaximumMaxSupply {
		return sdkerrors.Wrapf(ErrInvalidMaxSupply, "invalid token max supply %d, only accepts value [%d, %d]", token.MaxSupply, token.InitialSupply, MaximumMaxSupply)
	}

	if token.Scale > MaximumScale {
		return sdkerrors.Wrapf(ErrInvalidScale, "invalid token scale %d, only accepts value [0, %d]", token.Scale, MaximumScale)
	}

	return nil
}

// CheckSymbol checks if the given symbol is valid
func CheckSymbol(symbol string) error {
	if len(symbol) < MinimumSymbolLen || len(symbol) > MaximumSymbolLen {
		return sdkerrors.Wrapf(ErrInvalidSymbol, "invalid symbol: %s,  only accepts length [%d, %d]", symbol, MinimumSymbolLen, MaximumSymbolLen)
	}

	if !IsBeginWithAlpha(symbol) || !IsAlphaNumeric(symbol) {
		return sdkerrors.Wrapf(ErrInvalidSymbol, "invalid symbol: %s, only accepts alphanumeric characters, and begin with an english letter", symbol)
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

// UnmarshalJSON from using string
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
package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryToken  = "token"
	QueryTokens = "tokens"
	QueryFees   = "fees"
)

// QueryTokenParams is the query parameters for 'custom/token/token'
type QueryTokenParams struct {
	Denom string
}

// QueryTokensParams is the query parameters for 'custom/token/tokens'
type QueryTokensParams struct {
	Owner sdk.AccAddress
}

// QueryTokenFeesParams is the query parameters for 'custom/token/fees'
type QueryTokenFeesParams struct {
	Symbol string
}

// TokenFees is used for the token fees query output
type TokenFees struct {
	Exist    bool     `json:"exist"`     // indicate if the token already exists
	IssueFee sdk.Coin `json:"issue_fee"` // issue fee
	MintFee  sdk.Coin `json:"mint_fee"`  // mint fee
}

// String implements stringer
func (tfo TokenFees) String() string {
	var out strings.Builder
	if tfo.Exist {
		out.WriteString("The symbol already exists\n")
	}

	out.WriteString(fmt.Sprintf(`Fees:
  IssueFee: %s
  MintFee:  %s`,
		tfo.IssueFee.String(), tfo.MintFee.String()))

	return out.String()
}

package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

// Rest variable names
// nolint
const (
	RestParamDenom  = "denom"
	RestParamSymbol = "symbol"
	RestParamOwner  = "owner"
)

// RegisterRoutes registers token-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r)
}

type issueTokenReq struct {
	BaseReq       rest.BaseReq   `json:"base_req"`
	Owner         sdk.AccAddress `json:"owner"` // owner of the token
	Symbol        string         `json:"symbol"`
	Name          string         `json:"name"`
	Scale         uint8          `json:"scale"`
	MinUnit       string         `json:"min_unit"`
	InitialSupply uint64         `json:"initial_supply"`
	MaxSupply     uint64         `json:"max_supply"`
	Mintable      bool           `json:"mintable"`
}

type editTokenReq struct {
	BaseReq   rest.BaseReq   `json:"base_req"`
	Owner     sdk.AccAddress `json:"owner"` //  owner of the token
	MaxSupply uint64         `json:"max_supply"`
	Mintable  string         `json:"mintable"` // mintable of the token
	Name      string         `json:"name"`
}

type transferTokenOwnerReq struct {
	BaseReq  rest.BaseReq   `json:"base_req"`
	SrcOwner sdk.AccAddress `json:"src_owner"` // the current owner address of the token
	DstOwner sdk.AccAddress `json:"dst_owner"` // the new owner
}

type mintTokenReq struct {
	BaseReq rest.BaseReq   `json:"base_req"`
	Owner   sdk.AccAddress `json:"owner"`  // the current owner address of the token
	To      sdk.AccAddress `json:"to"`     // address of minting token to
	Amount  uint64         `json:"amount"` // amount of minting token
}

package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"

	"github.com/irismod/token/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// issue a token
	r.HandleFunc(
		fmt.Sprintf("/%s/tokens", types.ModuleName),
		issueTokenHandlerFn(cliCtx),
	).Methods("POST")

	// edit a token
	r.HandleFunc(
		fmt.Sprintf("/%s/tokens/{%s}", types.ModuleName, RestParamSymbol),
		editTokenHandlerFn(cliCtx),
	).Methods("PUT")

	// transfer owner
	r.HandleFunc(
		fmt.Sprintf("/%s/tokens/{%s}/transfer", types.ModuleName, RestParamSymbol),
		transferOwnerHandlerFn(cliCtx),
	).Methods("POST")

	// mint token
	r.HandleFunc(
		fmt.Sprintf("/%s/tokens/{%s}/mint", types.ModuleName, RestParamSymbol),
		mintTokenHandlerFn(cliCtx),
	).Methods("POST")
}

func issueTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req issueTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgIssueToken message
		msg := types.MsgIssueToken{
			Symbol:        req.Symbol,
			Name:          req.Name,
			Scale:         req.Scale,
			MinUnit:       req.MinUnit,
			InitialSupply: req.InitialSupply,
			MaxSupply:     req.MaxSupply,
			Mintable:      req.Mintable,
			Owner:         req.Owner,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		authclient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func editTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars[RestParamSymbol]

		var req editTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		mintable, err := types.ParseBool(req.Mintable)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the MsgEditToken message
		msg := types.NewMsgEditToken(req.Name, symbol, req.MaxSupply, mintable, req.Owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		authclient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func transferOwnerHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars[RestParamSymbol]

		var req transferTokenOwnerReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgTransferTokenOwner message
		msg := types.NewMsgTransferTokenOwner(req.SrcOwner, req.DstOwner, symbol)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		authclient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func mintTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars[RestParamSymbol]

		var req mintTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgMintToken message
		msg := types.NewMsgMintToken(symbol, req.Owner, req.To, req.Amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		authclient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

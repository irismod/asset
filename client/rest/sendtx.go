package lcd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	context "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github/irismod/token/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// issue a token
	r.HandleFunc(
		fmt.Sprintf("/%s", types.ModuleName),
		issueTokenHandlerFn(cliCtx),
	).Methods("POST")

	// edit a token
	r.HandleFunc(
		fmt.Sprintf("/%s/{%s}", types.ModuleName, RestParamSymbol),
		editTokenHandlerFn(cliCtx),
	).Methods("PUT")

	// transfer owner
	r.HandleFunc(
		fmt.Sprintf("/%s/{%s}/transfer", types.ModuleName, RestParamSymbol),
		transferOwnerHandlerFn(cliCtx),
	).Methods("POST")

	// mint token
	r.HandleFunc(
		fmt.Sprintf("/%s/{%s}/mint", types.ModuleName, RestParamSymbol),
		mintTokenHandlerFn(cliCtx),
	).Methods("POST")
}

func issueTokenHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req issueTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseTx.Sanitize()
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

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
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

		baseReq := req.BaseTx.Sanitize()
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

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
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

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgTransferTokenOwner message
		msg := types.NewMsgTransferTokenOwner(req.SrcOwner, req.DstOwner, symbol)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
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

		baseReq := req.BaseTx.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgMintToken message
		msg := types.NewMsgMintToken(symbol, req.Owner, req.To, req.Amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseTx, []sdk.Msg{msg})
	}
}

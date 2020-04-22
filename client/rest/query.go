package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/irismod/token/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	// Query token by symbol or minUnit
	r.HandleFunc(
		fmt.Sprintf("/%s/tokens/{%s}", types.ModuleName, RestParamDenom),
		queryTokenHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	// Query tokens by owner
	r.HandleFunc(
		fmt.Sprintf("/%s/tokens", types.ModuleName),
		queryTokensHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	// Query token fees
	r.HandleFunc(
		fmt.Sprintf("/%s/tokens/{%s}/fee", types.ModuleName, RestParamSymbol),
		queryTokenFeesHandlerFn(cliCtx, queryRoute),
	).Methods("GET")
}

// queryTokenHandlerFn is the HTTP request handler to query token
func queryTokenHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		params := types.QueryTokenParams{
			Denom: vars[RestParamDenom],
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryToken), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// queryTokensHandlerFn is the HTTP request handler to query tokens
func queryTokensHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerStr := r.FormValue(RestParamOwner)

		var err error
		var owner sdk.AccAddress

		if len(ownerStr) > 0 {
			owner, err = sdk.AccAddressFromBech32(ownerStr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := types.QueryTokensParams{
			Owner: owner,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryTokens), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// queryTokenFeesHandlerFn is the HTTP request handler to query token fees
func queryTokenFeesHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars[RestParamSymbol]

		if err := types.CheckSymbol(symbol); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.QueryTokenFeesParams{
			Symbol: symbol,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s/tokens", queryRoute, types.QueryFees), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

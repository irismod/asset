package cli

import (
	"fmt"

	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irismod/token/exported"
	"github.com/irismod/token/types"
)

// GetQueryCmd returns the query commands for the token module.
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the token module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(flags.GetCommands(
		getCmdQueryToken(queryRoute, cdc),
		getCmdQueryTokens(queryRoute, cdc),
		getCmdQueryFee(queryRoute, cdc),
		getCmdQueryParams(queryRoute, cdc),
	)...)

	return queryCmd
}

// getCmdQueryToken implements the query token command.
func getCmdQueryToken(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "token [denom]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a token by symbol or minUnit.
Example:
$ %s query token token <denom>
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			if err := types.CheckSymbol(args[0]); err != nil {
				return err
			}

			params := types.QueryTokenParams{
				Denom: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryToken), bz)
			if err != nil {
				return err
			}

			var tokens types.Token
			if err := cdc.UnmarshalJSON(res, &tokens); err != nil {
				return err
			}

			return cliCtx.PrintOutput(tokens)
		},
	}

	return cmd
}

// getCmdQueryTokens implements the query tokens command.
func getCmdQueryTokens(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "tokens [owner]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query token by the owner.
Example:
$ %s query token tokens <owner>
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var err error
			var owner sdk.AccAddress

			if len(args) > 0 {
				owner, err = sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return err
				}
			}

			params := types.QueryTokensParams{
				Owner: owner,
			}

			bz := cdc.MustMarshalJSON(params)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryTokens), bz)
			if err != nil {
				return err
			}

			var tokens []exported.TokenI
			if err := cdc.UnmarshalJSON(res, &tokens); err != nil {
				return err
			}

			return cliCtx.PrintOutput(tokens)
		},
	}

	return cmd
}

// getCmdQueryFee implements the query token related fees command.
func getCmdQueryFee(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "fee [symbol]",
		Args: cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the token related fees.
Example:
$ %s query token fee <symbol>
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			symbol := args[0]
			if err := types.CheckSymbol(symbol); err != nil {
				return err
			}

			// query token fees
			fees, err := queryTokenFees(cliCtx, queryRoute, symbol)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(fees)
		},
	}

	return cmd
}

// getCmdQueryParams implements the query token related param command.
func getCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "params",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as token parameters.
Example:
$ %s query token params
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParams), nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := cdc.UnmarshalJSON(res, &params); err != nil {
				return err
			}

			return cliCtx.PrintOutput(params)
		},
	}

	return cmd
}

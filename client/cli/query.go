package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/irismod/token/types"
)

// GetQueryCmd returns the query commands for the token module.
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the token module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		getCmdQueryToken(cdc),
		getCmdQueryTokens(cdc),
		getCmdQueryFee(cdc),
		getCmdQueryParams(cdc),
	)

	return queryCmd
}

// getCmdQueryToken implements the query token command.
func getCmdQueryToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "token [denom]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a token by symbol or minUnit.
Example:
$ %s query token token <denom>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

			if err := types.CheckSymbol(args[0]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Token(context.Background(), &types.QueryTokenRequest{
				Denom: args[0],
			})

			if err != nil {
				return err
			}

			var token types.TokenI
			err = clientCtx.InterfaceRegistry.UnpackAny(res.Token, &token)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(token)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// getCmdQueryTokens implements the query tokens command.
func getCmdQueryTokens(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "tokens [owner]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query token by the owner.
Example:
$ %s query token tokens <owner>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

			var err error
			var owner sdk.AccAddress

			if len(args) > 0 {
				owner, err = sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return err
				}
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Tokens(context.Background(), &types.QueryTokensRequest{
				Owner: owner,
			})

			if err != nil {
				return err
			}

			tokens := make([]types.TokenI, 0, len(res.Tokens))
			for _, eviAny := range res.Tokens {
				var evi types.TokenI
				err = clientCtx.InterfaceRegistry.UnpackAny(eviAny, &evi)
				if err != nil {
					return err
				}
				tokens = append(tokens, evi)
			}

			return clientCtx.PrintOutput(tokens)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// getCmdQueryFee implements the query token related fees command.
func getCmdQueryFee(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "fee [symbol]",
		Args: cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the token related fees.
Example:
$ %s query token fee <symbol>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

			symbol := args[0]
			if err := types.CheckSymbol(symbol); err != nil {
				return err
			}

			// query token fees
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Fees(context.Background(), &types.QueryFeesRequest{
				Symbol: symbol,
			})

			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// getCmdQueryParams implements the query token related param command.
func getCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "params",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as token parameters.
Example:
$ %s query token params
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})

			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.Params)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

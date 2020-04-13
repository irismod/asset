package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github/irismod/token/internal/types"
)

// GetTxCmd returns the transaction commands for the token module.
func GetTxCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Asset transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(flags.PostCommands(
		getCmdIssueToken(queryRoute, cdc),
		getCmdEditToken(cdc),
		getCmdMintToken(queryRoute, cdc),
		getCmdTransferTokenOwner(cdc),
	)...)

	return txCmd
}

// getCmdIssueToken implements the issue token command
func getCmdIssueToken(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "issue",
		Short:   "Issue a new token",
		Example: `token issue --name="Kitty Token" --symbol="kitty" --min-unit="kitty" --scale=0 --initial-supply=100000000000 --max-supply=1000000000000 --mintable=true --from=<key-name> --chain-id=<chain-id> --fee=<fee>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			owner := cliCtx.GetFromAddress()

			msg := types.MsgIssueToken{
				Symbol:        viper.GetString(FlagSymbol),
				Name:          viper.GetString(FlagName),
				MinUnit:       viper.GetString(FlagMinUnit),
				Scale:         uint8(viper.GetInt(FlagScale)),
				InitialSupply: uint64(viper.GetInt(FlagInitialSupply)),
				MaxSupply:     uint64(viper.GetInt(FlagMaxSupply)),
				Mintable:      viper.GetBool(FlagMintable),
				Owner:         owner,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = "The token issue transaction will consume extra fee"

			if !viper.GetBool(flags.FlagGenerateOnly) {
				// query fee
				fee, err1 := queryTokenFees(cliCtx, queryRoute, msg.Symbol)
				if err1 != nil {
					return fmt.Errorf("failed to query token issue fee: %s", err1.Error())
				}

				// append issue fee to prompt
				issueFeeMainUnit := sdk.Coins{fee.IssueFee}.String()
				prompt += fmt.Sprintf(": %s", issueFeeMainUnit)
			}

			// a confirmation is needed
			prompt += "\nAre you sure to proceed?"
			confirmed, err := input.GetConfirmation(prompt, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			if !confirmed {
				return fmt.Errorf("operation aborted")
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsIssueToken)
	_ = cmd.MarkFlagRequired(FlagSymbol)
	_ = cmd.MarkFlagRequired(FlagName)
	_ = cmd.MarkFlagRequired(FlagInitialSupply)
	_ = cmd.MarkFlagRequired(FlagScale)

	return cmd
}

// getCmdEditToken implements the edit token command
func getCmdEditToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit [symbol]",
		Short:   "Edit an existing token",
		Example: `token edit <symbol> --name="Cat Token" --max-supply=100000000000 --mintable=true --from=<key-name> --chain-id=<chain-id> --fee=<fee>`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			owner := cliCtx.GetFromAddress()

			name := viper.GetString(FlagName)
			maxSupply := uint64(viper.GetInt(FlagMaxSupply))

			mintable, err := types.ParseBool(viper.GetString(FlagMintable))
			if err != nil {
				return err
			}

			msg := types.NewMsgEditToken(name, args[0], maxSupply, mintable, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsEditToken)
	return cmd
}

func getCmdMintToken(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mint [symbol]",
		Short:   "Mint tokens to a specified address",
		Example: `token mint <symbol> --amount=<amount> --to=<to> --from=<key-name> --chain-id=<chain-id>  --fee=<fee>`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			owner := cliCtx.GetFromAddress()

			amount := uint64(viper.GetInt64(FlagAmount))

			var to sdk.AccAddress
			var err error
			addr := viper.GetString(FlagTo)
			if len(strings.TrimSpace(addr)) > 0 {
				to, err = sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgMintToken(
				args[0], owner, to, amount,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = "The token mint transaction will consume extra fee"

			if !viper.GetBool(flags.FlagGenerateOnly) {
				// query fee
				fee, err1 := queryTokenFees(cliCtx, queryRoute, args[0])
				if err1 != nil {
					return fmt.Errorf("failed to query token mint fee: %s", err1.Error())
				}

				// append mint fee to prompt
				mintFeeMainUnit := sdk.Coins{fee.MintFee}.String()
				prompt += fmt.Sprintf(": %s", mintFeeMainUnit)
			}

			// a confirmation is needed
			prompt += "\nAre you sure to proceed?"
			confirmed, err := input.GetConfirmation(prompt, bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}

			if !confirmed {
				return fmt.Errorf("operation aborted")
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsMintToken)
	_ = cmd.MarkFlagRequired(FlagAmount)

	return cmd
}

// getCmdTransferTokenOwner implements the transfer token owner command
func getCmdTransferTokenOwner(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transfer [symbol]",
		Short:   "Transfer the owner of a token to a new owner",
		Example: `token transfer <symbol> --to=<to> --from=<key-name> --chain-id=<chain-id> --fee=<fee>`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			owner := cliCtx.GetFromAddress()

			to, err := sdk.AccAddressFromBech32(viper.GetString(FlagTo))
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferTokenOwner(owner, to, args[0])

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsTransferTokenOwner)
	_ = cmd.MarkFlagRequired(FlagTo)

	return cmd
}

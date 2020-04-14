package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/irismod/token/types"
)

// queryTokenFees retrieves the fees of issuance and minting for the specified symbol
func queryTokenFees(cliCtx context.CLIContext, queryRoute string, symbol string) (types.TokenFees, error) {
	params := types.QueryTokenFeesParams{
		Symbol: symbol,
	}

	bz := cliCtx.Codec.MustMarshalJSON(params)

	route := fmt.Sprintf("custom/%s/fees/tokens", queryRoute)
	res, _, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return types.TokenFees{}, err
	}

	var out types.TokenFees
	err = cliCtx.Codec.UnmarshalJSON(res, &out)
	return out, err
}

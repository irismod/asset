package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github/irismod/asset/internal/types"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryToken:
			return queryToken(ctx, req, k)
		case types.QueryTokens:
			return queryTokens(ctx, req, k)
		case types.QueryFees:
			return queryFees(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown asset query endpoint")
		}
	}
}

func queryToken(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryTokenParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, err
	}

	if err := types.CheckSymbol(params.Symbol); err != nil {
		return nil, err
	}

	token, err := queryTokenBySymbol(ctx, keeper, strings.ToLower(params.Symbol))
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, token)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

func queryTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryTokensParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, err
	}

	tokens, err := queryTokensByOwner(ctx, keeper, params.Owner)
	if err != nil {
		return nil, err
	}

	bz, er := codec.MarshalJSONIndent(keeper.cdc, tokens)
	if er != nil {
		return nil, err
	}

	return bz, nil
}

func queryFees(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryTokenFeesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, err
	}

	if err := types.CheckSymbol(params.Symbol); err != nil {
		return nil, err
	}

	symbol := strings.ToLower(params.Symbol)
	issueFee := keeper.GetTokenIssueFee(ctx, symbol)
	mintFee := keeper.GetTokenMintFee(ctx, symbol)

	fees := types.TokenFees{
		Exist:    keeper.HasToken(ctx, symbol),
		IssueFee: issueFee,
		MintFee:  mintFee,
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, fees)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

func queryTokenBySymbol(ctx sdk.Context, keeper Keeper, symbol string) (types.FungibleToken, error) {
	token, err := keeper.getToken(ctx, symbol)
	if err != nil {
		return types.FungibleToken{}, err
	}

	return token, nil
}

func queryTokensByOwner(ctx sdk.Context, keeper Keeper, owner sdk.AccAddress) (tokens types.Tokens, err error) {
	if len(owner) == 0 {
		keeper.IterateTokens(ctx, func(token types.FungibleToken) (stop bool) {
			tokens = append(tokens, token)
			return false
		})
		return
	}

	keeper.iterateTokensWithOwner(ctx, owner, func(token types.FungibleToken) (stop bool) {
		tokens = append(tokens, token)
		return false
	})

	return
}

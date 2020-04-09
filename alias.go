package asset

import (
	"github/irismod/asset/internal/keeper"
	"github/irismod/asset/internal/types"
)

type (
	MsgIssueToken         = types.MsgIssueToken
	MsgEditToken          = types.MsgEditToken
	MsgMintToken          = types.MsgMintToken
	MsgTransferTokenOwner = types.MsgTransferTokenOwner
	Tokens                = types.Tokens
	Params                = types.Params
	FungibleToken         = types.FungibleToken
	QueryTokenParams      = types.QueryTokenParams
	QueryTokensParams     = types.QueryTokensParams
	QueryTokenFeesParams  = types.QueryTokenFeesParams
	TokenFeesOutput       = types.TokenFees
	GenesisState          = types.GenesisState

	Keeper = keeper.Keeper
)

const (
	ModuleName            = types.ModuleName
	StoreKey              = types.StoreKey
	QuerierRoute          = types.QuerierRoute
	RouterKey             = types.RouterKey
	DefaultParamspace     = types.DefaultParamspace
	MaximumAssetMaxSupply = types.MaximumAssetMaxSupply
)

var (
	ModuleCdc = types.ModuleCdc
	//QuerierRoute             = types.QuerierRoute
	RegisterCodec = types.RegisterCodec
	CheckSymbol   = types.CheckSymbol
	ParseBool     = types.ParseBool

	NewFungibleToken             = types.NewFungibleToken
	NewMsgEditToken              = types.NewMsgEditToken
	NewMsgMintToken              = types.NewMsgMintToken
	NewMsgTransferTokenOwner     = types.NewMsgTransferTokenOwner
	DefaultParams                = types.DefaultParams
	ValidateParams               = types.ValidateParams
	NewGenesisState              = types.NewGenesisState
	NewParams                    = types.NewParams
	NewValidateTokenFeeDecorator = keeper.NewValidateTokenFeeDecorator

	QueryToken  = types.QueryToken
	QueryTokens = types.QueryTokens
	QueryFees   = types.QueryFees
	NewKeeper   = keeper.NewKeeper
	NewQuerier  = keeper.NewQuerier
)

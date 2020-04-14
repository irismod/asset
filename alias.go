package token

import (
	"github/irismod/token/keeper"
	"github/irismod/token/types"
)

type (
	MsgIssueToken         = types.MsgIssueToken
	MsgEditToken          = types.MsgEditToken
	MsgMintToken          = types.MsgMintToken
	MsgTransferTokenOwner = types.MsgTransferTokenOwner
	Tokens                = types.Tokens
	Token                 = types.Token
	Params                = types.Params
	QueryTokenParams      = types.QueryTokenParams
	QueryTokensParams     = types.QueryTokensParams
	QueryTokenFeesParams  = types.QueryTokenFeesParams
	TokenFees             = types.TokenFees
	GenesisState          = types.GenesisState

	Keeper = keeper.Keeper
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	RouterKey         = types.RouterKey
	DefaultParamspace = types.DefaultParamspace
	MaximumMaxSupply  = types.MaximumMaxSupply
)

var (
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
	CheckSymbol   = types.CheckSymbol
	ParseBool     = types.ParseBool

	NewToken                     = types.NewToken
	NewMsgEditToken              = types.NewMsgEditToken
	NewMsgMintToken              = types.NewMsgMintToken
	NewMsgTransferTokenOwner     = types.NewMsgTransferTokenOwner
	DefaultParams                = types.DefaultParams
	SetNativeToken               = types.SetNativeToken
	GetNativeToken               = types.GetNativeToken
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

package asset

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
	abci "github.com/tendermint/tendermint/abci/types"

	"github/irismod/asset/client/cli"
	"github/irismod/asset/simulation"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the asset module.
type AppModuleBasic struct{}

// Name returns the asset module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the asset module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the asset module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the asset module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	if err := ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the asset module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	//rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns no root tx command for the asset module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(StoreKey, cdc)
}

// GetQueryCmd returns the root query command for the asset module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

//____________________________________________________________________________

// AppModule implements an application module for the asset module.
type AppModule struct {
	AppModuleBasic
	assetKeeper   Keeper
	accountKeeper auth.AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, ak auth.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		assetKeeper:    keeper,
		accountKeeper:  ak,
	}
}

// Name returns the asset module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the asset module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the asset module.
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the asset module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.assetKeeper)
}

// QuerierRoute returns the asset module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the asset module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.assetKeeper)
}

// InitGenesis performs genesis initialization for the asset module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.assetKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the asset module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.assetKeeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the asset module.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the asset module. It returns no validator updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ProposalContents(simState module.SimulationState) []sim.WeightedProposalContent {
	return nil
}

func (am AppModule) RandomizedParams(r *rand.Rand) []sim.ParamChange {
	return simulation.ParamChanges(r)
}

func (am AppModule) WeightedOperations(simState module.SimulationState) []sim.WeightedOperation {
	return simulation.WeightedOperations(simState.AppParams, simState.Cdc, am.assetKeeper, am.accountKeeper)
}

// RegisterStoreDecoder registers a decoder for asset module's types.
func (AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[StoreKey] = simulation.DecodeStore
}

// GenerateGenesisState creates a randomized GenState of the asset module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

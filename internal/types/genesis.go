package types

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params"` // asset params
	Tokens Tokens `json:"tokens"` // issued tokens
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, tokens Tokens) GenesisState {
	return GenesisState{
		Params: params,
		Tokens: tokens,
	}
}

var DefaultToken = FungibleToken{
	Symbol:        "iris",
	Name:          "IRIS Network",
	Scale:         18,
	MinUnit:       "atto",
	InitialSupply: 2000000000,
	MaxSupply:     10000000000,
	Mintable:      true,
}

package kudos

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	kudossimulation "kudoschain/x/kudos/simulation"
	"kudoschain/x/kudos/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	kudosGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&kudosGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgGiveKudos          = "op_weight_msg_kudos"
		defaultWeightMsgGiveKudos int = 100
	)

	var weightMsgGiveKudos int
	simState.AppParams.GetOrGenerate(opWeightMsgGiveKudos, &weightMsgGiveKudos, nil,
		func(_ *rand.Rand) {
			weightMsgGiveKudos = defaultWeightMsgGiveKudos
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgGiveKudos,
		kudossimulation.SimulateMsgGiveKudos(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}

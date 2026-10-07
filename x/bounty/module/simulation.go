package bounty

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	bountysimulation "kudoschain/x/bounty/simulation"
	"kudoschain/x/bounty/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	bountyGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&bountyGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateBounty          = "op_weight_msg_bounty"
		defaultWeightMsgCreateBounty int = 100
	)

	var weightMsgCreateBounty int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateBounty, &weightMsgCreateBounty, nil,
		func(_ *rand.Rand) {
			weightMsgCreateBounty = defaultWeightMsgCreateBounty
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateBounty,
		bountysimulation.SimulateMsgCreateBounty(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgClaimBounty          = "op_weight_msg_bounty"
		defaultWeightMsgClaimBounty int = 100
	)

	var weightMsgClaimBounty int
	simState.AppParams.GetOrGenerate(opWeightMsgClaimBounty, &weightMsgClaimBounty, nil,
		func(_ *rand.Rand) {
			weightMsgClaimBounty = defaultWeightMsgClaimBounty
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClaimBounty,
		bountysimulation.SimulateMsgClaimBounty(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}

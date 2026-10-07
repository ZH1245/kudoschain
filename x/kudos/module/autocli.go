package kudos

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"kudoschain/x/kudos/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod:      "ShowKudo",
					Use:            "show-kudo [id]",
					Short:          "Query show-kudo",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "GiveKudos",
					Use:            "give-kudos [receiver] [message]",
					Short:          "Send a give-kudos tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "receiver"}, {ProtoField: "message"}},
				},
			},
		},
	}
}

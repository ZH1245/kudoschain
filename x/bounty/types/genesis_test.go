package types_test

import (
	"testing"

	"kudoschain/x/bounty/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc:     "empty params are invalid",
			genState: &types.GenesisState{},
			valid:    false,
		},
		{
			desc:     "custom lifetime is valid",
			genState: &types.GenesisState{Params: types.NewParams(50)},
			valid:    true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

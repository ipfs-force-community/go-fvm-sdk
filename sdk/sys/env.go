package sys

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func isSimulatedEnv(ctx context.Context) (*simulated.Fsm, bool) {
	env, ok := ctx.Value(types.SimulatedEnvkey).(*simulated.Fsm) //nolint:govet
	return env, ok                                              //nolint:govet
}

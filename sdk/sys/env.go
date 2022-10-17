package sys

import (
	"context"
	"errors"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

var (
	ErrorEnvValid = errors.New("env is valid")
)

func tryGetSimulator(ctx context.Context) (*simulated.FvmSimulator, bool) {
	env, ok := ctx.Value(types.SimulatedEnvkey).(*simulated.FvmSimulator) //nolint:govet
	return env, ok                                                        //nolint:govet
}

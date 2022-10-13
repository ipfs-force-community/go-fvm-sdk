//go:build simulated
// +build simulated

package sys

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func VMContext(ctx context.Context) (*types.InvocationContext, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.VMContext()
	}
	return &types.InvocationContext{}, nil
}

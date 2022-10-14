//go:build simulate
// +build simulate

package sys

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func VMContext(ctx context.Context) (*types.InvocationContext, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.VMContext()
	}
	panic(ErrorEnvValid)
}

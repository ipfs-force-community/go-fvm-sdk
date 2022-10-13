//go:build simulated
// +build simulated

package sys

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func Send(ctx context.Context, to address.Address, method uint64, params uint32, value types.TokenAmount) (*types.Send, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.Send(to, method, params, value)
	}
	return &types.Send{}, nil
}

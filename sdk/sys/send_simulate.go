//go:build simulate
// +build simulate

package sys

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func Send(ctx context.Context, to address.Address, method abi.MethodNum, params uint32, value abi.TokenAmount) (*types.SendResult, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.Send(to, method, params, value)
	}
	panic(ErrorEnvValid)
}

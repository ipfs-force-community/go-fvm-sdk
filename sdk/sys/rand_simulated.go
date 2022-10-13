//go:build simulated
// +build simulated

package sys

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"
)

func GetChainRandomness(ctx context.Context, dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.GetChainRandomness(dst, round, entropy)
	}
	return abi.Randomness{}, nil

}

func GetBeaconRandomness(ctx context.Context, dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.GetBeaconRandomness(dst, round, entropy)
	}
	return abi.Randomness{}, nil

}

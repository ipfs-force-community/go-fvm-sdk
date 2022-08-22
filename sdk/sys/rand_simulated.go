//go:build simulate
// +build simulate

package sys

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/internal/simulated"
)

func GetChainRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	return simulated.SimulatedInstance.GetChainRandomness(dst, round, entropy)
}

func GetBeaconRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	return simulated.SimulatedInstance.GetBeaconRandomness(dst, round, entropy)
}

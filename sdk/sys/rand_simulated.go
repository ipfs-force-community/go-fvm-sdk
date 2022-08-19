//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/fvm"
)

func GetChainRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	return fvm.MockFvmInstance.GetChainRandomness(dst, round, entropy)
}

func GetBeaconRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	return fvm.MockFvmInstance.GetBeaconRandomness(dst, round, entropy)
}

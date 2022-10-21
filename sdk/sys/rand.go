//go:build !simulate
// +build !simulate

package sys

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

func GetChainRandomness(_ context.Context, dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	result := [abi.RandomnessLength]byte{}
	resultPtr, _ := GetSlicePointerAndLen(result[:])

	entropyPtr, entropyLen := GetSlicePointerAndLen(entropy[:])
	code := getChainRandomness(resultPtr, dst, round, entropyPtr, entropyLen)
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get chain randomness")
	}

	return result[:], nil
}

func GetBeaconRandomness(_ context.Context, dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	result := [abi.RandomnessLength]byte{}
	resultPtr, _ := GetSlicePointerAndLen(result[:])

	entropyPtr, entropyLen := GetSlicePointerAndLen(entropy[:])
	code := getBeaconRandomness(resultPtr, dst, round, entropyPtr, entropyLen)
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get beacon randomness")
	}

	return result[:], nil
}

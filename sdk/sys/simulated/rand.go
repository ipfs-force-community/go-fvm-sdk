package simulated

import "github.com/filecoin-project/go-state-types/abi"

func (fvmSimulator *FvmSimulator) GetChainRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	return makeRandomness(dst, round, entropy), nil
}

func (fvmSimulator *FvmSimulator) GetBeaconRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	return makeRandomness(dst, round, entropy), nil
}

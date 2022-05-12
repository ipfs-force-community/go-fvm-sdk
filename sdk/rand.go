package sdk

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

func GetChainRandomness(dst crypto.DomainSeparationTag, round abi.ChainEpoch, entropy []byte) (abi.Randomness, error) {
	return sys.GetChainRandomness(int64(dst), int64(round), entropy)
}

func GetBeaconRandomness(dst crypto.DomainSeparationTag, round abi.ChainEpoch, entropy []byte) (abi.Randomness, error) {
	return sys.GetBeaconRandomness(int64(dst), int64(round), entropy)
}

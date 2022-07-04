package sdk

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

// GetChainRandomness gets 32 bytes of randomness from the ticket chain.
// The supplied output buffer must have at least 32 bytes of capacity.
// If this syscall succeeds, exactly 32 bytes will be written starting at the
// supplied offset.
func GetChainRandomness(dst crypto.DomainSeparationTag, round abi.ChainEpoch, entropy []byte) (abi.Randomness, error) {
	return sys.GetChainRandomness(int64(dst), int64(round), entropy)
}

// GetBeaconRandomness gets 32 bytes of randomness from the beacon system (currently Drand).
// The supplied output buffer must have at least 32 bytes of capacity.
// If this syscall succeeds, exactly 32 bytes will be written starting at the
// supplied offset.
func GetBeaconRandomness(dst crypto.DomainSeparationTag, round abi.ChainEpoch, entropy []byte) (abi.Randomness, error) {
	return sys.GetBeaconRandomness(int64(dst), int64(round), entropy)
}

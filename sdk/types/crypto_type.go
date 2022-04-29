package types

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
)

type AggregateSealVerifyInfo struct {
	Number                abi.SectorNumber
	Randomness            abi.SealRandomness
	InteractiveRandomness abi.InteractiveSealRandomness

	// Safe because we get those from the miner actor
	SealedCID   cid.Cid `checked:"true"` // CommR
	UnsealedCID cid.Cid `checked:"true"` // CommD
}

type AggregateSealVerifyProofAndInfos struct {
	Miner          abi.ActorID
	SealProof      abi.RegisteredSealProof
	AggregateProof abi.RegisteredAggregationProof
	Proof          []byte
	Infos          []AggregateSealVerifyInfo
}

type ReplicaUpdateInfo struct {
	UpdateProofType      abi.RegisteredUpdateProof
	OldSealedSectorCID   cid.Cid
	NewSealedSectorCID   cid.Cid
	NewUnsealedSectorCID cid.Cid
	Proof                []byte
}

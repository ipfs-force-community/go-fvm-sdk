package proof

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
)

type PoStProof struct {
	PoStProof  abi.RegisteredPoStProof
	ProofBytes []byte
}

type SectorInfo struct {
	SealProof    abi.RegisteredSealProof // RegisteredProof used when sealing - needs to be mapped to PoSt registered proof when used to verify a PoSt
	SectorNumber abi.SectorNumber
	SealedCID    cid.Cid // CommR
}

type ExtendedSectorInfo struct {
	SealProof    abi.RegisteredSealProof // RegisteredProof used when sealing - needs to be mapped to PoSt registered proof when used to verify a PoSt
	SectorNumber abi.SectorNumber
	SectorKey    *cid.Cid
	SealedCID    cid.Cid // CommR
}

type WinningPoStVerifyInfo struct {
	Randomness        abi.PoStRandomness
	Proofs            []PoStProof
	ChallengedSectors []SectorInfo
	Prover            abi.ActorID // used to derive 32-byte prover ID
}

// Information needed to verify a Window PoSt submitted directly to a miner actor.
type WindowPoStVerifyInfo struct {
	Randomness        abi.PoStRandomness
	Proofs            []PoStProof
	ChallengedSectors []SectorInfo
	Prover            abi.ActorID // used to derive 32-byte prover ID
}

type SealVerifyInfo struct {
	SealProof abi.RegisteredSealProof
	abi.SectorID
	DealIDs               []abi.DealID
	Randomness            abi.SealRandomness
	InteractiveRandomness abi.InteractiveSealRandomness
	Proof                 []byte

	// Safe because we get those from the miner actor
	SealedCID   cid.Cid `checked:"true"` // CommR
	UnsealedCID cid.Cid `checked:"true"` // CommD
}

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

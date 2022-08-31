//go:build simulate
// +build simulate

package sys

import (
	address "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/specs-actors/v7/actors/runtime"
	"github.com/filecoin-project/specs-actors/v7/actors/runtime/proof"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func VerifySignature(
	signature *crypto.Signature,
	signer *address.Address,
	plaintext []byte,
) (bool, error) {
	return simulated.DefaultFsm.VerifySignature(signature, signer, plaintext)
}

func HashBlake2b(data []byte) ([32]byte, error) {
	return simulated.DefaultFsm.HashBlake2b(data)
}

func ComputeUnsealedSectorCid(
	proofType abi.RegisteredSealProof,
	pieces []abi.PieceInfo,
) (cid.Cid, error) {

	return simulated.DefaultFsm.ComputeUnsealedSectorCid(proofType, pieces)
}

// VerifySeal Verifies a sector seal proof.
func VerifySeal(info *proof.SealVerifyInfo) (bool, error) {
	return simulated.DefaultFsm.VerifySeal(info)
}

// VerifyPost Verifies a sector seal proof.
func VerifyPost(info *proof.WindowPoStVerifyInfo) (bool, error) {
	return simulated.DefaultFsm.VerifyPost(info)
}

func VerifyConsensusFault(
	h1 []byte,
	h2 []byte,
	extra []byte,
) (*runtime.ConsensusFault, error) {

	return simulated.DefaultFsm.VerifyConsensusFault(h1, h2, extra)
}

func VerifyAggregateSeals(info *types.AggregateSealVerifyProofAndInfos) (bool, error) {
	return simulated.DefaultFsm.VerifyAggregateSeals(info)
}

func VerifyReplicaUpdate(info *types.ReplicaUpdateInfo) (bool, error) {
	return simulated.DefaultFsm.VerifyReplicaUpdate(info)
}

func BatchVerifySeals(sealVerifyInfos []proof.SealVerifyInfo) ([]bool, error) {
	return simulated.DefaultFsm.BatchVerifySeals(sealVerifyInfos)
}
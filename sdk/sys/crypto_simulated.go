//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/fvm"
)

func VerifySignature(
	signature *crypto.Signature,
	signer *address.Address,
	plaintext []byte,
) (bool, error) {
	return fvm.MockFvmInstance.VerifySignature(signature, signer, plaintext)
}

func HashBlake2b(data []byte) ([32]byte, error) {
	return fvm.MockFvmInstance.HashBlake2b(data)
}

func ComputeUnsealedSectorCid(
	proofType abi.RegisteredSealProof,
	pieces []abi.PieceInfo,
) (cid.Cid, error) {

	return fvm.MockFvmInstance.ComputeUnsealedSectorCid(proofType, pieces)
}

// VerifySeal Verifies a sector seal proof.
func VerifySeal(info *proof.SealVerifyInfo) (bool, error) {
	return fvm.MockFvmInstance.VerifySeal(info)
}

// VerifyPost Verifies a sector seal proof.
func VerifyPost(info *proof.WindowPoStVerifyInfo) (bool, error) {
	return fvm.MockFvmInstance.VerifyPost(info)
}

func VerifyConsensusFault(
	h1 []byte,
	h2 []byte,
	extra []byte,
) (*runtime.ConsensusFault, error) {

	return fvm.MockFvmInstance.VerifyConsensusFault(h1, h2, extra)
}

func VerifyAggregateSeals(info *types.AggregateSealVerifyProofAndInfos) (bool, error) {
	return fvm.MockFvmInstance.VerifyAggregateSeals(info)
}

func VerifyReplicaUpdate(info *types.ReplicaUpdateInfo) (bool, error) {
	return fvm.MockFvmInstance.VerifyReplicaUpdate(info)
}

func BatchVerifySeals(sealVerifyInfos []proof.SealVerifyInfo) ([]bool, error) {
	return fvm.MockFvmInstance.BatchVerifySeals(sealVerifyInfos)
}

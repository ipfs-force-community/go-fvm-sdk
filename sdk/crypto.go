package sdk

import (
	address "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/specs-actors/v2/actors/runtime/proof"
	"github.com/filecoin-project/specs-actors/v7/actors/runtime"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func VerifySignature(
	signature *crypto.Signature,
	signer *address.Address,
	plainText []byte,
) (bool, error) {
	return sys.VerifySignature(signature, signer, plainText)
}

func HashBlake2b(data []byte) ([32]byte, error) {
	return sys.HashBlake2b(data)
}

func ComputeUnsealedSectorCid(
	proofType abi.RegisteredSealProof,
	pieces []abi.PieceInfo,
) (cid.Cid, error) {
	return sys.ComputeUnsealedSectorCid(proofType, pieces)
}

// VerifySeal verifies a sector seal proof.
func VerifySeal(info *proof.SealVerifyInfo) (bool, error) {
	return sys.VerifySeal(info)
}

//VerifyPost verifies a sector seal proof.
func VerifyPost(info *proof.WindowPoStVerifyInfo) (bool, error) {
	return sys.VerifyPost(info)
}

func VerifyConsensusFault(
	h1 []byte,
	h2 []byte,
	extra []byte,
) (*runtime.ConsensusFault, error) {
	return sys.VerifyConsensusFault(h1, h2, extra)
}

func VerifyAggregateSeals(info *types.AggregateSealVerifyProofAndInfos) (bool, error) {
	return sys.VerifyAggregateSeals(info)
}

func VerifyReplicaUpdate(info *types.ReplicaUpdateInfo) (bool, error) {
	return sys.VerifyReplicaUpdate(info)
}

func BatchVerifySeals(sealVerifyInfos []proof.SealVerifyInfo) ([]bool, error) {
	return sys.BatchVerifySeals(sealVerifyInfos)
}

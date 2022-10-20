package simulated

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/proof"
	"github.com/filecoin-project/specs-actors/actors/runtime"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func (s *FvmSimulator) VerifySignature(
	signature *crypto.Signature,
	signer *address.Address,
	plaintext []byte,
) (bool, error) {
	panic("This is not implement")
}

func (s *FvmSimulator) HashBlake2b(data []byte) ([32]byte, error) {
	result := blakehash(data)
	var temp [32]byte
	copy(temp[:], result[:32])
	return temp, nil
}

func (s *FvmSimulator) ComputeUnsealedSectorCid(
	proofType abi.RegisteredSealProof,
	pieces []abi.PieceInfo,
) (cid.Cid, error) {
	panic("This is not implement")
}

func (s *FvmSimulator) VerifySeal(info *proof.SealVerifyInfo) (bool, error) {
	panic("This is not implement")
}

func (s *FvmSimulator) VerifyPost(info *proof.WindowPoStVerifyInfo) (bool, error) {
	panic("This is not implement")
}

func (s *FvmSimulator) VerifyConsensusFault(h1 []byte, h2 []byte, extra []byte,
) (*runtime.ConsensusFault, error) {
	panic("This is not implement")
}

func (s *FvmSimulator) VerifyAggregateSeals(info *types.AggregateSealVerifyProofAndInfos) (bool, error) {
	panic("This is not implement")
}

func (s *FvmSimulator) VerifyReplicaUpdate(info *types.ReplicaUpdateInfo) (bool, error) {
	panic("This is not implement")
}
func (s *FvmSimulator) BatchVerifySeals(sealVerifyInfos []proof.SealVerifyInfo) ([]bool, error) {
	panic("This is not implement")
}

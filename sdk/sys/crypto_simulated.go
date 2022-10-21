//go:build simulate
// +build simulate

package sys

import (
	"context"

	address "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/proof"
	"github.com/filecoin-project/specs-actors/v7/actors/runtime"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func VerifySignature(
	ctx context.Context,
	signature *crypto.Signature,
	signer *address.Address,
	plaintext []byte,
) (bool, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.VerifySignature(signature, signer, plaintext)
	}
	panic(ErrorEnvValid)
}

func HashBlake2b(ctx context.Context, data []byte) ([32]byte, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.HashBlake2b(data)
	}
	panic(ErrorEnvValid)
}

func ComputeUnsealedSectorCid(
	ctx context.Context,
	proofType abi.RegisteredSealProof,
	pieces []abi.PieceInfo,
) (cid.Cid, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.ComputeUnsealedSectorCid(proofType, pieces)
	}
	panic(ErrorEnvValid)
}

// VerifySeal Verifies a sector seal proof.
func VerifySeal(ctx context.Context, info *proof.SealVerifyInfo) (bool, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.VerifySeal(info)
	}
	return false, nil
}

// VerifyPost Verifies a sector seal proof.
func VerifyPost(ctx context.Context, info *proof.WindowPoStVerifyInfo) (bool, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.VerifyPost(info)
	}
	panic(ErrorEnvValid)
}

func VerifyConsensusFault(
	ctx context.Context,
	h1 []byte,
	h2 []byte,
	extra []byte,
) (*runtime.ConsensusFault, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.VerifyConsensusFault(h1, h2, extra)
	}
	panic(ErrorEnvValid)
}

func VerifyAggregateSeals(ctx context.Context, info *types.AggregateSealVerifyProofAndInfos) (bool, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.VerifyAggregateSeals(info)
	}
	panic(ErrorEnvValid)
}

func VerifyReplicaUpdate(ctx context.Context, info *types.ReplicaUpdateInfo) (bool, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.VerifyReplicaUpdate(info)
	}
	panic(ErrorEnvValid)
}

func BatchVerifySeals(ctx context.Context, sealVerifyInfos []proof.SealVerifyInfo) ([]bool, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.BatchVerifySeals(sealVerifyInfos)
	}
	panic(ErrorEnvValid)
}

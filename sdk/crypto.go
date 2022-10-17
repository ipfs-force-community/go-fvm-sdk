package sdk

import (
	"context"

	address "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/specs-actors/v2/actors/runtime/proof"
	"github.com/filecoin-project/specs-actors/v7/actors/runtime"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

// VerifySignature verifies that a signature is valid for an address and plaintext.
func VerifySignature(
	ctx context.Context,
	signature *crypto.Signature,
	signer *address.Address,
	plainText []byte,
) (bool, error) {
	return sys.VerifySignature(ctx, signature, signer, plainText)
}

// HashBlake2b hashes input data using blake2b with 256 bit output.
func HashBlake2b(ctx context.Context, data []byte) ([32]byte, error) {
	return sys.HashBlake2b(ctx, data)
}

// ComputeUnsealedSectorCid computes an unsealed sector CID (CommD) from its constituent piece CIDs (CommPs) and sizes.
func ComputeUnsealedSectorCid(
	ctx context.Context,
	proofType abi.RegisteredSealProof,
	pieces []abi.PieceInfo,
) (cid.Cid, error) {
	return sys.ComputeUnsealedSectorCid(ctx, proofType, pieces)
}

// VerifySeal verifies a sector seal proof.
func VerifySeal(ctx context.Context, info *proof.SealVerifyInfo) (bool, error) {
	return sys.VerifySeal(ctx, info)
}

// VerifyPost verifies a sector seal proof.
func VerifyPost(ctx context.Context, info *proof.WindowPoStVerifyInfo) (bool, error) {
	return sys.VerifyPost(ctx, info)
}

// VerifyConsensusFault verifies that two block headers provide proof of a consensus fault:
// - both headers mined by the same actor
// - headers are different
// - first header is of the same or lower epoch as the second
// - at least one of the headers appears in the current chain at or after epoch `earliest`
// - the headers provide evidence of a fault (see the spec for the different fault types).
// The parameters are all serialized block headers. The third "extra" parameter is consulted only for
// the "parent grinding fault", in which case it must be the sibling of h1 (same parent tipset) and one of the
// blocks in the parent of h2 (i.e. h2's grandparent).
// Returns None and an error if the headers don't prove a fault.
func VerifyConsensusFault(
	ctx context.Context,
	h1 []byte,
	h2 []byte,
	extra []byte,
) (*runtime.ConsensusFault, error) {
	return sys.VerifyConsensusFault(ctx, h1, h2, extra)
}

// VerifyAggregateSeals verifies aggregate proof of replication of sectors
func VerifyAggregateSeals(ctx context.Context, info *types.AggregateSealVerifyProofAndInfos) (bool, error) {
	return sys.VerifyAggregateSeals(ctx, info)
}

// VerifyReplicaUpdate verifies sector replica update
func VerifyReplicaUpdate(ctx context.Context, info *types.ReplicaUpdateInfo) (bool, error) {
	return sys.VerifyReplicaUpdate(ctx, info)
}

// BatchVerifySeals batch verifies seals
func BatchVerifySeals(ctx context.Context, sealVerifyInfos []proof.SealVerifyInfo) ([]bool, error) {
	return sys.BatchVerifySeals(ctx, sealVerifyInfos)
}

//go:build !simulate
// +build !simulate

package sys

import (
	"bytes"
	"context"
	"fmt"
	"unsafe"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/proof"
	"github.com/filecoin-project/specs-actors/v7/actors/runtime"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func VerifySignature(
	_ context.Context,
	signature *crypto.Signature,
	signer *address.Address,
	plaintext []byte,
) (bool, error) {
	sigBuf := bytes.NewBuffer([]byte{})
	err := signature.MarshalCBOR(sigBuf)
	if err != nil {
		return false, fmt.Errorf("unable to marshal signature %w", err)
	}

	var result int32
	sigPtr, sigLen := GetSlicePointerAndLen(sigBuf.Bytes())
	signerPtr, signerLen := GetSlicePointerAndLen(signer.Bytes())
	plainTextPtr, plainTextLen := GetSlicePointerAndLen(plaintext)
	code := cryptoVerifySignature(uintptr(unsafe.Pointer(&result)), uint32(signature.Type), sigPtr, sigLen, signerPtr, signerLen, plainTextPtr, plainTextLen)
	if code != 0 {
		return false, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unable to verify signature")
	}
	return result == 0, nil
}

func HashBlake2b(_ context.Context, data []byte) ([32]byte, error) {
	dataPtr, dataLen := GetSlicePointerAndLen(data)
	digest := [32]byte{}
	digestPtr, digestLen := GetSlicePointerAndLen(digest)
	result := [32]byte{}
	resultPtr, _ := GetSlicePointerAndLen(result[:])
	code := cryptoHashBlake2b(resultPtr, 0xb220, dataPtr, dataLen, digestPtr, digestLen)
	if code != 0 {
		return result, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "failed to compute blake2b hash")
	}
	return result, nil
}

func ComputeUnsealedSectorCid(
	_ context.Context,
	proofType abi.RegisteredSealProof,
	pieces []abi.PieceInfo,
) (cid.Cid, error) {
	//todo need to be test
	buf := bytes.NewBuffer([]byte{})
	cw := cbg.NewCborWriter(buf)
	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(pieces))); err != nil {
		return cid.Undef, err
	}
	for _, piece := range pieces {
		if err := piece.MarshalCBOR(cw); err != nil {
			return cid.Undef, err
		}
	}
	piecesPtr, piecesLen := GetSlicePointerAndLen(buf.Bytes())
	cidBuf := make([]byte, types.MaxCidLen)
	cidBufPtr, _ := GetSlicePointerAndLen(cidBuf)
	var cidLen uint32
	code := cryptoComputeUnsealedSectorCid(uintptr(unsafe.Pointer(&cidLen)), int64(proofType), piecesPtr, piecesLen, cidBufPtr, types.MaxCidLen)
	if code != 0 {
		return cid.Undef, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unable to verify signature")
	}

	_, sID, err := cid.CidFromBytes(cidBuf[:cidLen])
	if err != nil {
		return cid.Undef, fmt.Errorf("unable to decode cid from compute unseal sector cid result, cid len %d, cid content %v %w", cidLen, cidBuf[:cidLen], err)
	}
	return sID, nil
}

// VerifySeal Verifies a sector seal proof.
func VerifySeal(_ context.Context, info *proof.SealVerifyInfo) (bool, error) {
	verifyBuf := bytes.NewBuffer([]byte{})
	err := info.MarshalCBOR(verifyBuf)
	if err != nil {
		return false, fmt.Errorf("unable to marshal signature %w", err)
	}
	var result int32
	verifyBufPtr, verifyBufLen := GetSlicePointerAndLen(verifyBuf.Bytes())
	code := cryptoVerifySeal(uintptr(unsafe.Pointer(&result)), verifyBufPtr, verifyBufLen)
	if code != 0 {
		return false, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unable to verify signature")
	}
	return result == 0, nil
}

// VerifyPost Verifies a sector seal proof.
func VerifyPost(_ context.Context, info *proof.WindowPoStVerifyInfo) (bool, error) {
	verifyBuf := bytes.NewBuffer([]byte{})
	err := info.MarshalCBOR(verifyBuf)
	if err != nil {
		return false, fmt.Errorf("unable to marshal signature %w", err)
	}
	var result int32
	verifyBufPtr, verifyBufLen := GetSlicePointerAndLen(verifyBuf.Bytes())
	code := cryptoVerifyPost(uintptr(unsafe.Pointer(&result)), verifyBufPtr, verifyBufLen)
	if code != 0 {
		return false, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unable to verify signature")
	}
	return result == 0, nil
}

func VerifyConsensusFault(
	_ context.Context,
	h1 []byte,
	h2 []byte,
	extra []byte,
) (*runtime.ConsensusFault, error) {
	h1Ptr, h1Len := GetSlicePointerAndLen(h1)
	h2Ptr, h2Len := GetSlicePointerAndLen(h2)
	extraPtr, extraLen := GetSlicePointerAndLen(extra)
	verifyFault := new(types.VerifyConsensusFault)
	code := cryptoVerifyConsensusFault(uintptr(unsafe.Pointer(verifyFault)), h1Ptr, h1Len, h2Ptr, h2Len, extraPtr, extraLen)
	if code != 0 {
		return nil, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unable to verify consensus signature")
	}
	if verifyFault.Fault == 0 {
		return nil, nil
	}

	faultType := runtime.ConsensusFaultType(verifyFault.Fault)
	if !types.ValidateConsensusFaultType(faultType) {
		return nil, ferrors.NewSysCallError(ferrors.ErrorNumber(code), fmt.Sprintf("received an invalid fault type (%d) from the runtime", faultType))
	}
	target, err := address.NewIDAddress(uint64(verifyFault.Target))
	if err != nil {
		return nil, ferrors.NewSysCallError(ferrors.AssertionFailed, fmt.Sprintf("unable to new id address for %d %v", verifyFault.Target, err))
	}
	return &runtime.ConsensusFault{
		Epoch:  abi.ChainEpoch(verifyFault.Epoch),
		Target: target,
		Type:   faultType,
	}, nil
}

func VerifyAggregateSeals(ctx context.Context, info *types.AggregateSealVerifyProofAndInfos) (bool, error) {
	aggregateSealBuf := bytes.NewBuffer([]byte{})
	err := info.MarshalCBOR(aggregateSealBuf)
	if err != nil {
		return false, fmt.Errorf("unable to marshal signature %w", err)
	}
	var result int32
	aggregateSealBufPtr, aggregateSealBufLen := GetSlicePointerAndLen(aggregateSealBuf.Bytes())
	code := cryptoVerifyAggregateSeals(uintptr(unsafe.Pointer(&result)), aggregateSealBufPtr, aggregateSealBufLen)
	if code != 0 {
		return false, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unable to verify aggregate seals")
	}
	return result == 0, nil
}

func VerifyReplicaUpdate(_ context.Context, info *types.ReplicaUpdateInfo) (bool, error) {
	replicaUpdateInfoBuf := bytes.NewBuffer([]byte{})
	err := info.MarshalCBOR(replicaUpdateInfoBuf)
	if err != nil {
		return false, fmt.Errorf("unable to marshal signature %w", err)
	}
	var result int32
	replicaUpdateInfoBufPtr, replicaUpdateInfoBufLen := GetSlicePointerAndLen(replicaUpdateInfoBuf.Bytes())
	code := cryptoVerifyReplicaUpdate(uintptr(unsafe.Pointer(&result)), replicaUpdateInfoBufPtr, replicaUpdateInfoBufLen)
	if code != 0 {
		return false, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unable to verify aggregate seals")
	}
	return result == 0, nil
}

func BatchVerifySeals(_ context.Context, sealVerifyInfos []proof.SealVerifyInfo) ([]bool, error) {
	//todo need to be test
	buf := bytes.NewBuffer([]byte{})
	cw := cbg.NewCborWriter(buf)
	batchCount := uint64(len(sealVerifyInfos))
	if err := cw.WriteMajorTypeHeader(cbg.MajArray, batchCount); err != nil {
		return nil, err
	}
	for _, sealVerifyInfo := range sealVerifyInfos {
		if err := sealVerifyInfo.MarshalCBOR(cw); err != nil {
			return nil, err
		}
	}
	sealInfoPtr, sealInfoLen := GetSlicePointerAndLen(buf.Bytes())
	verifyResult := make([]bool, batchCount)
	resultPtr, _ := GetSlicePointerAndLen(verifyResult)
	code := cryptoBatchVerifySeals(sealInfoPtr, sealInfoLen, resultPtr)
	if code != 0 {
		return nil, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unable to batch verify seal info")
	}
	return verifyResult, nil
}

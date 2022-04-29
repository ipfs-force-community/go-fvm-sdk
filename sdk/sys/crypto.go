package sys

import (
	"bytes"
	"fmt"
	"unsafe"

	address "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/specs-actors/v7/actors/runtime"
	"github.com/filecoin-project/specs-actors/v7/actors/runtime/proof"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
)

/// Verifies that a signature is valid for an address and plaintext.
///
/// Returns 0 on success, or -1 if the signature fails to validate.
///
/// # Errors
///
/// | Error             | Reason                                               |
/// |-------------------|------------------------------------------------------|
/// | `IllegalArgument` | signature, address, or plaintext buffers are invalid |
//go:wasm-module crypto
//export verify_signature
func cryptoVerifySignature(ret uintptr, sigOff uintptr, sigLen uint32, addrOff uintptr, addrLen uint32, plainTextOff uintptr, plainTextLen uint32) uint32

/// Hashes input data using blake2b with 256 bit output.
///
/// The output buffer must be sized to a minimum of 32 bytes.
///
/// # Errors
///
/// | Error             | Reason                                          |
/// |-------------------|-------------------------------------------------|
/// | `IllegalArgument` | the input buffer does not point to valid memory |
//go:wasm-module crypto
//export hash_blake2b
func cryptoHashBlake2b(ret uintptr, dataOff uintptr, dataLen uint32) uint32

/// Computes an unsealed sector CID (CommD) from its constituent piece CIDs
/// (CommPs) and sizes.
///
/// Writes the CID in the provided output buffer, and returns the length of
/// the written CID.
///
/// # Errors
///
/// | Error             | Reason                   |
/// |-------------------|--------------------------|
/// | `IllegalArgument` | an argument is malformed |
//go:wasm-module crypto
//export hash_blake2b
func cryptoComputeUnsealedSectorCid(ret uintptr, proofType int64, piecesOff uintptr, pieceLen uint32, cidPtr uintptr, cidLen uint32) uint32

/// Verifies a sector seal proof.
///
/// Returns 0 to indicate that the proof was valid, -1 otherwise.
///
/// # Errors
///
/// | Error             | Reason                   |
/// |-------------------|--------------------------|
/// | `IllegalArgument` | an argument is malformed |
//go:wasm-module crypto
//export verify_seal
func cryptoVerifySeal(ret uintptr, infoOff uintptr, infoLen uint32) uint32

/// Verifies a window proof of spacetime.
///
/// Returns 0 to indicate that the proof was valid, -1 otherwise.
///
/// # Errors
///
/// | Error             | Reason                   |
/// |-------------------|--------------------------|
/// | `IllegalArgument` | an argument is malformed |
//go:wasm-module crypto
//export verify_post
func cryptoVerifyPost(ret uintptr, infoOff uintptr, infoLen uint32) uint32

/// Verifies that two block headers provide proof of a consensus fault.
///
/// Returns a 0 status if a consensus fault was recognized, along with the
/// BlockId containing the fault details. Otherwise, a -1 status is returned,
/// and the second result parameter must be ignored.
///
/// # Errors
///
/// | Error             | Reason                                |
/// |-------------------|---------------------------------------|
/// | `LimitExceeded`   | exceeded lookback limit finding block |
/// | `IllegalArgument` | an argument is malformed              |
//go:wasm-module crypto
//export verify_consensus_fault
func cryptoVerifyConsensusFault(ret uintptr, h1Off uintptr, h1Len uint32, h2Off uintptr, h2Len uint32, extraOff uintptr, extraLen uint32) uint32

/// Verifies an aggregated batch of sector seal proofs.
///
/// Returns 0 to indicate that the proof was valid, -1 otherwise.
///
/// # Errors
///
/// | Error             | Reason                         |
/// |-------------------|--------------------------------|
/// | `LimitExceeded`   | exceeds seal aggregation limit |
/// | `IllegalArgument` | an argument is malformed       |
//go:wasm-module crypto
//export verify_aggregate_seals
func cryptoVerifyAggregateSeals(ret uintptr, aggOff uintptr, aggLen uint32) uint32

/// Verifies a replica update proof.
///
/// Returns 0 to indicate that the proof was valid, -1 otherwise.
///
/// # Errors
///
/// | Error             | Reason                        |
/// |-------------------|-------------------------------|
/// | `LimitExceeded`   | exceeds replica update limit  |
/// | `IllegalArgument` | an argument is malformed      |
//go:wasm-module crypto
//export verify_replica_update
func cryptoVerifyReplicaUpdate(ret uintptr, repOff uintptr, repLen uint32) uint32

/// Verifies an aggregated batch of sector seal proofs.
///
/// # Errors
///
/// | Error             | Reason                   |
/// |-------------------|--------------------------|
/// | `IllegalArgument` | an argument is malformed |
//go:wasm-module crypto
//export batch_verify_seals
func cryptoBatchVerifySeals(batchOff uintptr, batLen uint32, resultOff uintptr) uint32

func VerifySignature(
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
	code := cryptoVerifySignature(uintptr(unsafe.Pointer(&result)), sigPtr, sigLen, signerPtr, signerLen, plainTextPtr, plainTextLen)
	if code != 0 {
		return false, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to verify signature")
	}
	return result == 0, nil
}

func HashBlake2b(data []byte) ([32]byte, error) {
	dataPtr, dataLen := GetSlicePointerAndLen(data)
	result := [32]byte{}
	resultLen, _ := GetSlicePointerAndLen(result[:])
	code := cryptoHashBlake2b(resultLen, dataPtr, dataLen)
	if code != 0 {
		return result, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to compute blak2b hash")
	}
	return result, nil
}

func ComputeUnsealedSectorCid(
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
	cidBuf := make([]byte, types.MAX_CID_LEN)
	cidBufPtr, _ := GetSlicePointerAndLen(cidBuf)
	var cidLen uint32
	code := cryptoComputeUnsealedSectorCid(uintptr(unsafe.Pointer(&cidLen)), int64(proofType), piecesPtr, piecesLen, cidBufPtr, types.MAX_CID_LEN)
	if code != 0 {
		return cid.Undef, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to verify signature")
	}

	_, sId, err := cid.CidFromBytes(cidBuf[:cidLen])
	if err != nil {
		return cid.Undef, fmt.Errorf("unable to decode cid from compute unseal sector cid result, cid len %d, cid content %v %w", cidLen, cidBuf[:cidLen], err)
	}
	return sId, nil
}

/// Verifies a sector seal proof.
func VerifySeal(info *proof.SealVerifyInfo) (bool, error) {
	verifyBuf := bytes.NewBuffer([]byte{})
	err := info.MarshalCBOR(verifyBuf)
	if err != nil {
		return false, fmt.Errorf("unable to marshal signature %w", err)
	}
	var result int32
	verifyBufPtr, verifyBufLen := GetSlicePointerAndLen(verifyBuf.Bytes())
	code := cryptoVerifySeal(uintptr(unsafe.Pointer(&result)), verifyBufPtr, verifyBufLen)
	if code != 0 {
		return false, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to verify signature")
	}
	return result == 0, nil
}

/// Verifies a sector seal proof.
func VerifyPost(info *proof.WindowPoStVerifyInfo) (bool, error) {
	verifyBuf := bytes.NewBuffer([]byte{})
	err := info.MarshalCBOR(verifyBuf)
	if err != nil {
		return false, fmt.Errorf("unable to marshal signature %w", err)
	}
	var result int32
	verifyBufPtr, verifyBufLen := GetSlicePointerAndLen(verifyBuf.Bytes())
	code := cryptoVerifyPost(uintptr(unsafe.Pointer(&result)), verifyBufPtr, verifyBufLen)
	if code != 0 {
		return false, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to verify signature")
	}
	return result == 0, nil
}

func VerifyConsensusFault(
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
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to verify consensus signature")
	}
	if verifyFault.Fault == 0 {
		return nil, nil
	}

	faultType := runtime.ConsensusFaultType(verifyFault.Fault)
	if !types.ValidateConsensusFaultType(faultType) {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("received an invalid fault type (%d) from the runtime", faultType))
	}
	target, err := address.NewIDAddress(verifyFault.Target)
	if err != nil {
		return nil, ferrors.NewFvmError(ferrors.SYS_ASSERTION_FAILED, fmt.Sprintf("unable to new id address for %d %w", verifyFault.Target, err))
	}
	return &runtime.ConsensusFault{
		Epoch:  abi.ChainEpoch(verifyFault.Epoch),
		Target: target,
		Type:   faultType,
	}, nil
}

func VerifyAggregateSeals(info *types.AggregateSealVerifyProofAndInfos) (bool, error) {
	aggregateSealBuf := bytes.NewBuffer([]byte{})
	err := info.MarshalCBOR(aggregateSealBuf)
	if err != nil {
		return false, fmt.Errorf("unable to marshal signature %w", err)
	}
	var result int32
	aggregateSealBufPtr, aggregateSealBufLen := GetSlicePointerAndLen(aggregateSealBuf.Bytes())
	code := cryptoVerifyAggregateSeals(uintptr(unsafe.Pointer(&result)), aggregateSealBufPtr, aggregateSealBufLen)
	if code != 0 {
		return false, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to verify aggregate seals")
	}
	return result == 0, nil
}

func VerifyReplicaUpdate(info *types.ReplicaUpdateInfo) (bool, error) {
	replicaUpdateInfoBuf := bytes.NewBuffer([]byte{})
	err := info.MarshalCBOR(replicaUpdateInfoBuf)
	if err != nil {
		return false, fmt.Errorf("unable to marshal signature %w", err)
	}
	var result int32
	replicaUpdateInfoBufPtr, replicaUpdateInfoBufLen := GetSlicePointerAndLen(replicaUpdateInfoBuf.Bytes())
	code := cryptoVerifyReplicaUpdate(uintptr(unsafe.Pointer(&result)), replicaUpdateInfoBufPtr, replicaUpdateInfoBufLen)
	if code != 0 {
		return false, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to verify aggregate seals")
	}
	return result == 0, nil
}

func BatchVerifySeals(sealVerifyInfos []proof.SealVerifyInfo) ([]bool, error) {
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
	verifyResult := make([]bool, batchCount, batchCount)
	resultPtr, _ := GetSlicePointerAndLen(verifyResult)
	code := cryptoBatchVerifySeals(sealInfoPtr, sealInfoLen, resultPtr)
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to batch verify seal info")
	}
	return verifyResult, nil
}

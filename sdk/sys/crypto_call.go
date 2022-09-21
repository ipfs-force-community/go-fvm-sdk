//go:build tinygo.wasm
// +build tinygo.wasm

package sys

// / Verifies that a signature is valid for an address and plaintext.
// /
// / Returns 0 on success, or -1 if the signature fails to validate.
// /
// / # Errors
// /
// / | Error             | Reason                                               |
// / |-------------------|------------------------------------------------------|
// / | `IllegalArgument` | signature, address, or plaintext buffers are invalid |
//
//go:wasm-module crypto
//export verify_signature
func cryptoVerifySignature(ret uintptr, sigOff uintptr, sigLen uint32, addrOff uintptr, addrLen uint32, plainTextOff uintptr, plainTextLen uint32) uint32

// / Hashes input data using blake2b with 256 bit output.
// /
// / The output buffer must be sized to a minimum of 32 bytes.
// /
// / # Errors
// /
// / | Error             | Reason                                          |
// / |-------------------|-------------------------------------------------|
// / | `IllegalArgument` | the input buffer does not point to valid memory |
//
//go:wasm-module crypto
//export hash_blake2b
func cryptoHashBlake2b(ret uintptr, dataOff uintptr, dataLen uint32) uint32

// / Computes an unsealed sector CID (CommD) from its constituent piece CIDs
// / (CommPs) and sizes.
// /
// / Writes the CID in the provided output buffer, and returns the length of
// / the written CID.
// /
// / # Errors
// /
// / | Error             | Reason                   |
// / |-------------------|--------------------------|
// / | `IllegalArgument` | an argument is malformed |
//
//go:wasm-module crypto
//export compute_unsealed_sector_cid
func cryptoComputeUnsealedSectorCid(ret uintptr, proofType int64, piecesOff uintptr, pieceLen uint32, cidPtr uintptr, cidLen uint32) uint32

// / Verifies a sector seal proof.
// /
// / Returns 0 to indicate that the proof was valid, -1 otherwise.
// /
// / # Errors
// /
// / | Error             | Reason                   |
// / |-------------------|--------------------------|
// / | `IllegalArgument` | an argument is malformed |
//
//go:wasm-module crypto
//export verify_seal
func cryptoVerifySeal(ret uintptr, infoOff uintptr, infoLen uint32) uint32

// / Verifies a window proof of spacetime.
// /
// / Returns 0 to indicate that the proof was valid, -1 otherwise.
// /
// / # Errors
// /
// / | Error             | Reason                   |
// / |-------------------|--------------------------|
// / | `IllegalArgument` | an argument is malformed |
//
//go:wasm-module crypto
//export verify_post
func cryptoVerifyPost(ret uintptr, infoOff uintptr, infoLen uint32) uint32

// / Verifies that two block headers provide proof of a consensus fault.
// /
// / Returns a 0 status if a consensus fault was recognized, along with the
// / BlockId containing the fault details. Otherwise, a -1 status is returned,
// / and the second result parameter must be ignored.
// /
// / # Errors
// /
// / | Error             | Reason                                |
// / |-------------------|---------------------------------------|
// / | `LimitExceeded`   | exceeded lookback limit finding block |
// / | `IllegalArgument` | an argument is malformed              |
//
//go:wasm-module crypto
//export verify_consensus_fault
func cryptoVerifyConsensusFault(ret uintptr, h1Off uintptr, h1Len uint32, h2Off uintptr, h2Len uint32, extraOff uintptr, extraLen uint32) uint32

// / Verifies an aggregated batch of sector seal proofs.
// /
// / Returns 0 to indicate that the proof was valid, -1 otherwise.
// /
// / # Errors
// /
// / | Error             | Reason                         |
// / |-------------------|--------------------------------|
// / | `LimitExceeded`   | exceeds seal aggregation limit |
// / | `IllegalArgument` | an argument is malformed       |
//
//go:wasm-module crypto
//export verify_aggregate_seals
func cryptoVerifyAggregateSeals(ret uintptr, aggOff uintptr, aggLen uint32) uint32

// / Verifies a replica update proof.
// /
// / Returns 0 to indicate that the proof was valid, -1 otherwise.
// /
// / # Errors
// /
// / | Error             | Reason                        |
// / |-------------------|-------------------------------|
// / | `LimitExceeded`   | exceeds replica update limit  |
// / | `IllegalArgument` | an argument is malformed      |
//
//go:wasm-module crypto
//export verify_replica_update
func cryptoVerifyReplicaUpdate(ret uintptr, repOff uintptr, repLen uint32) uint32

// / Verifies an aggregated batch of sector seal proofs.
// /
// / # Errors
// /
// / | Error             | Reason                   |
// / |-------------------|--------------------------|
// / | `IllegalArgument` | an argument is malformed |
//
//go:wasm-module crypto
//export batch_verify_seals
func cryptoBatchVerifySeals(batchOff uintptr, batLen uint32, resultOff uintptr) uint32

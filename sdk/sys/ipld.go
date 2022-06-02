package sys

import (
	"fmt"
	"unsafe"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

/// Opens a block from the "reachable" set, returning an ID for the block, its codec, and its
/// size in bytes.
///
/// - The reachable set is initialized to the root.
/// - The reachable set is extended to include the direct children of loaded blocks until the
///   end of the invocation.
///
/// # Arguments
///
/// - `cid` the location of the input CID (in wasm memory).
///
/// # Errors
///
/// | Error               | Reason                                      |
/// |---------------------|---------------------------------------------|
/// | [`NotFound`]        | the target block isn't in the reachable set |
/// | [`IllegalArgument`] | there's something wrong with the CID        |
//go:wasm-module ipld
//export block_open
func ipldOpen(ret uintptr, cid uintptr) uint32

/// Creates a new block, returning the block's ID. The block's children must be in the reachable
/// set. The new block isn't added to the reachable set until the CID is computed.
///
/// # Arguments
///
/// - `codec` is the codec of the block.
/// - `data` and `len` specify the location and length of the block data.
///
/// # Errors
///
/// | Error               | Reason                                                  |
/// |---------------------|---------------------------------------------------------|
/// | [`LimitExceeded`]   | the block is too big                                    |
/// | [`NotFound`]        | one of the blocks's children isn't in the reachable set |
/// | [`IllegalCodec`]    | the passed codec isn't supported                        |
/// | [`Serialization`]   | the passed block doesn't match the passed codec         |
/// | [`IllegalArgument`] | the block isn't in memory, etc.                         |
//go:wasm-module ipld
//export block_create
func ipldCreate(ret uintptr, codec uint64, data uintptr, len uint32) uint32

/// Reads the block identified by `id` into `obuf`, starting at `offset`, reading _at most_
/// `max_len` bytes.
///
/// Returns the difference between the length of the block and `offset + max_len`. This can be
/// used to find the end of the block relative to the buffer the block is being read into:
///
/// - A zero return value means that the block was read into the output buffer exactly.
/// - A positive return value means that that many more bytes need to be read.
/// - A negative return value means that the buffer should be truncated by the return value.
///
/// # Arguments
///
/// - `id` is ID of the block to read.
/// - `offset` is the offset in the block to start reading.
/// - `obuf` is the output buffer (in wasm memory) where the FVM will write the block data.
/// - `max_len` is the maximum amount of block data to read.
///
/// Passing a length/offset that exceed the length of the block will not result in an error, but
/// will result in no data being read and a negative return value indicating where the block
/// actually ended (relative to `offset + max_len`).
///
/// # Errors
///
/// | Error               | Reason                                            |
/// |---------------------|---------------------------------------------------|
/// | [`InvalidHandle`]   | if the handle isn't known.                        |
/// | [`IllegalArgument`] | if the passed buffer isn't valid, in memory, etc. |
//go:wasm-module ipld
//export block_read
func ipldRead(ret uintptr, id uint32, offset uint32, obuf uintptr, max_len uint32) uint32

/// Returns the codec and size of the specified block.
///
/// # Errors
///
/// | Error             | Reason                     |
/// |-------------------|----------------------------|
/// | [`InvalidHandle`] | if the handle isn't known. |
//go:wasm-module ipld
//export block_stat
func ipldStat(ret uintptr, id uint32) uint32

/// Computes the given block's CID, writing the resulting CID into `cid`.
///
/// The returned CID is added to the reachable set.
///
/// # Arguments
///
/// - `id` is ID of the block to link.
/// - `hash_fun` is the multicodec of the hash function to use.
/// - `hash_len` is the desired length of the hash digest.
/// - `cid` is the output buffer (in wasm memory) where the FVM will write the resulting cid.
/// - `cid_max_length` is the length of the output CID buffer.
///
/// # Returns
///
/// The length of the CID.
///
/// # Errors
///
/// | Error               | Reason                                            |
/// |---------------------|---------------------------------------------------|
/// | [`InvalidHandle`]   | if the handle isn't known.                        |
/// | [`IllegalCid`]      | hash code and/or hash length aren't supported.    |
/// | [`IllegalArgument`] | if the passed buffer isn't valid, in memory, etc. |
//go:wasm-module ipld
//export block_link
func ipldLink(ret uintptr, id uint32, hash_fun uint64, hash_len uint32, cid uintptr, cid_max_len uint32) uint32

func Open(id cid.Cid) (*types.IpldOpen, error) {
	cidBuf := make([]byte, types.MAX_CID_LEN)
	copy(cidBuf, id.Bytes())
	cidBufPtr, _ := GetSlicePointerAndLen(cidBuf)

	result := new(types.IpldOpen)
	code := ipldOpen(uintptr(unsafe.Pointer(result)), cidBufPtr)
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("unable to open ipld %s", id.String()))
	}
	return result, nil
}

func Create(codec uint64, data []byte) (uint32, error) {
	result := uint32(0)
	dataPtr, dataLen := GetSlicePointerAndLen(data)
	code := ipldCreate(uintptr(unsafe.Pointer(&result)), codec, dataPtr, dataLen)
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}
	return result, nil
}

func Read(id uint32, offset uint32, buf []byte) (uint32, error) {
	result := uint32(0)
	bufPtr, bufLen := GetSlicePointerAndLen(buf)
	code := ipldRead(uintptr(unsafe.Pointer(&result)), id, offset, bufPtr, bufLen)
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to read ipld ")
	}
	return result, nil
}

func Stat(id uint32) (*types.IpldStat, error) {
	result := new(types.IpldStat)
	code := ipldStat(uintptr(unsafe.Pointer(result)), id)
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to read ipld ")
	}

	return result, nil
}

func BlockLink(id uint32, hash_fun uint64, hash_len uint32, cidBuf []byte) (uint32, error) {
	result := uint32(0)
	cidBufPtr, cidBufLen := GetSlicePointerAndLen(cidBuf)
	code := ipldLink(uintptr(unsafe.Pointer(&result)), id, hash_fun, hash_len, cidBufPtr, cidBufLen)
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to read ipld ")
	} else {
		return result, nil
	}
}

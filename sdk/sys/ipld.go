package sys

import (
	"fmt"
	"unsafe"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

//go:wasm-module ipld
//export open
func ipldOpen(ret uintptr, cid uintptr) uint32

//go:wasm-module ipld
//export create
func ipldCreate(ret uintptr, codec uint64, data uintptr, len uint32) uint32

//go:wasm-module ipld
//export read
func ipldRead(ret uintptr, id uint32, offset uint32, obuf uintptr, max_len uint32) uint32

//go:wasm-module ipld
//export stat
func ipldStat(ret uintptr, id uint32) uint32

//go:wasm-module ipld
//export cid
func ipldCid(ret uintptr, id uint32, hash_fun uint64, hash_len uint32, cid uintptr, cid_max_len uint32) uint32

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

func Cid(id uint32, hash_fun uint64, hash_len uint32, cidBuf []byte) (uint32, error) {
	result := uint32(0)
	cidBufPtr, cidBufLen := GetSlicePointerAndLen(cidBuf)
	code := ipldCid(uintptr(unsafe.Pointer(&result)), id, hash_fun, hash_len, cidBufPtr, cidBufLen)
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to read ipld ")
	} else {
		return result, nil
	}
}

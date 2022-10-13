//go:build !simulated
// +build !simulated

package sys

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func Open(ctx context.Context, id cid.Cid) (*types.IpldOpen, error) {
	cidBuf := make([]byte, types.MaxCidLen)
	copy(cidBuf, id.Bytes())
	cidBufPtr, _ := GetSlicePointerAndLen(cidBuf)

	result := new(types.IpldOpen)
	code := ipldOpen(uintptr(unsafe.Pointer(result)), cidBufPtr)
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("unable to open ipld %s", id.String()))
	}
	return result, nil
}

func Create(ctx context.Context, codec uint64, data []byte) (uint32, error) {
	result := uint32(0)
	dataPtr, dataLen := GetSlicePointerAndLen(data)
	code := ipldCreate(uintptr(unsafe.Pointer(&result)), codec, dataPtr, dataLen)
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}
	return result, nil
}

func Read(ctx context.Context, id uint32, offset, size uint32) ([]byte, uint32, error) {
	result := uint32(0)
	buf := make([]byte, size)
	bufPtr, bufLen := GetSlicePointerAndLen(buf)
	code := ipldRead(uintptr(unsafe.Pointer(&result)), id, offset, bufPtr, bufLen)
	if code != 0 {
		return nil, 0, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to read ipld ")
	}
	return buf, result, nil
}

func Stat(ctx context.Context, id uint32) (*types.IpldStat, error) {
	result := new(types.IpldStat)
	code := ipldStat(uintptr(unsafe.Pointer(result)), id)
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to read ipld ")
	}
	return result, nil
}

func BlockLink(ctx context.Context, id uint32, hashFun uint64, hashLen uint32, cidBuf []byte) (cid.Cid, error) {
	result := uint32(0)
	cidBufPtr, cidBufLen := GetSlicePointerAndLen(cidBuf)
	code := ipldLink(uintptr(unsafe.Pointer(&result)), id, hashFun, hashLen, cidBufPtr, cidBufLen)
	if code != 0 {
		return cid.Undef, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to read ipld ")
	}
	if int(cidBufLen) > len(cidBuf) {
		panic(fmt.Sprintf("CID too big: %d > %d", cidBufLen, len(cidBuf)))
	}
	_, cid, err := cid.CidFromBytes(cidBuf)
	return cid, err
}

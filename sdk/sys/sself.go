//go:build !simulate
// +build !simulate

package sys

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/filecoin-project/go-state-types/abi"

	addr "github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func SelfRoot(_ context.Context) (cid.Cid, error) {
	// I really hate this CID interface. Why can't I just have bytes?
	result := uint32(0)
	cidBuf := make([]byte, types.MaxCidLen)
	cidBufPtr, cidBufLen := GetSlicePointerAndLen(cidBuf)
	code := sselfRoot(uintptr(unsafe.Pointer(&result)), cidBufPtr, cidBufLen)
	if code != 0 {
		return cid.Undef, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unexpected error from `self::root` syscall:")
	}
	if int(cidBufLen) > len(cidBuf) {
		// TODO: re-try with a larger buffer?
		panic(fmt.Sprintf("CID too big: %d > %d", cidBufLen, len(cidBuf)))
	}
	_, cid, err := cid.CidFromBytes(cidBuf)
	return cid, err
}

func SelfSetRoot(_ context.Context, id cid.Cid) error {
	buf := make([]byte, types.MaxCidLen)
	copy(buf, id.Bytes())
	cidBufPtr, _ := GetSlicePointerAndLen(buf)
	code := sselfSetRoot(cidBufPtr)
	if code != 0 {
		return ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unexpected error from `self::set_root` syscall:")
	}
	return nil

}

func SelfCurrentBalance(_ context.Context) (*abi.TokenAmount, error) {
	result := new(fvmTokenAmount)
	code := selfCurrentBalance(uintptr(unsafe.Pointer(result)))
	if code != 0 {
		return nil, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unable to create ipld")
	}
	return result.TokenAmount(), nil
}

func SelfDestruct(_ context.Context, addr addr.Address) error {
	addrPtr, addrLen := GetSlicePointerAndLen(addr.Bytes())
	code := selfDestruct(addrPtr, addrLen)
	if code != 0 {
		return ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unexpected error from `self::self_destruct` syscall:")
	}

	return nil
}

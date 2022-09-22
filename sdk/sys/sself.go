//go:build !simulated
// +build !simulated

package sys

import (
	"fmt"
	"unsafe"

	addr "github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func SelfRoot() (cid.Cid, error) {
	// I really hate this CID interface. Why can't I just have bytes?
	result := uint32(0)
	cidBuf := make([]byte, types.MaxCidLen)
	cidBufPtr, cidBufLen := GetSlicePointerAndLen(cidBuf)
	code := sselfRoot(uintptr(unsafe.Pointer(&result)), cidBufPtr, cidBufLen)
	if code != 0 {
		return cid.Undef, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}
	if int(cidBufLen) > len(cidBuf) {
		// TODO: re-try with a larger buffer?
		panic(fmt.Sprintf("CID too big: %d > %d", cidBufLen, len(cidBuf)))
	}
	_, cid, err := cid.CidFromBytes(cidBuf)
	return cid, err
}

func SelfSetRoot(id cid.Cid) error {
	buf := make([]byte, types.MaxCidLen)
	copy(buf, id.Bytes())
	cidBufPtr, _ := GetSlicePointerAndLen(buf)
	code := sselfSetRoot(cidBufPtr)
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}
	return nil

}

func SelfCurrentBalance() (*types.TokenAmount, error) {
	result := new(types.TokenAmount)
	code := selfCurrentBalance(uintptr(unsafe.Pointer(result)))
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}
	return result, nil
}

func SelfDestruct(addr addr.Address) error {
	addrPtr, addrLen := GetSlicePointerAndLen(addr.Bytes())
	code := selfDestruct(addrPtr, addrLen)
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}

	return nil
}

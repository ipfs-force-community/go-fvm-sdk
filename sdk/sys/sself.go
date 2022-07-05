package sys

import (
	"unsafe"

	addr "github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func SelfRoot(cidBuf []byte) (uint32, error) {
	result := uint32(0)
	cidBufPtr, cidBufLen := GetSlicePointerAndLen(cidBuf)
	code := sselfRoot(uintptr(unsafe.Pointer(&result)), cidBufPtr, cidBufLen)
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}
	return result, nil
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

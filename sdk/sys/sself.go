//go:build !simulated
// +build !simulated

package sys

import (
	"context"
	"fmt"
	"unsafe"

	addr "github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func SelfRoot(ctx context.Context) (cid.Cid, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.SelfRoot()
	}

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

func SelfSetRoot(ctx context.Context, id cid.Cid) error {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.SelfSetRoot(id)
	}

	buf := make([]byte, types.MaxCidLen)
	copy(buf, id.Bytes())
	cidBufPtr, _ := GetSlicePointerAndLen(buf)
	code := sselfSetRoot(cidBufPtr)
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}
	return nil

}

func SelfCurrentBalance(ctx context.Context) (*types.TokenAmount, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.SelfCurrentBalance()
	}

	result := new(types.TokenAmount)
	code := selfCurrentBalance(uintptr(unsafe.Pointer(result)))
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}
	return result, nil
}

func SelfDestruct(ctx context.Context, addr addr.Address) error {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.SelfDestruct(addr)
	}

	addrPtr, addrLen := GetSlicePointerAndLen(addr.Bytes())
	code := selfDestruct(addrPtr, addrLen)
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}

	return nil
}

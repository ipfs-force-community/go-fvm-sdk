//go:build !simulated
// +build !simulated

// Package sys : a go-fvm-sdk system calls
package sys

import (
	"context"
	"fmt"
	"unsafe"

	address "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func ResolveAddress(ctx context.Context, addr address.Address) (abi.ActorID, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.ResolveAddress(addr)
	}

	if addr.Protocol() == address.ID {
		actid, err := address.IDFromAddress(addr)
		return abi.ActorID(actid), err
	}
	addrBufPtr, addrBufLen := GetSlicePointerAndLen(addr.Bytes())
	var result abi.ActorID
	code := actorResolveAddress(uintptr(unsafe.Pointer(&result)), addrBufPtr, addrBufLen)
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to resolve address")
	}
	return result, nil
}

func GetActorCodeCid(ctx context.Context, addr address.Address) (*cid.Cid, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.GetActorCodeCid(addr)
	}

	addrBufPtr, addrBufLen := GetSlicePointerAndLen(addr.Bytes())
	buf := make([]byte, types.MaxCidLen)
	bufPtr, bufLen := GetSlicePointerAndLen(buf)
	var result int32
	code := actorGetActorCodeCid(uintptr(unsafe.Pointer(&result)), addrBufPtr, addrBufLen, bufPtr, bufLen)
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("unable to get actor code id from address %s", addr))
	}

	if result == 0 {
		_, result, err := cid.CidFromBytes(buf)
		if err != nil {
			return nil, err
		}
		return &result, nil
	} else {
		return nil, nil
	}
}

func ResolveBuiltinActorType(ctx context.Context, codeCid cid.Cid) (types.ActorType, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.ResolveBuiltinActorType(codeCid)
	}

	addrBufPtr, _ := GetSlicePointerAndLen(codeCid.Bytes())
	var result types.ActorType
	code := actorResolveBuiltinActorType(uintptr(unsafe.Pointer(&result)), addrBufPtr)
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("unable to resolve builtin actor type for cid %s", codeCid))
	}
	return result, nil
}

func GetCodeCidForType(ctx context.Context, actorT types.ActorType) (cid.Cid, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.GetCodeCidForType(actorT)
	}

	buf := make([]byte, types.MaxCidLen)
	bufPtr, bufLen := GetSlicePointerAndLen(buf)

	var cidLen int32
	code := actorGetCodeCidForType(uintptr(unsafe.Pointer(&cidLen)), int32(actorT), bufPtr, bufLen)
	if code != 0 {
		return cid.Undef, ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("unable to get code cid for type %d", actorT))
	}
	_, result, err := cid.CidFromBytes(buf[:cidLen])
	if err != nil {
		return cid.Undef, err
	}
	return result, nil
}

func NewActorAddress(ctx context.Context) (address.Address, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.NewActorAddress()
	}

	buf := make([]byte, types.MaxActorAddrLen)
	bufPtr, bufLen := GetSlicePointerAndLen(buf)

	var addrLen uint32
	code := actorNewActorAddress(uintptr(unsafe.Pointer(&addrLen)), bufPtr, bufLen)
	if code != 0 {
		return address.Undef, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create actor address")
	}
	return address.NewFromBytes(buf[:addrLen])
}

func CreateActor(ctx context.Context, actorID abi.ActorID, codeCid cid.Cid) error {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.CreateActor(actorID, codeCid)
	}

	addrBufPtr, _ := GetSlicePointerAndLen(codeCid.Bytes())
	code := actorCreateActor(uint64(actorID), addrBufPtr)
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("unable to create actor type %d code cid %s", actorID, codeCid))
	}
	return nil
}

//go:build !simulate
// +build !simulate

package sys

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func ResolveAddress(_ context.Context, addr address.Address) (abi.ActorID, error) {
	if addr.Protocol() == address.ID {
		actid, err := address.IDFromAddress(addr)
		return abi.ActorID(actid), err
	}
	addrBufPtr, addrBufLen := GetSlicePointerAndLen(addr.Bytes())
	var result abi.ActorID
	code := actorResolveAddress(uintptr(unsafe.Pointer(&result)), addrBufPtr, addrBufLen)
	if code != 0 {
		return 0, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unexpected address resolution failure: "+ferrors.EnToString(code))
	}
	return result, nil
}

func LookupAddress(_ context.Context, actorID abi.ActorID) (address.Address, error) {

	buf := make([]byte, types.MaxCidLen)
	bufPtr, bufLen := GetSlicePointerAndLen(buf)
	var addrLen uint32
	code := actorLookupAddress(uintptr(unsafe.Pointer(&addrLen)), uint64(actorID), bufPtr, bufLen)
	if code != 0 {
		return address.Undef, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unexpected address resolution failure: "+ferrors.EnToString(code))
	}
	//
	addr, err := address.NewFromBytes(buf[:addrLen])
	if err != nil {
		Abort(context.Background(), uint32(ferrors.NotFound), fmt.Sprintf("%v", buf[:addrLen]))
	}
	return addr, err
}

func GetActorCodeCid(ctx context.Context, addr address.Address) (*cid.Cid, error) {
	actorID, err := ResolveAddress(ctx, addr)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, types.MaxActorAddrLen)
	bufPtr, bufLen := GetSlicePointerAndLen(buf)
	var cidLen int32

	code := actorGetActorCodeCid(uintptr(unsafe.Pointer(&cidLen)), uint64(actorID), bufPtr, bufLen)
	if code != 0 {
		return nil, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unexpected code cid resolution failure: "+ferrors.EnToString(code))
	}

	_, result, err := cid.CidFromBytes(buf)
	if err != nil {
		return nil, err
	}
	return &result, nil

}

func ResolveBuiltinActorType(_ context.Context, codeCid cid.Cid) (types.ActorType, error) {
	addrBufPtr, _ := GetSlicePointerAndLen(codeCid.Bytes())
	var result types.ActorType
	code := actorResolveBuiltinActorType(uintptr(unsafe.Pointer(&result)), addrBufPtr)
	if code != 0 {
		return 0, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "failed to determine if CID belongs to builtin actor:"+ferrors.EnToString(code))
	}
	return result, nil
}

func GetCodeCidForType(_ context.Context, actorT types.ActorType) (cid.Cid, error) {
	buf := make([]byte, types.MaxCidLen)
	bufPtr, bufLen := GetSlicePointerAndLen(buf)

	var cidLen int32
	code := actorGetCodeCidForType(uintptr(unsafe.Pointer(&cidLen)), int32(actorT), bufPtr, bufLen)
	if code != 0 {
		return cid.Undef, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "failed to get CodeCID for type: "+ferrors.EnToString(code))
	}
	_, result, err := cid.CidFromBytes(buf[:cidLen])
	if err != nil {
		return cid.Undef, err
	}
	return result, nil
}

func NewActorAddress(_ context.Context) (address.Address, error) {
	buf := make([]byte, types.MaxActorAddrLen)
	bufPtr, bufLen := GetSlicePointerAndLen(buf)

	var addrLen uint32
	code := actorNewActorAddress(uintptr(unsafe.Pointer(&addrLen)), bufPtr, bufLen)
	if code != 0 {
		return address.Undef, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "failed to create a new actor address "+ferrors.EnToString(code))
	}
	return address.NewFromBytes(buf[:addrLen])
}

func CreateActor(_ context.Context, actorID abi.ActorID, codeCid cid.Cid, address address.Address) error {
	cidBufPtr, _ := GetSlicePointerAndLen(codeCid.Bytes())
	addrBufPtr, addrBufLen := GetSlicePointerAndLen(address.Bytes())
	code := actorCreateActor(uint64(actorID), cidBufPtr, addrBufPtr, addrBufLen)
	if code != 0 {
		return ferrors.NewSysCallError(ferrors.ErrorNumber(code), fmt.Sprintf("unable to create actor type %d code cid %s address %s", actorID, codeCid, address))
	}
	return nil
}

func BalanceOf(_ context.Context, actorID abi.ActorID) (*abi.TokenAmount, error) {
	tokenAmount := new(fvmTokenAmount)
	code := actorBalanceOf(uintptr(unsafe.Pointer(tokenAmount)), uint64(actorID))
	if code != 0 {
		return nil, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unexpected error: "+ferrors.EnToString(code))
	}
	return tokenAmount.TokenAmount(), nil
}

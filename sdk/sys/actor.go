package sys

import (
	"fmt"
	"unsafe"

	address "github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

/// Resolves the ID address of an actor.
///
/// # Errors
///
/// | Error             | Reason                                                    |
/// |-------------------|-----------------------------------------------------------|
/// | `NotFound`        | target actor doesn't exist                                |
/// | `IllegalArgument` | if the passed address buffer isn't valid, in memory, etc. |
//go:wasm-module actor
//export resolve_address
func actorResolveAddress(ret uintptr, addr_off uintptr, addr_len uint32) uint32

/// Gets the CodeCID of an actor by address.
///
/// Returns the
///
/// # Errors
///
/// | Error             | Reason                                                    |
/// |-------------------|-----------------------------------------------------------|
/// | `NotFound`        | target actor doesn't exist                                |
/// | `IllegalArgument` | if the passed address buffer isn't valid, in memory, etc. |
//go:wasm-module actor
//export get_actor_code_cid
func actorGetActorCodeCid(ret uintptr, addr_off uintptr, addr_len uint32, obuf_off uintptr, obuf_len uint32) uint32

/// Determines whether the specified CodeCID belongs to that of a builtin
/// actor and which. Returns 0 if unrecognized. Can only fail due to
/// internal errors.
//go:wasm-module actor
//export resolve_builtin_actor_type
func actorResolveBuiltinActorType(ret uintptr, cid_off uintptr) uint32

/// Returns the CodeCID for the given built-in actor type. Aborts with exit
/// code IllegalArgument if the supplied type is invalid. Returns the
/// length of the written CID written to the output buffer. Can only
/// return a failure due to internal errors.
//go:wasm-module actor
//export get_code_cid_for_type
func actorGetCodeCidForType(ret uintptr, typ int32, obuf_off uintptr, obuf_len uint32) uint32

/// Generates a new actor address for an actor deployed
/// by the calling actor.
///
/// **Privledged:** May only be called by the init actor.
//go:wasm-module actor
//export new_actor_address
func actorNewActorAddress(ret uintptr, obuf_off uintptr, obuf_len uint32) uint32

/// Creates a new actor of the specified type in the state tree, under
/// the provided address.
///
/// **Privledged:** May only be called by the init actor.
//go:wasm-module actor
//export create_actor
func actorCreateActor(actor_id uint64, typ_off uintptr) uint32

func ResolveAddress(addr address.Address) (types.ActorId, error) {
	if addr.Protocol() == address.ID {
		return address.IDFromAddress(addr)
	}
	addrBufPtr, addrBufLen := GetSlicePointerAndLen(addr.Bytes())
	var result types.ActorId
	code := actorResolveAddress(uintptr(unsafe.Pointer(&result)), addrBufPtr, addrBufLen)
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("unable to resolve address %s"))
	}
	return result, nil
}

func GetActorCodeCid(addr address.Address) (*cid.Cid, error) {
	addrBufPtr, addrBufLen := GetSlicePointerAndLen(addr.Bytes())
	buf := make([]byte, types.MAX_CID_LEN)
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

func ResolveBuiltinActorType(codeCid cid.Cid) (types.ActorType, error) {
	addrBufPtr, _ := GetSlicePointerAndLen(codeCid.Bytes())
	var result types.ActorType
	code := actorResolveBuiltinActorType(uintptr(unsafe.Pointer(&result)), addrBufPtr)
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("unable to resolve builtin actor type for cid %s", codeCid))
	}
	return result, nil
}

func GetCodeCidForType(actorT types.ActorType) (cid.Cid, error) {
	buf := make([]byte, types.MAX_CID_LEN)
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

func NewActorAddress() (address.Address, error) {
	buf := make([]byte, types.MAX_ACTOR_ADDR_LEN)
	bufPtr, bufLen := GetSlicePointerAndLen(buf)

	var addrLen uint32
	code := actorNewActorAddress(uintptr(unsafe.Pointer(&addrLen)), bufPtr, bufLen)
	if code != 0 {
		return address.Undef, ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("unable to create actor address"))
	}
	return address.NewFromBytes(buf[:addrLen])
}

func CreateActor(actorId types.ActorId, codeCid cid.Cid) error {
	addrBufPtr, _ := GetSlicePointerAndLen(codeCid.Bytes())
	code := actorCreateActor(actorId, addrBufPtr)
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("unable to create actor type %d code cid %s", actorId, codeCid))
	}
	return nil
}

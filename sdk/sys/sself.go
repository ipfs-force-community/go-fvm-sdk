package sys

import (
	"unsafe"

	addr "github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

/// Gets the current root for the calling actor.
///
/// If the CID doesn't fit in the specified maximum length (and/or the length is 0), this
/// function returns the required size and does not update the cid buffer.
///
/// # Errors
///
/// | Error              | Reason                                             |
/// |--------------------|----------------------------------------------------|
/// | `IllegalOperation` | actor hasn't set the root yet, or has been deleted |
/// | `IllegalArgument`  | if the passed buffer isn't valid, in memory, etc.  |
//go:wasm-module self
//export root
func sselfRoot(ret uintptr, cid uintptr, cidMaxLen uint32) uint32

/// Sets the root CID for the calling actor. The new root must be in the reachable set.
///
/// # Errors
///
/// | Error              | Reason                                         |
/// |--------------------|------------------------------------------------|
/// | `IllegalOperation` | actor has been deleted                         |
/// | `NotFound`         | specified root CID is not in the reachable set |
//go:wasm-module self
//export set_root
func sselfSetRoot(cid uintptr) uint32

/// Gets the current balance for the calling actor.
///
/// # Errors
///
/// None.
//go:wasm-module self
//export current_balance
func selfCurrentBalance(ret uintptr) uint32

/// Destroys the calling actor, sending its current balance
/// to the supplied address, which cannot be itself.
///
/// # Errors
///
/// | Error             | Reason                                                         |
/// |-------------------|----------------------------------------------------------------|
/// | `NotFound`        | beneficiary isn't found                                        |
/// | `Forbidden`       | beneficiary is not allowed (usually means beneficiary is self) |
/// | `IllegalArgument` | if the passed address buffer isn't valid, in memory, etc.      |
//go:wasm-module self
//export self_destruct
func selfDestruct(addrOff uintptr, addrLen uint32) uint32

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
	buf := make([]byte, types.MAX_CID_LEN)
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
	code := selfDestruct(addrPtr, uint32(addrLen))
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}

	return nil
}

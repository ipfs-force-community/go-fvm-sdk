package sys

import (
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

// TODO: name for debugging & tracing?
// We could also _not_ feed that through to the outside?

/// Charge gas.
///
/// # Arguments
///
/// - `name_off` and `name_len` specify the location and length of the "name" of the gas charge,
///   for debugging.
/// - `amount` is the amount of gas to charge.
///
/// # Errors
///
/// | Error               | Reason               |
/// |---------------------|----------------------|
/// | [`IllegalArgument`] | invalid name buffer. |
//go:wasm-module gas
//export charge
func gasCharge(name_off uintptr, name_len uint32, amount uint64) uint32

/// Returns the amount of gas remaining.
/// TODO not implemented.
///go:wasm-module gas
///export remaining
///func gasRemaining(ret uintptr) uint32

/// Charge gas for the operation identified by name.
func Charge(name string, compute uint64) error {
	nameBufPtr, nameBufLen := GetStringPointerAndLen(name)
	code := gasCharge(nameBufPtr, nameBufLen, compute)
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("charge gas to %s", name))
	}
	return nil
}

//go:build tinygo.wasm
// +build tinygo.wasm

package sys

// TODO: name for debugging & tracing?
// We could also _not_ feed that through to the outside?

// Charge gas.
// /
// # Arguments
// /
//   - `name_off` and `name_len` specify the location and length of the "name" of the gas charge,
//     for debugging.
//   - `amount` is the amount of gas to charge.
//
// /
// # Errors
// /
// | Error               | Reason               |
// |---------------------|----------------------|
// | [`IllegalArgument`] | invalid name buffer. |
//
//go:wasm-module gas
//export charge
func gasCharge(name_off uintptr, name_len uint32, amount uint64) uint32

// / Returns the amount of gas remaining.
//
//go:wasm-module gas
//export available
func gasAvailable(ret uintptr) uint32

//go:build tinygo.wasm
// +build tinygo.wasm

package sys

//! Syscalls related to eventing.
/// Emits an actor event to be recorded in the receipt.
///
/// Expects a DAG-CBOR representation of the ActorEvent struct.
///
/// # Errors
///
/// | Error               | Reason                                                              |
/// |---------------------|---------------------------------------------------------------------|
/// | [`IllegalArgument`] | entries failed to validate due to improper encoding or invalid data |

//go:wasm-module event
//export emit_event
func emitEvent(evt_ptr uintptr, evt_len uint32) uint32

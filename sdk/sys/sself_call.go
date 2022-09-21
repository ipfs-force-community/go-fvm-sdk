//go:build tinygo.wasm
// +build tinygo.wasm

package sys

// / Gets the current root for the calling actor.
// /
// / If the CID doesn't fit in the specified maximum length (and/or the length is 0), this
// / function returns the required size and does not update the cid buffer.
// /
// / # Errors
// /
// / | Error              | Reason                                             |
// / |--------------------|----------------------------------------------------|
// / | `IllegalOperation` | actor hasn't set the root yet, or has been deleted |
// / | `IllegalArgument`  | if the passed buffer isn't valid, in memory, etc.  |
//
//go:wasm-module self
//export root
func sselfRoot(ret uintptr, cid uintptr, cidMaxLen uint32) uint32

// / Sets the root CID for the calling actor. The new root must be in the reachable set.
// /
// / # Errors
// /
// / | Error              | Reason                                         |
// / |--------------------|------------------------------------------------|
// / | `IllegalOperation` | actor has been deleted                         |
// / | `NotFound`         | specified root CID is not in the reachable set |
//
//go:wasm-module self
//export set_root
func sselfSetRoot(cid uintptr) uint32

// / Gets the current balance for the calling actor.
// /
// / # Errors
// /
// / None.
//
//go:wasm-module self
//export current_balance
func selfCurrentBalance(ret uintptr) uint32

// / Destroys the calling actor, sending its current balance
// / to the supplied address, which cannot be itself.
// /
// / # Errors
// /
// / | Error             | Reason                                                         |
// / |-------------------|----------------------------------------------------------------|
// / | `NotFound`        | beneficiary isn't found                                        |
// / | `Forbidden`       | beneficiary is not allowed (usually means beneficiary is self) |
// / | `IllegalArgument` | if the passed address buffer isn't valid, in memory, etc.      |
//
//go:wasm-module self
//export self_destruct
func selfDestruct(addrOff uintptr, addrLen uint32) uint32

//go:build tinygo.wasm
// +build tinygo.wasm

package sys

// Resolves the ID address of an actor.
// /
// # Errors
// /
// | Error             | Reason                                                    |
// |-------------------|-----------------------------------------------------------|
// | `NotFound`        | target actor doesn't exist                                |
// | `IllegalArgument` | if the passed address buffer isn't valid, in memory, etc. |
//
//go:wasm-module actor
//export resolve_address
func actorResolveAddress(ret uintptr, addr_off uintptr, addr_len uint32) uint32

// / Looks up the "delegated" (f4) address of the target actor (if any).
// /
// / # Arguments
// /
// / `addr_buf_off` and `addr_buf_len` specify the location and length of the output buffer in
// / which to store the address.
// /
// / # Returns
// /
// / The length of the address written to the output buffer, or 0 if the target actor has no
// / delegated (f4) address.
// /
// / # Errors
// /
// / | Error               | Reason                                                           |
// / |---------------------|------------------------------------------------------------------|
// / | [`NotFound`]        | if the target actor does not exist                               |
// / | [`BufferTooSmall`]  | if the output buffer isn't large enough to fit the address       |
// / | [`IllegalArgument`] | if the output buffer isn't valid, in memory, etc.                |
//
//go:wasm-module actor
//export lookup_delegated_address
func actorLookupDelegatedAddress(ret uintptr, actor_id uint64, addr_buf_off uintptr, addr_buf_len uint32) uint32

// Gets the CodeCID of an actor by address.
// /
// Returns the
// /
// # Errors
// /
// | Error             | Reason                                                    |
// |-------------------|-----------------------------------------------------------|
// | `NotFound`        | target actor doesn't exist                                |
// | `IllegalArgument` | if the passed address buffer isn't valid, in memory, etc. |
//
//go:wasm-module actor
//export get_actor_code_cid
func actorGetActorCodeCid(ret uintptr, actor_id uint64, obuf_off uintptr, obuf_len uint32) uint32

// Determines whether the specified CodeCID belongs to that of a builtin
// actor and which. Returns 0 if unrecognized. Can only fail due to
// internal errors.
//
//go:wasm-module actor
//export get_builtin_actor_type
func actorGetBuiltinActorType(ret uintptr, cid_off uintptr) uint32

// Returns the CodeCID for the given built-in actor type. Aborts with exit
// code IllegalArgument if the supplied type is invalid. Returns the
// length of the written CID written to the output buffer. Can only
// return a failure due to internal errors.
//
//go:wasm-module actor
//export get_code_cid_for_type
func actorGetCodeCidForType(ret uintptr, typ int32, obuf_off uintptr, obuf_len uint32) uint32

// Generates a new actor address for an actor deployed
// by the calling actor.
// /
// **Privledged:** May only be called by the init actor.
//
//go:wasm-module actor
//export next_actor_address
func actorNextActorAddress(ret uintptr, obuf_off uintptr, obuf_len uint32) uint32

// Creates a new actor of the specified type in the state tree, under
// the provided address.
// /
// **Privledged:** May only be called by the init actor.
//
//go:wasm-module actor
//export create_actor
func actorCreateActor(actor_id uint64, typ_off uintptr, predictable_addr_off uintptr, predictable_addr_len uint32) uint32

//go:wasm-module actor
//export balance_of
func actorBalanceOf(ret uintptr, actor_id uint64) uint32

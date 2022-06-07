//go:build tinygo.wasm
// +build tinygo.wasm

package sys

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

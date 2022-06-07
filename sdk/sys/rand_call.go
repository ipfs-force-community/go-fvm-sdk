//go:build tinygo.wasm
// +build tinygo.wasm

package sys

/// Gets 32 bytes of randomness from the ticket chain.
///
/// # Arguments
///
/// - `tag` is the "domain separation tag" for distinguishing between different categories of
///    randomness. Think of it like extra, structured entropy.
/// - `epoch` is the epoch to pull the randomness from.
/// - `entropy_off` and `entropy_len` specify the location and length of the entropy buffer that
///    will be mixed into the system randomness.
///
/// # Errors
///
/// | Error               | Reason                  |
/// |---------------------|-------------------------|
/// | [`LimitExceeded`]   | lookback exceeds limit. |
/// | [`IllegalArgument`] | invalid buffer, etc.    |
//go:wasm-module rand
//export get_chain_randomness
func getChainRandomness(ret uintptr, tag int64, epoch int64, entropy_off uintptr, entropy_len uint32) uint32

/// Gets 32 bytes of randomness from the beacon system (currently Drand).
///
/// # Arguments
///
/// - `tag` is the "domain separation tag" for distinguishing between different categories of
///    randomness. Think of it like extra, structured entropy.
/// - `epoch` is the epoch to pull the randomness from.
/// - `entropy_off` and `entropy_len` specify the location and length of the entropy buffer that
///    will be mixed into the system randomness.
///
/// # Errors
///
/// | Error               | Reason                  |
/// |---------------------|-------------------------|
/// | [`LimitExceeded`]   | lookback exceeds limit. |
/// | [`IllegalArgument`] | invalid buffer, etc.    |
//go:wasm-module rand
//export get_beacon_randomness
func getBeaconRandomness(ret uintptr, tag int64, epoch int64, entropy_off uintptr, entropy_len uint32) uint32

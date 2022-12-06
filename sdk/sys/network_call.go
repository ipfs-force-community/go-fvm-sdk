//go:build tinygo.wasm
// +build tinygo.wasm

package sys

// Gets the circulating supply.
// /
// # Errors
// /
// None
//
//go:wasm-module network
//export total_fil_circ_supply
func networkTotalFilCircSupply(ret uintptr) uint32

// / Retrieves a tipset's CID within the last finality, if available
// /
// / # Arguments
// /
// / - `epoch` the epoch being queried.
// / - `ret_off` and `ret_len` specify the location and length of the buffer into which the
// /   tipset CID will be written.
// /
// / # Returns
// /
// / Returns the length of the CID written to the output buffer.
// /
// / # Errors
// /
// / | Error               | Reason                                       |
// / |---------------------|----------------------------------------------|
// / | [`IllegalArgument`] | specified epoch is negative or in the future |
// / | [`LimitExceeded`]   | specified epoch exceeds finality
//
//go:wasm-module network
//export tipset_cid
func networkTipsetCid(ret uintptr, epoch int64, ret_off uintptr, ret_len uint32) uint32

/// Returns the details about the network.
///
/// # Errors
///
/// None

//go:wasm-module network
//export context
func networkContext(ret uintptr) uint32

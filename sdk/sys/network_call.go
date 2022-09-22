//go:build tinygo.wasm
// +build tinygo.wasm

package sys

// Gets the base fee for the current epoch.
// /
// # Errors
// /
// None
//
//go:wasm-module network
//export base_fee
func networkBaseFee(ret uintptr) uint32

// Gets the circulating supply.
// /
// # Errors
// /
// None
//
//go:wasm-module network
//export total_fil_circ_supply
func networkTotalFilCircSupply(ret uintptr) uint32

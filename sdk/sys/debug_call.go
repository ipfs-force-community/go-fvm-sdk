//go:build tinygo.wasm
// +build tinygo.wasm

package sys

// Returns if we're in debug mode. A zero or positive return value means
// yes, a negative return value means no.
//
//go:wasm-module debug
//export enabled
func debugEnabled(ret uintptr) uint32

// Logs a message on the node.
//
//go:wasm-module debug
//export log
func debugLog(message uintptr, message_len uint32) uint32

// store_artifact   Saves an artifact to the host env. New artifacts with the same name will overwrite old ones
//
//go:wasm-module debug
//export store_artifact
func debugStoreArtifact(name_off uintptr, name_len uint32, data_off uintptr, data_len uint32) uint32

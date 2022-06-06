package main

import (
	"unsafe"
)

/// Logs a message on the node.
//go:wasm-module debug
//export log
func debugLog(message uintptr, message_len uint32) uint32 //nolint

//go:export fd_write
func fd_write(id uint32, iovs *__wasi_iovec_t, iovs_len uint, nwritten *uint) uint { //nolint
	//only support println in fvm
	errno := debugLog(uintptr(iovs.buf), uint32(iovs.bufLen))
	return uint(errno)
}

func main() {}

type __wasi_iovec_t struct { //nolint
	buf    unsafe.Pointer
	bufLen uint
}

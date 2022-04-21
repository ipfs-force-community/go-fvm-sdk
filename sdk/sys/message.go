package sys

import (
	"unsafe"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

//go:wasm-module message
//export caller
func messageCaller(ret uintptr) uint32

//go:wasm-module message
//export receiver
func messageReceiver(ret uintptr) uint32

//go:wasm-module message
//export method_number
func messageMethodNumber(ret uintptr) uint32

//go:wasm-module message
//export value_received
func messageValueReceived(ret uintptr) uint32

func Caller() (uint64, error) {
	result := uint64(0)
	code := messageCaller(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}
	return result, nil
}

func Receiver() (uint64, error) {
	result := uint64(0)
	code := messageReceiver(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}
	return result, nil
}

func MethodNumber() (uint64, error) {
	result := uint64(0)
	code := messageMethodNumber(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}
	return result, nil
}

func ValueReceived() (*types.TokenAmount, error) {
	result := new(types.TokenAmount)
	code := messageValueReceived(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to create ipld")
	}
	return result, nil
}

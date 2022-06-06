package main

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

func main() {} //nolint

//go:export invoke
func Invoke(_ uint32) uint32 { //nolint
	method_num, err := sdk.MethodNumber()
	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to get method number")
	}
	switch method_num {
	case 1:
		sdk.Abort(ferrors.USR_ILLEGAL_ARGUMENT, "test_abort USR_ILLEGAL_ARGUMENT")
	case 2:
		sdk.Abort(ferrors.SYS_SENDER_STATE_INVALID, "test_abort SYS_SENDER_STATE_INVALID")
	}
	return 0
}

package sdk

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

func Caller() (uint64, error) {
	return sys.Caller()
}

func Receiver() (uint64, error) {
	return sys.Receiver()
}

func MethodNumber() (uint64, error) {
	return sys.MethodNumber()
}

func ValueReceived() (*types.TokenAmount, error) {
	return sys.ValueReceived()
}

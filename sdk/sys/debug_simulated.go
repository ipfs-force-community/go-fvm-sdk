//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/fvm"
)

func Enabled() (bool, error) {
	return fvm.MockFvmInstance.Enabled()
}

func Log(msg string) error {
	return fvm.MockFvmInstance.Log(msg)
}

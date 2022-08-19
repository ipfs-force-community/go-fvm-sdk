//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/fvm"
)

// Charge charge gas for the operation identified by name.
func Charge(name string, compute uint64) error {
	return fvm.MockFvmInstance.Charge(name, compute)
}

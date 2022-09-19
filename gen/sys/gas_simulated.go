//go:build simulatedd
// +build simulatedd

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
)

// Charge charge gas for the operation identified by name.
func Charge(name string, compute uint64) error {
	return simulated.DefaultFsm.Charge(name, compute)
}

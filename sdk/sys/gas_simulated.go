//go:build simulate
// +build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/internal/simulated"
)

// Charge charge gas for the operation identified by name.
func Charge(name string, compute uint64) error {
	return simulated.SimulatedInstance.Charge(name, compute)
}

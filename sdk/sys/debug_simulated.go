//go:build simulate
// +build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/internal/simulated"
)

func Enabled() (bool, error) {
	return simulated.SimulatedInstance.Enabled()
}

func Log(msg string) error {
	return simulated.SimulatedInstance.Log(msg)
}

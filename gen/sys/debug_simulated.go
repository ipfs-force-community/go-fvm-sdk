//go:build simulatedd
// +build simulatedd

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
)

func Enabled() (bool, error) {
	return simulated.DefaultFsm.Enabled()
}

func Log(msg string) error {
	return simulated.DefaultFsm.Log(msg)
}

//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/internal/simulated"
)

func Begin() {
	simulated.Begin()
}

func End() {
	simulated.End()
}

func GetSimulated() *simulated.MockSimulated {
	return simulated.SimulatedInstance
}

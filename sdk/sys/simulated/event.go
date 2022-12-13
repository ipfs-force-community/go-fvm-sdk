package simulated

import "github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

func (fvmSimulator *FvmSimulator) AppendEvent(event types.ActorEvent) {
	fvmSimulator.events = append(fvmSimulator.events, event)
}

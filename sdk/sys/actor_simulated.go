//go:build simulate
// +build simulate

package sys

import (
	address "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/internal/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func ResolveAddress(addr address.Address) (abi.ActorID, error) {
	return simulated.SimulatedInstance.ResolveAddress(addr)
}

func GetActorCodeCid(addr address.Address) (*cid.Cid, error) {
	return simulated.SimulatedInstance.GetActorCodeCid(addr)
}

func ResolveBuiltinActorType(codeCid cid.Cid) (types.ActorType, error) {
	return simulated.SimulatedInstance.ResolveBuiltinActorType(codeCid)
}

func GetCodeCidForType(actorT types.ActorType) (cid.Cid, error) {
	return simulated.SimulatedInstance.GetCodeCidForType(actorT)
}

func NewActorAddress() (address.Address, error) {
	return simulated.SimulatedInstance.NewActorAddress()
}

func CreateActor(actorID abi.ActorID, codeCid cid.Cid) error {
	return simulated.SimulatedInstance.CreateActor(actorID, codeCid)
}

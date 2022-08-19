//go:build simulate

package sys

import (
	address "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/internal/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func ResolveAddress(addr address.Address) (abi.ActorID, error) {
	return simulated.MockFvmInstance.ResolveAddress(addr)
}

func GetActorCodeCid(addr address.Address) (*cid.Cid, error) {
	return simulated.MockFvmInstance.GetActorCodeCid(addr)
}

func ResolveBuiltinActorType(codeCid cid.Cid) (types.ActorType, error) {
	return simulated.MockFvmInstance.ResolveBuiltinActorType(codeCid)
}

func GetCodeCidForType(actorT types.ActorType) (cid.Cid, error) {
	return simulated.MockFvmInstance.GetCodeCidForType(actorT)
}

func NewActorAddress() (address.Address, error) {
	return simulated.MockFvmInstance.NewActorAddress()
}

func CreateActor(actorID abi.ActorID, codeCid cid.Cid) error {
	return simulated.MockFvmInstance.CreateActor(actorID, codeCid)
}

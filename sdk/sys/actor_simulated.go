//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/fvm"
)

func ResolveAddress(addr address.Address) (abi.ActorID, error) {
	return fvm.MockFvmInstance.ResolveAddress(addr)
}

func GetActorCodeCid(addr address.Address) (*cid.Cid, error) {
	return fvm.MockFvmInstance.GetActorCodeCid(addr)
}

func ResolveBuiltinActorType(codeCid cid.Cid) (types.ActorType, error) {
	return fvm.MockFvmInstance.ResolveBuiltinActorType(codeCid)
}

func GetCodeCidForType(actorT types.ActorType) (cid.Cid, error) {
	return fvm.MockFvmInstance.GetCodeCidForType(actorT)
}

func NewActorAddress() (address.Address, error) {
	return fvm.MockFvmInstance.NewActorAddress()
}

func CreateActor(actorID abi.ActorID, codeCid cid.Cid) error {
	return fvm.MockFvmInstance.CreateActor(actorID, codeCid)
}

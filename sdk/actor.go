package sdk

import (
	address "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

/// Resolves the ID address of an actor. Returns `None` if the address cannot be resolved.
/// Successfully resolving an address doesn't necessarily mean the actor exists (e.g., if the
/// addresss was already an actor ID).
func ResolveAddress(addr address.Address) (abi.ActorID, error) {
	return sys.ResolveAddress(addr)
}

/// Look up the code ID at an actor address. Returns `None` if the actor cannot be found.
func GetActorCodeCid(addr address.Address) (*cid.Cid, error) {
	return sys.GetActorCodeCid(addr)
}

/// Generates a new actor address for an actor deployed
/// by the calling actor.
func NewActorAddress() (address.Address, error) {
	return sys.NewActorAddress()
}

/// Creates a new actor of the specified type in the state tree, under
/// the provided address.
/// TODO this syscall will change to calculate the address internally.
func CreateActor(actorId abi.ActorID, codeCid cid.Cid) error {
	return sys.CreateActor(actorId, codeCid)
}

/// Determines whether the supplied CodeCID belongs to a built-in actor type,
/// and to which.
func ResolveBuiltinActorType(codeCid cid.Cid) (types.ActorType, error) {
	return sys.ResolveBuiltinActorType(codeCid)
}

/// Returns the CodeCID for a built-in actor type. Aborts with IllegalArgument
/// if the supplied type is invalid.
func GetCodeCidForType(actorT types.ActorType) (cid.Cid, error) {
	return sys.GetCodeCidForType(actorT)
}

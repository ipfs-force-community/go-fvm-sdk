// Package sdk : a go-fvm-sdk for creation actors.
package sdk

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

// ResolveAddress resolves the ID address of an actor. Returns `None` if the address cannot be resolved.
// Successfully resolving an address doesn't necessarily mean the actor exists (e.g., if the
// addresss was already an actor ID).
func ResolveAddress(ctx context.Context, addr address.Address) (abi.ActorID, error) {
	return sys.ResolveAddress(ctx, addr)
}

// GetActorCodeCid look up the code ID at an actor address. Returns `None` if the actor cannot be found.
func GetActorCodeCid(ctx context.Context, addr address.Address) (*cid.Cid, error) {
	return sys.GetActorCodeCid(ctx, addr)
}

// NewActorAddress generates a new actor address for an actor deployed
// by the calling actor.
func NewActorAddress(ctx context.Context) (address.Address, error) {
	return sys.NewActorAddress(ctx)
}

// CreateActor Creates a new actor of the specified type in the state tree, under
// the provided address.
// TODO this syscall will change to calculate the address internally.
func CreateActor(ctx context.Context, actorID abi.ActorID, codeCid cid.Cid) error {
	return sys.CreateActor(ctx, actorID, codeCid)
}

// ResolveBuiltinActorType determines whether the supplied CodeCID belongs to a built-in actor type,
// and to which.
func ResolveBuiltinActorType(ctx context.Context, codeCid cid.Cid) (types.ActorType, error) {
	return sys.ResolveBuiltinActorType(ctx, codeCid)
}

// GetCodeCidForType Returns the CodeCID for a built-in actor type. Aborts with IllegalArgument
// if the supplied type is invalid.
func GetCodeCidForType(ctx context.Context, actorT types.ActorType) (cid.Cid, error) {
	return sys.GetCodeCidForType(ctx, actorT)
}

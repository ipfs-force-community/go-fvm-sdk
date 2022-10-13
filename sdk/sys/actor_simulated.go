//go:build simulated
// +build simulated

package sys

import (
	"context"

	address "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func ResolveAddress(ctx context.Context, addr address.Address) (abi.ActorID, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.ResolveAddress(addr)
	}
	return abi.ActorID(0), nil
}

func GetActorCodeCid(ctx context.Context, addr address.Address) (*cid.Cid, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.GetActorCodeCid(addr)
	}
	return &cid.Undef, nil
}

func ResolveBuiltinActorType(ctx context.Context, codeCid cid.Cid) (types.ActorType, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.ResolveBuiltinActorType(codeCid)
	}
	return types.ActorType(0), nil
}

func GetCodeCidForType(ctx context.Context, actorT types.ActorType) (cid.Cid, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.GetCodeCidForType(actorT)
	}

	return cid.Undef, nil
}

func NewActorAddress(ctx context.Context) (address.Address, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.NewActorAddress()
	}
	return address.Undef, nil
}

func CreateActor(ctx context.Context, actorID abi.ActorID, codeCid cid.Cid) error {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.CreateActor(actorID, codeCid)
	}
	return nil
}

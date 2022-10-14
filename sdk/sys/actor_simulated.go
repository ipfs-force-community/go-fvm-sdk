//go:build simulate
// +build simulate

package sys

import (
	"context"

	address "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func ResolveAddress(ctx context.Context, addr address.Address) (abi.ActorID, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.ResolveAddress(addr)
	}
	panic(ErrorEnvValid)
}

func GetActorCodeCid(ctx context.Context, addr address.Address) (*cid.Cid, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.GetActorCodeCid(addr)
	}
	panic(ErrorEnvValid)
}

func ResolveBuiltinActorType(ctx context.Context, codeCid cid.Cid) (types.ActorType, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.ResolveBuiltinActorType(codeCid)
	}
	panic(ErrorEnvValid)
}

func GetCodeCidForType(ctx context.Context, actorT types.ActorType) (cid.Cid, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.GetCodeCidForType(actorT)
	}
	panic(ErrorEnvValid)
}

func NewActorAddress(ctx context.Context) (address.Address, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.NewActorAddress()
	}
	panic(ErrorEnvValid)
}

func CreateActor(ctx context.Context, actorID abi.ActorID, codeCid cid.Cid) error {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.CreateActor(actorID, codeCid)
	}
	panic(ErrorEnvValid)
}
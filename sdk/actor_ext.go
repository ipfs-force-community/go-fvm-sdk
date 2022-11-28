package sdk

import (
	"context"
	"errors"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

func IsAccountAddress(ctx context.Context, addr address.Address) bool {
	codeCid, err := GetActorCodeCid(ctx, addr)
	if err != nil {
		return false
	}

	actorTytp, err := ResolveBuiltinActorType(ctx, *codeCid)
	if err != nil {
		return false
	}

	return actorTytp == types.Account
}

func ResolveOrInitAddress(ctx context.Context, addr address.Address) (abi.ActorID, error) {
	actorId, err := sys.ResolveAddress(ctx, addr)
	if err == nil {
		return actorId, nil
	}
	if errors.Is(err, ferrors.NotFound) {
		return InitializeAccount(ctx, addr)
	}
	return 0, err
}

func InitializeAccount(ctx context.Context, addr address.Address) (abi.ActorID, error) {
	_, err := Send(ctx, addr, 0, nil, big.Zero())
	if err != nil {
		return 0, err
	}
	return sys.ResolveAddress(ctx, addr)
}

func SameAddress(ctx context.Context, addrA, addrB address.Address) bool {
	protocolA := addrA.Protocol()
	protocolB := addrB.Protocol()
	if protocolA == protocolB {
		return addrA == addrB
	} else {
		// attempt to resolve both to ActorID
		idA, err := ResolveAddress(ctx, addrA)
		if err != nil {
			return false
		}

		idB, err := ResolveAddress(ctx, addrB)
		if err != nil {
			return false
		}
		return idA == idB
	}
}

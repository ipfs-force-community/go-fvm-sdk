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

// IsAccountAddress use check whether specific address is  action type
func IsAccountAddress(ctx context.Context, addr address.Address) bool {
	codeCid, err := GetActorCodeCid(ctx, addr)
	if err != nil {
		return false
	}

	actorType, err := GetBuiltinActorType(ctx, codeCid)
	if err != nil {
		return false
	}

	return actorType == types.Account
}

// ResolveOrInitAddress get actor id from address, if not found create one
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

// InitializeAccount create an account actor for address
func InitializeAccount(ctx context.Context, addr address.Address) (abi.ActorID, error) {
	_, err := Send(ctx, addr, 0, nil, big.Zero())
	if err != nil {
		return 0, err
	}
	return sys.ResolveAddress(ctx, addr)
}

// SameAddress check if two address is the same actor
func SameAddress(ctx context.Context, addrA, addrB address.Address) bool {
	protocolA := addrA.Protocol()
	protocolB := addrB.Protocol()
	if protocolA == protocolB {
		return addrA == addrB
	}

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

/*
// LookupAddress get address of the actor, including f1/f3/f4
func LookupAddress(ctx context.Context, actorID abi.ActorID) (address.Address, error) {
	actorAddr, err := address.NewIDAddress(uint64(actorID))
	if err != nil {
		return address.Undef, err
	}

	codeCid, err := GetActorCodeCid(ctx, actorAddr)
	if err != nil {
		return address.Undef, err
	}

	actorType, err := GetBuiltinActorType(ctx, codeCid)
	if err != nil {
		return address.Undef, err
	}

	if actorType == types.Account {
		receipt, err := Send(ctx, actorAddr, 2, nil, big.Zero())
		if err != nil {
			return address.Undef, err
		}

		if receipt.ExitCode.IsSuccess() {
			return address.Undef, fmt.Errorf("expect success send result but got %d", receipt.ExitCode)
		}

		addr := &address.Undef
		err = addr.UnmarshalCBOR(bytes.NewReader(receipt.ReturnData))
		if err != nil {
			return address.Undef, err
		}
		return *addr, nil
	}
	if actorType == types.PlaceHolder || actorType == types.EthAccount {
		return LookupDelegatedAddress(ctx, actorID)
	}
	return address.Undef, fmt.Errorf("unsupport account code %d", actorType)
}

*/

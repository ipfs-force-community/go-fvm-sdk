package simulated

import (
	"time"

	"github.com/filecoin-project/go-state-types/builtin"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs/go-cid"
)

func (fvmSimulator *FvmSimulator) SetActor(actorID abi.ActorID, addr address.Address, actor builtin.Actor) {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()

	fvmSimulator.actorsMap[actorID] = actor
	fvmSimulator.addressMap[addr] = actorID
}

func (fvmSimulator *FvmSimulator) LookupDelegatedAddress(actorID abi.ActorID) (address.Address, error) {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()
	for k, v := range fvmSimulator.addressMap {
		if v == actorID {
			return k, nil
		}
	}
	return address.Undef, ferrors.NotFound
}

func (fvmSimulator *FvmSimulator) ResolveAddress(addr address.Address) (abi.ActorID, error) {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()
	if addr.Protocol() == address.ID {
		id, err := address.IDFromAddress(addr)
		if err != nil {
			return 0, err
		}
		return abi.ActorID(id), nil
	}
	id, ok := fvmSimulator.addressMap[addr]
	if !ok {
		return 0, ferrors.NotFound
	}
	return id, nil
}

func (fvmSimulator *FvmSimulator) NextActorAddress() (address.Address, error) {
	seed := time.Now().String()
	return address.NewActorAddress([]byte(seed))
}

// CreateActor this is api can only create builtin actor
func (fvmSimulator *FvmSimulator) CreateActor(actorID abi.ActorID, codeCid cid.Cid) error {
	fvmSimulator.SetActor(actorID, address.Address{}, builtin.Actor{Code: codeCid})
	return nil
}

func (fvmSimulator *FvmSimulator) GetActorCodeCid(addr address.Address) (cid.Cid, error) {
	acstat, err := fvmSimulator.getActorWithAddress(addr)
	if err != nil {
		return cid.Undef, err
	}
	return acstat.Code, nil
}

func (fvmSimulator *FvmSimulator) BalanceOf(actorID abi.ActorID) (*abi.TokenAmount, error) {
	if v, ok := fvmSimulator.actorsMap[actorID]; ok {
		return &v.Balance, nil
	}
	return nil, ferrors.NotFound
}

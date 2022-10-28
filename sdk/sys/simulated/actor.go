package simulated

import (
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v9/migration"
	"github.com/ipfs/go-cid"
)

func (fvmSimulator *FvmSimulator) SetActor(actorID abi.ActorID, addr address.Address, actor migration.Actor) {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()

	fvmSimulator.actorsMap[actorID] = actor
	fvmSimulator.addressMap[addr] = actorID
}

func (fvmSimulator *FvmSimulator) LookupAddress(actorID abi.ActorID) (address.Address, error) {
	for k, v := range fvmSimulator.addressMap {
		if v == actorID {
			return k, nil
		}
	}
	return address.Undef, ErrorNotFound
}

func (fvmSimulator *FvmSimulator) ResolveAddress(addr address.Address) (abi.ActorID, error) {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()
	id, ok := fvmSimulator.addressMap[addr]
	if !ok {
		return 0, ErrorNotFound
	}
	return id, nil
}

func (fvmSimulator *FvmSimulator) NewActorAddress() (address.Address, error) {
	seed := time.Now().String()
	return address.NewActorAddress([]byte(seed))
}

// CreateActor this is api can only create builtin actor
func (fvmSimulator *FvmSimulator) CreateActor(actorID abi.ActorID, codeCid cid.Cid) error {
	fvmSimulator.SetActor(actorID, address.Address{}, migration.Actor{Code: codeCid})
	return nil
}

func (fvmSimulator *FvmSimulator) GetActorCodeCid(addr address.Address, actorID abi.ActorID) (*cid.Cid, error) {
	acstat, err := fvmSimulator.getActorWithAddress(addr)
	if err != nil {
		return nil, err
	}
	return &acstat.Code, nil
}

func (fvmSimulator *FvmSimulator) BalanceOf(addr address.Address, actorID abi.ActorID) (abi.TokenAmount, error) {
	if v, ok := fvmSimulator.actorsMap[actorID]; ok {
		return v.Balance, nil
	}
	return abi.NewTokenAmount(0), ErrorNotFound
}

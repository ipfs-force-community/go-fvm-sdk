package simulated

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
)

func (fvmSimulator *FvmSimulator) SelfRoot() (cid.Cid, error) {
	return fvmSimulator.rootCid, nil
}

func (fvmSimulator *FvmSimulator) SelfSetRoot(id cid.Cid) error {
	fvmSimulator.rootCid = id
	return nil
}

func (fvmSimulator *FvmSimulator) SelfCurrentBalance() (*abi.TokenAmount, error) {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()

	actor, ok := fvmSimulator.actorsMap[fvmSimulator.messageCtx.Caller]
	if !ok {
		return nil, ErrorNotFound
	}
	return &actor.Balance, nil
}

func (fvmSimulator *FvmSimulator) SelfDestruct(addr address.Address) error {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()

	actorId, ok := fvmSimulator.addressMap[addr]
	if !ok {
		return ErrorNotFound
	}
	delete(fvmSimulator.actorsMap, actorId)
	return nil
}

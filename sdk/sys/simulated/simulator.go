package simulated

import (
	"fmt"

	"github.com/filecoin-project/go-state-types/builtin/v9/migration"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func (fvmSimulator *FvmSimulator) GetActor(addr address.Address) (migration.Actor, error) {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()
	actorId, err := fvmSimulator.ResolveAddress(addr) //nolint
	if err != nil {
		return migration.Actor{}, err
	}
	actor, ok := fvmSimulator.actorsMap[actorId]
	if !ok {
		return migration.Actor{}, ErrorNotFound
	}
	return actor, nil
}

func (fvmSimulator *FvmSimulator) ResolveBuiltinActorType(codeCid cid.Cid) (types.ActorType, error) {
	for k, v := range EmbeddedBuiltinActors {
		if v == codeCid {
			av, err := stringToactorType(k)
			return av, err
		}
	}
	return types.ActorType(0), ErrorNotFound
}

func (fvmSimulator *FvmSimulator) GetCodeCidForType(actorT types.ActorType) (cid.Cid, error) {
	actstr, err := actorTypeTostring(actorT)
	if err != nil {
		return cid.Undef, err
	}
	return EmbeddedBuiltinActors[actstr], nil
}

func (fvmSimulator *FvmSimulator) Abort(code uint32, msg string) {
	panic(fmt.Sprintf("%d:%sfvmSimulator", code, msg))
}

func (fvmSimulator *FvmSimulator) Enabled() (bool, error) {
	return true, nil
}

func (fvmSimulator *FvmSimulator) Log(msg string) error {
	fmt.Println(msg)
	return nil
}

func (fvmSimulator *FvmSimulator) SetCallContext(callContext *types.InvocationContext) {
	fvmSimulator.callContext = callContext
}

func (fvmSimulator *FvmSimulator) VMContext() (*types.InvocationContext, error) {
	return fvmSimulator.callContext, nil
}

func (fvmSimulator *FvmSimulator) SetBaseFee(ta abi.TokenAmount) {
	fvmSimulator.baseFee = ta
}

func (fvmSimulator *FvmSimulator) BaseFee() (abi.TokenAmount, error) {
	return fvmSimulator.baseFee, nil
}

func (fvmSimulator *FvmSimulator) Charge(_ string, _ uint64) error {
	return nil
}

func (fvmSimulator *FvmSimulator) SetTotalFilCircSupply(amount abi.TokenAmount) {
	fvmSimulator.totalFilCircSupply = amount
}

func (fvmSimulator *FvmSimulator) TotalFilCircSupply() (abi.TokenAmount, error) {
	return fvmSimulator.totalFilCircSupply, nil
}

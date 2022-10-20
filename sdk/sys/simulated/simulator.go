package simulated

import (
	"fmt"

	"github.com/filecoin-project/go-state-types/builtin/v9/migration"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func (fvmSimulator *FvmSimulator) Open(id cid.Cid) (*types.IpldOpen, error) {
	blockid, blockstat := fvmSimulator.blockOpen(id)
	return &types.IpldOpen{ID: blockid, Size: blockstat.size, Codec: blockstat.codec}, nil
}

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

func (fvmSimulator *FvmSimulator) SelfRoot() (cid.Cid, error) {
	return fvmSimulator.rootCid, nil
}

func (fvmSimulator *FvmSimulator) SelfSetRoot(id cid.Cid) error {
	fvmSimulator.rootCid = id
	return nil
}

func (fvmSimulator *FvmSimulator) SelfCurrentBalance() (abi.TokenAmount, error) {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()

	actor, ok := fvmSimulator.actorsMap[fvmSimulator.callContext.Caller]
	if !ok {
		return abi.TokenAmount{}, ErrorNotFound
	}
	return actor.Balance, nil
}

func (fvmSimulator *FvmSimulator) SelfDestruct(addr address.Address) error {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()

	actorId, ok := fvmSimulator.addressMap[addr] //nolint
	if !ok {
		return ErrorNotFound
	}
	delete(fvmSimulator.actorsMap, actorId)
	return nil
}

func (fvmSimulator *FvmSimulator) Create(codec uint64, data []byte) (uint32, error) {
	index := fvmSimulator.blockCreate(codec, data)
	return index, nil
}

func (fvmSimulator *FvmSimulator) Read(id uint32, offset, size uint32) ([]byte, uint32, error) {
	data, err := fvmSimulator.blockRead(id, offset)
	if err != nil {
		return nil, 0, err
	}
	if size < uint32(len(data)) {
		return data[:size], uint32(len(data)) - size, nil
	}
	return data, 0, nil
}

func (fvmSimulator *FvmSimulator) Stat(id uint32) (*types.IpldStat, error) {
	return fvmSimulator.blockStat(id)
}

func (fvmSimulator *FvmSimulator) BlockLink(id uint32, hashFun uint64, hashLen uint32, cidBuf []byte) (cid.Cid, error) {
	return fvmSimulator.blockLink(id, hashFun, hashLen)
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

func (fvmSimulator *FvmSimulator) GetChainRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	return makeRandomness(dst, round, entropy), nil
}

func (fvmSimulator *FvmSimulator) GetBeaconRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	return makeRandomness(dst, round, entropy), nil
}

func (fvmSimulator *FvmSimulator) SetCallContext(callcontext *types.InvocationContext) {
	fvmSimulator.callContext = callcontext
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

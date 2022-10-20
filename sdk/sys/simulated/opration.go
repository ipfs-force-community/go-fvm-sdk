package simulated

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func (s *FvmSimulator) Open(id cid.Cid) (*types.IpldOpen, error) {
	blockid, blockstat := s.blockOpen(id)
	return &types.IpldOpen{ID: blockid, Size: blockstat.size, Codec: blockstat.codec}, nil
}

func (s *FvmSimulator) SelfRoot() (cid.Cid, error) {
	return s.rootCid, nil
}

func (s *FvmSimulator) SelfSetRoot(id cid.Cid) error {
	s.rootCid = id
	return nil
}

func (s *FvmSimulator) SelfCurrentBalance() (*types.TokenAmount, error) {
	s.actorLk.Lock()
	defer s.actorLk.Unlock()

	actor, ok := s.actorsMap[s.callContext.Caller]
	if !ok {
		return nil, ErrorNotFound
	}
	balance := types.FromBig(&actor.Balance) //todo change TokenAmount to abi
	return &balance, nil
}

func (s *FvmSimulator) SelfDestruct(addr address.Address) error {
	s.actorLk.Lock()
	defer s.actorLk.Unlock()

	actorId, ok := s.addressMap[addr]
	if !ok {
		return ErrorNotFound
	}
	delete(s.actorsMap, actorId)
	return nil
}

func (s *FvmSimulator) Create(codec uint64, data []byte) (uint32, error) {
	index := s.blockCreate(codec, data)
	return uint32(index), nil
}

func (s *FvmSimulator) Read(id uint32, offset, size uint32) ([]byte, uint32, error) {
	data, err := s.blockRead(id, offset)
	return data, 0, err
}

func (s *FvmSimulator) Stat(id uint32) (*types.IpldStat, error) {
	return s.blockStat(id)
}

func (s *FvmSimulator) BlockLink(id uint32, hashFun uint64, hashLen uint32, cidBuf []byte) (cided cid.Cid, err error) {
	return s.blockLink(id, hashFun, hashLen)
}

func (s *FvmSimulator) ResolveBuiltinActorType(codeCid cid.Cid) (types.ActorType, error) {
	for k, v := range EmbeddedBuiltinActors {
		if v == codeCid {
			av, err := stringToactorType(k)
			return av, err
		}
	}
	return types.ActorType(0), ErrorNotFound
}

func (s *FvmSimulator) GetCodeCidForType(actorT types.ActorType) (cid.Cid, error) {
	actstr, err := actorTypeTostring(actorT)
	if err != nil {
		return cid.Undef, err
	}
	return EmbeddedBuiltinActors[actstr], nil
}

func (s *FvmSimulator) Abort(code uint32, msg string) {
	panic(fmt.Sprintf("%d:%s", code, msg))
}

func (s *FvmSimulator) Enabled() (bool, error) {
	return true, nil
}

func (s *FvmSimulator) Log(msg string) error {
	fmt.Println(msg)
	return nil
}

func (s *FvmSimulator) GetChainRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	h := makeRandomness(dst, round, entropy)
	return abi.Randomness(h), nil
}

func (s *FvmSimulator) GetBeaconRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	h := makeRandomness(dst, round, entropy)
	return abi.Randomness(h), nil
}

func (s *FvmSimulator) SetCallContext(callcontext *types.InvocationContext) {
	s.callContext = callcontext
}

func (s *FvmSimulator) VMContext() (*types.InvocationContext, error) {
	return s.callContext, nil
}

func (s *FvmSimulator) SetBaseFee(ta big.Int) {
	amount, _ := types.FromString(ta.String())
	s.baseFee = amount
}

func (s *FvmSimulator) BaseFee() (*types.TokenAmount, error) {
	return &s.baseFee, nil
}

func (s *FvmSimulator) Charge(name string, compute uint64) error {
	return nil
}

func (s *FvmSimulator) SetTotalFilCircSupply(ta big.Int) {
	amount, _ := types.FromString(ta.String())
	s.totalFilCircSupply = amount
}

func (s *FvmSimulator) TotalFilCircSupply() (*types.TokenAmount, error) {
	return &s.totalFilCircSupply, nil
}

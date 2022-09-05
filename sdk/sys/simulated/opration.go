package simulated

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/specs-actors/v7/actors/runtime"
	"github.com/filecoin-project/specs-actors/v7/actors/runtime/proof"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func (s *Fsm) Open(id cid.Cid) (*types.IpldOpen, error) {

	result := new(types.IpldOpen)
	blockid, blockstat := s.blockOpen(id)
	result.ID = blockid
	result.Size = blockstat.size
	result.Codec = blockstat.codec
	return result, nil
}

func (s *Fsm) SelfRoot() (cid.Cid, error) {
	return s.rootCid, nil
}

func (s *Fsm) SelfSetRoot(id cid.Cid) error {
	s.rootCid = id
	return nil
}

func (s *Fsm) SelfCurrentBalance() (*types.TokenAmount, error) {
	return s.currentBalance, nil
}

func (s *Fsm) SelfDestruct(addr address.Address) error {
	s.actorMutex.Lock()
	defer s.actorMutex.Unlock()
	actorid, ok := s.address.Load(addr)
	if !ok {
		return ErrorNotFound
	}
	s.actors.Delete(actorid)
	return nil

}

func (s *Fsm) Create(codec uint64, data []byte) (uint32, error) {
	index := s.blockCreate(codec, data)
	return uint32(index), nil
}

func (s *Fsm) Read(id uint32, offset, size uint32) ([]byte, uint32, error) {

	data, err := s.blockRead(id, offset)
	return data, 0, err
}

func (s *Fsm) Stat(id uint32) (*types.IpldStat, error) {
	return s.blockStat(id)
}
func (s *Fsm) BlockLink(id uint32, hashFun uint64, hashLen uint32, cidBuf []byte) (cided cid.Cid, err error) {
	return s.blockLink(id, hashFun, hashLen)
}
func (s *Fsm) ResolveAddress(addr address.Address) (abi.ActorID, error) {

	id, ok := s.address.Load(addr)
	if !ok {
		return 0, ErrorNotFound
	}
	idu32 := id.(uint32)
	return abi.ActorID(uint32(idu32)), nil
}
func (s *Fsm) NewActorAddress() (address.Address, error) {
	return address.NewActorAddress([]byte{})
}

func (s *Fsm) GetActorCodeCid(addr address.Address) (*cid.Cid, error) {
	acstat, err := s.getActorWithAddress(addr)
	if err != nil {
		return nil, err
	}
	return &acstat.Code, nil
}
func (s *Fsm) ResolveBuiltinActorType(codeCid cid.Cid) (types.ActorType, error) {
	for k, v := range EmbeddedBuiltinActors {
		if v == codeCid {
			av, err := stringToactorType(k)
			return av, err
		}
	}
	return types.ActorType(0), ErrorNotFound
}

func (s *Fsm) GetCodeCidForType(actorT types.ActorType) (cid.Cid, error) {
	actstr, err := actorTypeTostring(actorT)
	if err != nil {
		return cid.Undef, err
	}
	cid := EmbeddedBuiltinActors[actstr]
	return cid, nil
}

func (s *Fsm) CreateActor(actorID abi.ActorID, codeCid cid.Cid) error {

	SetActorAndAddress(uint32(actorID), ActorState{Code: codeCid}, address.Address{})
	return nil
}

func (s *Fsm) Abort(code uint32, msg string) {
	panic(fmt.Sprintf("%v:%v", code, msg))
}

func (s *Fsm) VerifySignature(
	signature *crypto.Signature,
	signer *address.Address,
	plaintext []byte,
) (bool, error) {
	panic("This is not implement")
}

func (s *Fsm) HashBlake2b(data []byte) ([32]byte, error) {
	result := blakehash(data)
	var temp [32]byte

	copy(temp[:], result[:32])
	return temp, nil
}
func (s *Fsm) ComputeUnsealedSectorCid(
	proofType abi.RegisteredSealProof,
	pieces []abi.PieceInfo,
) (cid.Cid, error) {
	panic("This is not implement")
}
func (s *Fsm) VerifySeal(info *proof.SealVerifyInfo) (bool, error) {
	panic("This is not implement")
}
func (s *Fsm) VerifyPost(info *proof.WindowPoStVerifyInfo) (bool, error) {
	panic("This is not implement")
}
func (s *Fsm) VerifyConsensusFault(
	h1 []byte,
	h2 []byte,
	extra []byte,
) (*runtime.ConsensusFault, error) {
	panic("This is not implement")
}
func (s *Fsm) VerifyAggregateSeals(info *types.AggregateSealVerifyProofAndInfos) (bool, error) {
	panic("This is not implement")
}
func (s *Fsm) VerifyReplicaUpdate(info *types.ReplicaUpdateInfo) (bool, error) {
	panic("This is not implement")
}
func (s *Fsm) BatchVerifySeals(sealVerifyInfos []proof.SealVerifyInfo) ([]bool, error) {
	panic("This is not implement")
}

func (s *Fsm) VMContext() (*types.InvocationContext, error) {
	return &types.InvocationContext{}, nil
}

func (s *Fsm) Enabled() (bool, error) {
	return true, nil
}
func (s *Fsm) Log(msg string) error {
	fmt.Println(msg)
	return nil
}

func (s *Fsm) Send(to address.Address, method uint64, params uint32, value types.TokenAmount) (*types.Send, error) {
	panic("This is not implement")
}

func (s *Fsm) GetChainRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	h := BeaconRandomness(dst, round, entropy)
	return abi.Randomness(h), nil
}
func (s *Fsm) GetBeaconRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	h := BeaconRandomness(dst, round, entropy)
	return abi.Randomness(h), nil
}

func (s *Fsm) BaseFee() (*types.TokenAmount, error) {
	return &types.TokenAmount{Lo: 0, Hi: 0}, nil
}
func (s *Fsm) TotalFilCircSupply() (*types.TokenAmount, error) {
	return &types.TokenAmount{Lo: 0, Hi: 0}, nil
}

func (s *Fsm) Charge(name string, compute uint64) error {
	return nil
}

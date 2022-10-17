package simulated

import (
	"fmt"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v9/migration"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/specs-actors/v7/actors/runtime"
	"github.com/filecoin-project/specs-actors/v7/actors/runtime/proof"
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
	return &s.baseFee, nil
}

func (s *FvmSimulator) SelfDestruct(addr address.Address) error {
	s.actorMutex.Lock()
	defer s.actorMutex.Unlock()

	actorid, ok := s.addressMap.Load(addr)
	if !ok {
		return ErrorNotFound
	}
	s.actorsMap.Delete(actorid)
	return nil
}

// nolint
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

func (s *FvmSimulator) ResolveAddress(addr address.Address) (abi.ActorID, error) {
	id, ok := s.addressMap.Load(addr)
	if !ok {
		return 0, ErrorNotFound
	}
	idu32, ok := id.(uint32)
	if !ok {
		return abi.ActorID(0), ErrorKeyTypeException
	}
	return abi.ActorID(idu32), nil
}

func (s *FvmSimulator) NewActorAddress() (address.Address, error) {
	seed := time.Now().String()
	return address.NewActorAddress([]byte(seed))
}

func (s *FvmSimulator) GetActorCodeCid(addr address.Address) (*cid.Cid, error) {
	acstat, err := s.getActorWithAddress(addr)
	if err != nil {
		return nil, err
	}
	return &acstat.Code, nil
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

func (s *FvmSimulator) CreateActor(actorID abi.ActorID, codeCid cid.Cid) error {
	s.SetActorAndAddress(uint32(actorID), migration.Actor{Code: codeCid}, address.Address{})
	return nil
}

// nolint
func (s *FvmSimulator) Abort(code uint32, msg string) {
	panic(fmt.Sprintf("%d:%s", code, msg))
}

func (s *FvmSimulator) VerifySignature(
	signature *crypto.Signature,
	signer *address.Address,
	plaintext []byte,
) (bool, error) {
	panic("This is not implement")
}

func (s *FvmSimulator) HashBlake2b(data []byte) ([32]byte, error) {
	result := blakehash(data)
	var temp [32]byte
	copy(temp[:], result[:32])
	return temp, nil
}

func (s *FvmSimulator) ComputeUnsealedSectorCid(
	proofType abi.RegisteredSealProof,
	pieces []abi.PieceInfo,
) (cid.Cid, error) {
	panic("This is not implement")
}

func (s *FvmSimulator) VerifySeal(info *proof.SealVerifyInfo) (bool, error) {
	panic("This is not implement")
}

func (s *FvmSimulator) VerifyPost(info *proof.WindowPoStVerifyInfo) (bool, error) {
	panic("This is not implement")
}

func (s *FvmSimulator) VerifyConsensusFault(
	h1 []byte,
	h2 []byte,
	extra []byte,
) (*runtime.ConsensusFault, error) {
	panic("This is not implement")
}

func (s *FvmSimulator) VerifyAggregateSeals(info *types.AggregateSealVerifyProofAndInfos) (bool, error) {
	panic("This is not implement")
}

func (s *FvmSimulator) VerifyReplicaUpdate(info *types.ReplicaUpdateInfo) (bool, error) {
	panic("This is not implement")
}
func (s *FvmSimulator) BatchVerifySeals(sealVerifyInfos []proof.SealVerifyInfo) ([]bool, error) {
	panic("This is not implement")
}

func (s *FvmSimulator) VMContext() (*types.InvocationContext, error) {
	return s.callContext, nil
}

func (s *FvmSimulator) Enabled() (bool, error) {
	return true, nil
}

func (s *FvmSimulator) Log(msg string) error {
	fmt.Println(msg)
	return nil
}

func (s *FvmSimulator) Send(to address.Address, method uint64, params uint32, value types.TokenAmount) (*types.Send, error) {
	send, ok := s.sendMatch(to, method, params, *value.Big())
	if ok {
		return send, nil
	}
	return nil, ErrorKeyMatchFail
}

func (s *FvmSimulator) GetChainRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	h := makeRandomness(dst, round, entropy)
	return abi.Randomness(h), nil
}

func (s *FvmSimulator) GetBeaconRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error) {
	h := makeRandomness(dst, round, entropy)
	return abi.Randomness(h), nil
}

func (s *FvmSimulator) BaseFee() (*types.TokenAmount, error) {
	return &s.baseFee, nil
}

func (s *FvmSimulator) TotalFilCircSupply() (*types.TokenAmount, error) {
	return &s.totalFilCircSupply, nil
}

func (s *FvmSimulator) Charge(name string, compute uint64) error {
	return nil
}

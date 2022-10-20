//nolint:unparam
package simulated

import (
	"context"
	"sync"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v9/migration"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

// nolint
type IpldOpen struct {
	codec uint64
	id    uint32
	size  uint32
}

type block struct {
	codec uint64
	data  []byte
}
type BlockStat struct {
	codec uint64
	size  uint32
}

func (s *block) stat() BlockStat {
	return BlockStat{codec: s.codec, size: uint32(len(s.data))}
}

type blocks []block

// CreateSimulateEnv new context of simulated
func CreateSimulateEnv(callContext *types.InvocationContext, baseFee big.Int, totalFilCircSupply big.Int) (*FvmSimulator, context.Context) {
	fsm := NewFvmSimulator(callContext, baseFee, totalFilCircSupply)
	return fsm, fsm.Context
}

type FvmSimulator struct {
	Context     context.Context
	blocksMutex sync.Mutex
	blocks      blocks
	blockid     uint32
	ipld        sync.Map
	actorLk     sync.Mutex
	// actorid->ActorState
	actorsMap map[abi.ActorID]migration.Actor
	// address->actorid
	addressMap map[address.Address]abi.ActorID

	callContext        *types.InvocationContext
	rootCid            cid.Cid
	baseFee            types.TokenAmount
	totalFilCircSupply types.TokenAmount
	currentBalance     types.TokenAmount
	sendList           []SendMock
}

func NewFvmSimulator(callContext *types.InvocationContext, baseFee big.Int, totalFilCircSupply big.Int) *FvmSimulator {
	fsm := &FvmSimulator{
		blockid:            1,
		callContext:        callContext,
		baseFee:            types.FromBig(&baseFee),
		totalFilCircSupply: types.FromBig(&totalFilCircSupply),
	}
	fsm.Context = context.WithValue(context.Background(), types.SimulatedEnvkey, fsm)
	return fsm
}

func (a *FvmSimulator) sendMatch(to address.Address, method uint64, params uint32, value big.Int) (*types.Send, bool) {
	for i, v := range a.sendList {
		if to != v.to {
			continue
		}
		if method != v.method {
			continue
		}
		if params != v.params {
			continue
		}
		if !value.Equals(v.value) {
			continue
		}
		if i == len(a.sendList)-1 {
			a.sendList = a.sendList[0 : i-1]
		} else {
			a.sendList = append(a.sendList[:i], a.sendList[i+1:]...)
		}

		return &v.out, true
	}
	return nil, false
}

func (s *FvmSimulator) blockLink(blockid uint32, hashfun uint64, hashlen uint32) (blkCid cid.Cid, err error) {
	block, err := s.getBlock(blockid)
	if err != nil {
		return cid.Undef, err
	}

	Mult, _ := mh.Sum(block.data, hashfun, int(hashlen))

	blkCid = cid.NewCidV1(block.codec, Mult)
	s.putData(blkCid, block.data)
	return
}

func (s *FvmSimulator) blockCreate(codec uint64, data []byte) uint32 {
	s.putBlock(block{codec: codec, data: data})
	return uint32(len(s.blocks) - 1)
}

func (s *FvmSimulator) blockOpen(id cid.Cid) (blockID uint32, blockStat BlockStat) {
	data, _ := s.getData(id)
	block := block{data: data, codec: id.Prefix().GetCodec()}

	stat := block.stat()
	bid := s.putBlock(block)
	return bid, stat

}

func (s *FvmSimulator) blockRead(id uint32, offset uint32) ([]byte, error) {
	block, err := s.getBlock(id)
	if err != nil {
		return nil, err
	}
	data := block.data

	if offset >= uint32(len(data)) {
		return nil, ErrorIDValid
	}
	return data[offset:], nil
}

func (s *FvmSimulator) blockStat(blockID uint32) (*types.IpldStat, error) {
	b, err := s.getBlock(blockID)
	if err != nil {
		return nil, ErrorNotFound
	}
	return &types.IpldStat{Size: b.stat().size, Codec: b.codec}, ErrorNotFound
}

func (s *FvmSimulator) putData(key cid.Cid, value []byte) {
	s.ipld.Store(key, value)
}

func (s *FvmSimulator) getData(key cid.Cid) ([]byte, error) {
	value, ok := s.ipld.Load(key)
	if ok {
		return value.([]byte), nil
	}
	return nil, ErrorNotFound
}

func (s *FvmSimulator) putBlock(block block) uint32 {
	s.blocksMutex.Lock()
	defer s.blocksMutex.Unlock()

	s.blocks = append(s.blocks, block)

	return uint32(len(s.blocks) - 1)
}

func (s *FvmSimulator) getBlock(blockID uint32) (block, error) {
	s.blocksMutex.Lock()
	defer s.blocksMutex.Unlock()

	if blockID >= uint32(len(s.blocks)) {
		return block{}, ErrorNotFound
	}
	return s.blocks[blockID], nil
}

// nolint
func (s *FvmSimulator) putActor(actorID abi.ActorID, actor migration.Actor) error {
	_, err := s.getActorWithActorid(actorID)
	if err == nil {
		return ErrorKeyExists
	}
	s.ipld.Store(actorID, actor)
	return nil
}

// nolint
func (s *FvmSimulator) getActorWithActorid(actorID abi.ActorID) (migration.Actor, error) {
	s.actorLk.Lock()
	defer s.actorLk.Unlock()

	actor, ok := s.actorsMap[actorID]
	if ok {
		return actor, nil
	}
	return migration.Actor{}, ErrorNotFound
}

// nolint
func (s *FvmSimulator) getActorWithAddress(addr address.Address) (migration.Actor, error) {
	s.actorLk.Lock()
	defer s.actorLk.Unlock()

	actorId, ok := s.addressMap[addr]
	if ok {
		return migration.Actor{}, ErrorNotFound
	}
	as, ok := s.actorsMap[actorId]
	if !ok {
		return migration.Actor{}, ErrorNotFound
	}
	return as, nil
}

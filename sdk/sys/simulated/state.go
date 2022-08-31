package simulated

import (
	"sync"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

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

type Actor struct {
	ActorID    uint64
	ActorState ActorState
}

type ActorState struct {
	code     cid.Cid
	state    cid.Cid
	sequence uint64
	balance  big.Int
}

func newActorState(code_id cid.Cid) ActorState {
	Mult, _ := mh.Sum([]byte{}, mh.BLAKE2B_MAX, 32)
	return ActorState{code: code_id, sequence: 0, balance: big.NewInt(0), state: cid.NewCidV1(cid.DagCBOR, Mult)}
}

var DefaultFsm *Fsm

func init() {
	Begin()
}

func Begin() {
	DefaultFsm = newSate()
}
func End() {
	DefaultFsm = newSate()
}

type Fsm struct {
	blocksMutex sync.Mutex
	blocks
	blockid    uint32
	Ipld       sync.Map
	actorMutex sync.Mutex
	actors     sync.Map
	address    sync.Map

	rootCid            cid.Cid
	baseFee            *types.TokenAmount
	totalFilCircSupply *types.TokenAmount
	currentBalance     *types.TokenAmount
}

func newSate() *Fsm {
	return &Fsm{blockid: 1, Ipld: sync.Map{}}
}

func (s *Fsm) blockLink(blockid uint32, hash_fun uint64, hash_len uint32) (cided cid.Cid, err error) {
	block, err := s.getBlock((blockid))
	if err != nil {
		return cid.Undef, err
	}
	Mult, _ := mh.Sum(block.data, hash_fun, int(hash_len))
	cided = cid.NewCidV1(block.codec, Mult)
	s.putData(cided, block.data)
	return
}

func (s *Fsm) blockCreate(codec uint64, data []byte) uint32 {
	s.putBlock(block{codec: codec, data: data})
	return uint32(len(s.blocks) - 1)
}

func (s *Fsm) blockOpen(id cid.Cid) (BlockId uint32, BlockStat BlockStat) {
	data, _ := s.getData(id)
	block := block{data: data, codec: id.Prefix().GetCodec()}

	stat := block.stat()
	bid := s.putBlock(block)
	return bid, stat

}

func (s *Fsm) blockRead(id uint32, offset uint32) ([]byte, error) {
	block, err := s.getBlock(id)
	if err != nil {
		return nil, err
	}
	data := block.data

	if offset >= uint32(len(data)) {
		return nil, ErrorIdValid
	}
	return data[offset:], nil
}

func (s *Fsm) blockStat(blockId uint32) (*types.IpldStat, error) {
	b, err := s.getBlock(blockId)
	if err != nil {
		return nil, ErrorNotFound
	}
	return &types.IpldStat{Size: b.stat().size, Codec: b.codec}, ErrorNotFound
}

func (s *Fsm) putData(key cid.Cid, value []byte) {
	s.Ipld.Store(key, value)
}

func (s *Fsm) getData(key cid.Cid) ([]byte, error) {
	value, ok := s.Ipld.Load(key)
	if ok {
		return value.([]byte), nil
	}
	return nil, ErrorNotFound
}

func (s *Fsm) putBlock(block block) uint32 {
	s.blocksMutex.Lock()
	defer s.blocksMutex.Unlock()

	s.blocks = append(s.blocks, block)
	return uint32(len(s.blocks) - 1)
}

func (s *Fsm) getBlock(blockId uint32) (block, error) {
	s.blocksMutex.Lock()
	defer s.blocksMutex.Unlock()

	if blockId >= uint32(len(s.blocks)) {
		return block{}, ErrorNotFound
	}
	return s.blocks[blockId], nil
}

func (s *Fsm) putActor(actorID uint64, actor ActorState) error {
	_, err := s.getActorWithActorid(uint32(actorID))
	if err == nil {
		return ErrorKeyExists
	}
	s.Ipld.Store(actorID, actor)
	return nil
}

func (s *Fsm) getActorWithActorid(actorID uint32) (ActorState, error) {
	actor, ok := s.actors.Load(actorID)
	if ok {
		return actor.(ActorState), nil
	}
	return ActorState{}, ErrorNotFound
}

func (s *Fsm) getActorWithAddress(addr address.Address) (ActorState, error) {
	s.actorMutex.Lock()
	defer s.actorMutex.Unlock()

	actorid, ok := s.address.Load(addr)
	if ok {
		return ActorState{}, ErrorNotFound
	}
	as, ok := s.actors.Load(actorid)
	if !ok {
		return ActorState{}, ErrorNotFound
	}
	return as.(ActorState), nil
}

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
	baseFee            abi.TokenAmount
	totalFilCircSupply abi.TokenAmount
	sendList           []SendMock
}

func NewFvmSimulator(callContext *types.InvocationContext, baseFee abi.TokenAmount, totalFilCircSupply abi.TokenAmount) *FvmSimulator {
	fsm := &FvmSimulator{
		blockid:            1,
		callContext:        callContext,
		baseFee:            baseFee,
		totalFilCircSupply: totalFilCircSupply,
		actorsMap:          make(map[abi.ActorID]migration.Actor),
		addressMap:         make(map[address.Address]abi.ActorID),
	}
	fsm.Context = context.WithValue(context.Background(), types.SimulatedEnvkey, fsm)
	return fsm
}

func (fvmSimulator *FvmSimulator) sendMatch(to address.Address, method uint64, params uint32, value big.Int) (*types.Send, bool) {
	for i, v := range fvmSimulator.sendList {
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
		if i == len(fvmSimulator.sendList)-1 {
			fvmSimulator.sendList = fvmSimulator.sendList[0 : i-1]
		} else {
			fvmSimulator.sendList = append(fvmSimulator.sendList[:i], fvmSimulator.sendList[i+1:]...)
		}

		return &v.out, true
	}
	return nil, false
}

func (fvmSimulator *FvmSimulator) blockLink(blockid uint32, hashfun uint64, hashlen uint32) (blkCid cid.Cid, err error) {
	block, err := fvmSimulator.getBlock(blockid)
	if err != nil {
		return cid.Undef, err
	}

	Mult, _ := mh.Sum(block.data, hashfun, int(hashlen))

	blkCid = cid.NewCidV1(block.codec, Mult)
	fvmSimulator.putData(blkCid, block.data)
	return
}

func (fvmSimulator *FvmSimulator) blockCreate(codec uint64, data []byte) uint32 {
	fvmSimulator.putBlock(block{codec: codec, data: data})
	return uint32(len(fvmSimulator.blocks) - 1)
}

func (fvmSimulator *FvmSimulator) blockOpen(id cid.Cid) (blockID uint32, blockStat BlockStat) {
	data, _ := fvmSimulator.getData(id)
	block := block{data: data, codec: id.Prefix().GetCodec()}

	stat := block.stat()
	bid := fvmSimulator.putBlock(block)
	return bid, stat

}

func (fvmSimulator *FvmSimulator) blockRead(id uint32, offset uint32) ([]byte, error) {
	block, err := fvmSimulator.getBlock(id)
	if err != nil {
		return nil, err
	}
	data := block.data

	if offset >= uint32(len(data)) {
		return nil, ErrorIDValid
	}
	return data[offset:], nil
}

func (fvmSimulator *FvmSimulator) blockStat(blockID uint32) (*types.IpldStat, error) {
	b, err := fvmSimulator.getBlock(blockID)
	if err != nil {
		return nil, ErrorNotFound
	}
	return &types.IpldStat{Size: b.stat().size, Codec: b.codec}, ErrorNotFound
}

func (fvmSimulator *FvmSimulator) putData(key cid.Cid, value []byte) {
	fvmSimulator.ipld.Store(key, value)
}

func (fvmSimulator *FvmSimulator) getData(key cid.Cid) ([]byte, error) {
	value, ok := fvmSimulator.ipld.Load(key)
	if ok {
		return value.([]byte), nil
	}
	return nil, ErrorNotFound
}

func (fvmSimulator *FvmSimulator) putBlock(block block) uint32 {
	fvmSimulator.blocksMutex.Lock()
	defer fvmSimulator.blocksMutex.Unlock()

	fvmSimulator.blocks = append(fvmSimulator.blocks, block)

	return uint32(len(fvmSimulator.blocks) - 1)
}

func (fvmSimulator *FvmSimulator) getBlock(blockID uint32) (block, error) {
	fvmSimulator.blocksMutex.Lock()
	defer fvmSimulator.blocksMutex.Unlock()

	if blockID >= uint32(len(fvmSimulator.blocks)) {
		return block{}, ErrorNotFound
	}
	return fvmSimulator.blocks[blockID], nil
}

// nolint
func (fvmSimulator *FvmSimulator) putActor(actorID abi.ActorID, actor migration.Actor) error {
	_, err := fvmSimulator.getActorWithActorid(actorID)
	if err == nil {
		return ErrorKeyExists
	}
	fvmSimulator.ipld.Store(actorID, actor)
	return nil
}

// nolint
func (fvmSimulator *FvmSimulator) getActorWithActorid(actorID abi.ActorID) (migration.Actor, error) {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()

	actor, ok := fvmSimulator.actorsMap[actorID]
	if ok {
		return actor, nil
	}
	return migration.Actor{}, ErrorNotFound
}

// nolint
func (fvmSimulator *FvmSimulator) getActorWithAddress(addr address.Address) (migration.Actor, error) {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()

	actorId, ok := fvmSimulator.addressMap[addr]
	if ok {
		return migration.Actor{}, ErrorNotFound
	}
	as, ok := fvmSimulator.actorsMap[actorId]
	if !ok {
		return migration.Actor{}, ErrorNotFound
	}
	return as, nil
}

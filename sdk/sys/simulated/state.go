//nolint:unparam
package simulated

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/filecoin-project/go-state-types/builtin"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
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

// CreateEmptySimulator new context of simulated
func CreateEmptySimulator() (*FvmSimulator, context.Context) {
	fsm := NewFvmSimulator(&types.MessageContext{}, &types.NetworkContext{}, big.Zero())
	return fsm, fsm.Context
}

// CreateSimulateEnv new context of simulated
func CreateSimulateEnv(callContext *types.MessageContext, networkContext *types.NetworkContext, totalFilCircSupply big.Int) (*FvmSimulator, context.Context) {
	fsm := NewFvmSimulator(callContext, networkContext, totalFilCircSupply)
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
	actorsMap map[abi.ActorID]builtin.Actor
	// address->actorid
	addressMap map[address.Address]abi.ActorID

	messageCtx         *types.MessageContext
	networkCtx         *types.NetworkContext
	rootCid            cid.Cid
	tipsetCidLk        sync.Mutex
	tipsetCids         map[abi.ChainEpoch]*cid.Cid
	totalFilCircSupply abi.TokenAmount
	sendList           []SendMock
	events             []types.ActorEvent
}

func NewFvmSimulator(callContext *types.MessageContext, networkContext *types.NetworkContext, totalFilCircSupply abi.TokenAmount) *FvmSimulator {
	fsm := &FvmSimulator{
		blockid:            1,
		messageCtx:         callContext,
		totalFilCircSupply: totalFilCircSupply,
		actorsMap:          make(map[abi.ActorID]builtin.Actor),
		addressMap:         make(map[address.Address]abi.ActorID),
	}
	fsm.Context = context.WithValue(context.Background(), types.SimulatedEnvkey, fsm)
	return fsm
}

func (fvmSimulator *FvmSimulator) sendMatch(to address.Address, method abi.MethodNum, paramsId uint32, value big.Int) (*types.SendResult, error) {
	rawParams, err := fvmSimulator.getBlock(paramsId)
	if err != nil {
		return nil, err
	}

	if len(fvmSimulator.sendList) == 0 {
		return nil, fmt.Errorf("no expect send for(to: %s method: %d params %v value %s", to, method, rawParams, value)
	}
	send := fvmSimulator.sendList[0]

	if to != send.To {
		return nil, fmt.Errorf("send to not match expect: %s actual %s", send.To, to)
	}
	if method != send.Method {
		return nil, fmt.Errorf("send method not match expect: %d actual %d", send.Method, method)
	}
	if !bytes.Equal(rawParams.data, send.Params) {
		return nil, fmt.Errorf("send to not match expect: %v actual %v", send.Params, rawParams.data)
	}

	if !value.Equals(send.Value) {
		return nil, fmt.Errorf("send value not match expect: %s actual %s", send.Value, value)
	}

	fvmSimulator.sendList = fvmSimulator.sendList[1:]
	return &send.Out, nil
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
func (fvmSimulator *FvmSimulator) putActor(actorID abi.ActorID, actor builtin.Actor) error {
	_, err := fvmSimulator.getActorWithActorid(actorID)
	if err == nil {
		return ErrorKeyExists
	}
	fvmSimulator.ipld.Store(actorID, actor)
	return nil
}

// nolint
func (fvmSimulator *FvmSimulator) getActorWithActorid(actorID abi.ActorID) (builtin.Actor, error) {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()

	actor, ok := fvmSimulator.actorsMap[actorID]
	if ok {
		return actor, nil
	}
	return builtin.Actor{}, ErrorNotFound
}

// nolint
func (fvmSimulator *FvmSimulator) getActorWithAddress(addr address.Address) (builtin.Actor, error) {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()

	actorId, ok := fvmSimulator.addressMap[addr]
	if !ok {
		if SimulateDebug {
			for addr, _ := range fvmSimulator.addressMap {
				fmt.Println("Has:", addr)
			}
			fmt.Println("not found", addr)
		}
		return builtin.Actor{}, ErrorNotFound
	}
	as, ok := fvmSimulator.actorsMap[actorId]
	if !ok {
		return builtin.Actor{}, ErrorNotFound
	}
	return as, nil
}

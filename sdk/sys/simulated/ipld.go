package simulated

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func (fvmSimulator *FvmSimulator) Open(id cid.Cid) (*types.IpldOpen, error) {
	blockid, blockstat := fvmSimulator.blockOpen(id)
	return &types.IpldOpen{ID: blockid, Size: blockstat.size, Codec: blockstat.codec}, nil
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

func (fvmSimulator *FvmSimulator) BlockLink(id uint32, hashFun uint64, hashLen uint32) (cid.Cid, error) {
	return fvmSimulator.blockLink(id, hashFun, hashLen)
}

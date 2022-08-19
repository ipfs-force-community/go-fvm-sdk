//go:build simulate

package sys

import (

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/internal/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func Open(id cid.Cid) (*types.IpldOpen, error) {
	
	return simulated.MockFvmInstance.Open(id)
}

func Create(codec uint64, data []byte) (uint32, error) {
	return simulated.MockFvmInstance.Create(codec, data)
}

func Read(id uint32, offset uint32, buf []byte) (uint32, error) {
	return simulated.MockFvmInstance.Read(id, offset, buf)
}

func Stat(id uint32) (*types.IpldStat, error) {
	return simulated.MockFvmInstance.Stat(id)
}

func BlockLink(id uint32, hashFun uint64, hashLen uint32, cidBuf []byte) (uint32, error) {
	return simulated.MockFvmInstance.BlockLink(id, hashFun, hashLen, cidBuf)
}

//go:build simulate
// +build simulate

package sys

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func Open(ctx context.Context, id cid.Cid) (*types.IpldOpen, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.Open(id)
	}
	panic(ErrorEnvValid)
}

func Create(ctx context.Context, codec uint64, data []byte) (uint32, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		a, v := env.Create(codec, data)
		return a, v
	}
	panic(ErrorEnvValid)
}

func Read(ctx context.Context, id uint32, offset, size uint32) ([]byte, uint32, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.Read(id, offset, size)
	}
	panic(ErrorEnvValid)
}

func Stat(ctx context.Context, id uint32) (*types.IpldStat, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.Stat(id)
	}
	panic(ErrorEnvValid)
}

func BlockLink(ctx context.Context, id uint32, hashFun uint64, hashLen uint32, cidBuf []byte) (cid.Cid, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.BlockLink(id, hashFun, hashLen, cidBuf)
	}
	panic(ErrorEnvValid)
}

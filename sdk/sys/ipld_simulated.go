//go:build simulated
// +build simulated

package sys

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func Open(ctx context.Context, id cid.Cid) (*types.IpldOpen, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.Open(id)
	}
	return &types.IpldOpen{}, nil
}

func Create(ctx context.Context, codec uint64, data []byte) (uint32, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		a, v := env.Create(codec, data)
		return a, v
	}
	return 0, nil
}

func Read(ctx context.Context, id uint32, offset, size uint32) ([]byte, uint32, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.Read(id, offset, size)
	}
	return []byte{}, 0, nil
}

func Stat(ctx context.Context, id uint32) (*types.IpldStat, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.Stat(id)
	}
	return &types.IpldStat{}, nil

}

func BlockLink(ctx context.Context, id uint32, hashFun uint64, hashLen uint32, cidBuf []byte) (cid.Cid, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.BlockLink(id, hashFun, hashLen, cidBuf)
	}

	return cid.Undef, nil

}



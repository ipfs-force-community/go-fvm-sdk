package sdk

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

// Put store a block. The block will only be persisted in the state-tree if the CID is "linked in" to
// the actor's state-tree before the end of the current invocation.
func Put(ctx context.Context, mhCode uint64, mhSize uint32, codec uint64, data []byte) (cid.Cid, error) {
	id, err := sys.Create(ctx, codec, data)
	if err != nil {
		return cid.Undef, err
	}

	// I really hate this CID interface. Why can't I just have bytes?
	return sys.BlockLink(ctx, id, mhCode, mhSize)
}

// Get get a block. It's valid to call this on:
//
// 1. All CIDs returned by prior calls to `get_root`...
// 2. All CIDs returned by prior calls to `put`...
// 3. Any children of a blocks returned by prior calls to `get`...
//
// ...during the current invocation.
func Get(ctx context.Context, cid cid.Cid) ([]byte, error) {
	// TODO: Check length of cid?
	result, err := sys.Open(ctx, cid)
	if err != nil {
		return nil, err
	}

	return GetBlock(ctx, result.ID, &result.Size)
}

// GetBlock gets the data of the block referenced by BlockId. If the caller knows the size, this function
// will read the block in a single syscall. Otherwise, any block over 1KiB will take two syscalls.
func GetBlock(ctx context.Context, id types.BlockID, size *uint32) ([]byte, error) {
	if id == types.UNIT {
		return []byte{}, nil
	}

	var size1 uint32
	if size != nil {
		size1 = *size
	} else {
		size1 = 1024
	}

	// block := make([]byte, size1)
	block, remaining, err := sys.Read(ctx, id, 0, size1) //only set len and slice
	if err != nil {
		return nil, err
	}

	if remaining > 0 { //more than 1KiB
		sencondPart, remaining, err := sys.Read(ctx, id, uint32(len(block)), remaining) //only set len and slice
		if err != nil {
			return nil, err
		}
		if remaining > 0 {
			panic("should have read whole block")
		}
		block = append(block, sencondPart...)
	}
	return block, nil
}

// PutBlock writes the supplied block and returns the BlockId.
func PutBlock(ctx context.Context, codec types.Codec, data []byte) (types.BlockID, error) {
	return sys.Create(ctx, codec, data)
}

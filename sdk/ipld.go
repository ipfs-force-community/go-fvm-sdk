package sdk

import (
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

/// The ID of the "unit" block (or void for C programmers).
const UNIT uint32 = 0

/// Store a block. The block will only be persisted in the state-tree if the CID is "linked in" to
/// the actor's state-tree before the end of the current invocation.
func Put(mh_code uint64, mh_size uint32, codec uint64, data []byte) (cid.Cid, error) {
	id, err := sys.Create(codec, data)
	if err != nil {
		return cid.Undef, err
	}

	// I really hate this CID interface. Why can't I just have bytes?
	buf := [types.MAX_CID_LEN]byte{}
	cidLen, err := sys.Cid(id, mh_code, mh_size, buf[:])
	if err != nil {
		return cid.Undef, err
	}
	if int(cidLen) > len(buf) {
		// TODO: re-try with a larger buffer?
		panic(fmt.Sprintf("CID too big: %d > %d", cidLen, len(buf)))
	}
	_, result, err := cid.CidFromBytes(buf[:cidLen])
	if err != nil {
		return cid.Undef, err
	}
	return result, err
}

/// Get a block. It's valid to call this on:
///
/// 1. All CIDs returned by prior calls to `get_root`...
/// 2. All CIDs returned by prior calls to `put`...
/// 3. Any children of a blocks returned by prior calls to `get`...
///
/// ...during the current invocation.
func Get(cid cid.Cid) ([]byte, error) {
	// TODO: Check length of cid?
	result, err := sys.Open(cid)
	if err != nil {
		return nil, err
	}
	return GetBlock(result.Id, &result.Size)
}

/// Gets the data of the block referenced by BlockId. If the caller knows the
/// size, this function will avoid statting the block.
func GetBlock(id types.BlockId, size *uint32) ([]byte, error) {
	var size1 uint32
	if size != nil {
		size1 = *size
	} else {
		stat, err := sys.Stat(id)
		if err != nil {
			return nil, err
		}
		size1 = stat.Size
	}

	block := make([]byte, size1, size1)
	bytesRead, err := sys.Read(id, 0, block)
	if err != nil {
		return nil, err
	}
	if bytesRead != size1 {
		panic(fmt.Sprintf("read an unexpected number of bytes expect %d but got %d", size1, bytesRead))
	}
	return block[:bytesRead], nil
}

/// Writes the supplied block and returns the BlockId.
func PutBlock(codec types.Codec, data []byte) (types.BlockId, error) {
	return sys.Create(codec, data)
}

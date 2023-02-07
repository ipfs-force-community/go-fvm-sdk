package types

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/specs-actors/v7/actors/runtime"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

const (
	/// DagCBOR should be used for all IPLD-CBOR data where CIDs need to be traversable.
	DAGCBOR uint64 = 0x71
	CBOR    uint64 = 0x51
	IPLDRAW uint64 = 0x55
)

type Codec = uint64

// NoDataBlockID specify noblock in params and return
const NoDataBlockID uint32 = 0

type BlockID = uint32

func ValidateConsensusFaultType(c runtime.ConsensusFaultType) bool {
	return 0 <= c && c <= 3
}

type IpldOpen struct {
	Codec Codec
	ID    uint32
	Size  uint32
}

type IpldStat struct {
	Codec Codec
	Size  uint32
}

type ResolveAddress struct {
	Resolved int32
	Value    uint64
}

type SendResult struct {
	ExitCode    ferrors.ExitCode
	ReturnID    BlockID
	ReturnCodec uint64
	ReturnSize  uint32
}

type VerifyConsensusFault struct {
	Epoch  int64
	Target abi.ActorID
	Fault  uint32
}
type ParamsRaw struct {
	Codec Codec
	Raw   []byte
}

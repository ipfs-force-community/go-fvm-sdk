package types

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/specs-actors/v7/actors/runtime"
)

const (
	BLAKE2B256 uint64 = 0xb220
	BLAKE2BLEN uint32 = 32
)

type BlockID = uint32
type Codec = uint64

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

type Send struct {
	ExitCode uint32
	ReturnID BlockID
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

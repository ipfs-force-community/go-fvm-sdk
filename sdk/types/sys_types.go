package types

import "github.com/filecoin-project/specs-actors/v7/actors/runtime"

type BlockId = uint32
type Codec = uint64
type ActorId = uint64

func ValidateConsensusFaultType(c runtime.ConsensusFaultType) bool {
	return 0 <= c && c <= 3
}

type IpldOpen struct {
	Codec Codec
	Id    uint32
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
	ReturnID BlockId
}

type VerifyConsensusFault struct {
	Epoch  int64
	Target ActorId
	Fault  uint32
}

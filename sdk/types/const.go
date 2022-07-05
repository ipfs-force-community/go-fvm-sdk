package types

const (
	//MaxCidLen The maximum supported CID size. (SPEC_AUDIT)
	MaxCidLen = 100

	//MaxActorAddrLen The maximum actor address length (class 2 addresses).
	MaxActorAddrLen = 21

	//DAGCbor dag codec
	DAGCbor uint64 = 0x71 // TODO find something to reference.

	//NoDataBlockID specify noblock in params and return
	NoDataBlockID uint32 = 0

	//UNIT block unit
	UNIT uint32 = 0
)

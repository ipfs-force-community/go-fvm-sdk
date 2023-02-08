package types

const (
	//MaxCidLen The maximum supported CID size. (SPEC_AUDIT)
	MaxCidLen = 100

	//MaxActorAddrLen The maximum actor address length (class 2 addresses).
	MaxActorAddrLen = 21

	//UNIT block unit
	UNIT uint32 = 0

	BLAKE2B256 uint64 = 0xb220
	BLAKE2BLEN uint32 = 32
)

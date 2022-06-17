package types

const (
	//MAX_CID_LEN The maximum supported CID size. (SPEC_AUDIT)
	MAX_CID_LEN = 100

	//MAX_ACTOR_ADDR_LEN The maximum actor address length (class 2 addresses).
	MAX_ACTOR_ADDR_LEN = 21

	//DAG_CBOR dag codec
	DAG_CBOR uint64 = 0x71 // TODO find something to reference.

	//NO_DATA_BLOCK_ID specify noblock in params and return
	NO_DATA_BLOCK_ID uint32 = 0

	//UNIT block unit
	UNIT uint32 = 0
)

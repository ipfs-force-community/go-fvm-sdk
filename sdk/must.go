package sdk

import (
	"bytes"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/cbor"
)

func MustCborMarshal(obj cbor.Marshaler) []byte {
	buf := bytes.NewBuffer(nil)
	err := obj.MarshalCBOR(buf)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func MustAddressFromActorId(actorId abi.ActorID) address.Address {
	addr, _ := address.NewIDAddress(uint64(actorId))
	return addr
}

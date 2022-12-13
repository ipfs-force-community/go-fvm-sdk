package sdk

import (
	"bytes"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/cbor"
)

// MustCborMarshal marshal obj to bytes, panic if marshal fail, used it carefully
func MustCborMarshal(obj cbor.Marshaler) []byte {
	if obj == nil {
		return nil
	}
	buf := bytes.NewBuffer(nil)
	err := obj.MarshalCBOR(buf)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// MustAddressFromActorId convert actor id to address, panic if marshal fail, used it carefully
func MustAddressFromActorId(actorId abi.ActorID) address.Address {
	addr, _ := address.NewIDAddress(uint64(actorId))
	return addr
}

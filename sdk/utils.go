package sdk

import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-state-types/cbor"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func SaveState(state cbor.Marshaler) cid.Cid {
	buf := bytes.NewBuffer([]byte{})
	err := state.MarshalCBOR(buf)
	if err != nil {
		Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get root: %v", err))
	}
	stBytes := buf.Bytes()
	stCid, err := Put(0xb220, 32, types.DAGCbor, stBytes)
	if err != nil {
		Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get root: %v", err))
	}

	err = SetRoot(stCid)
	if err != nil {
		Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get root: %v", err))
	}
	return stCid
}

func Constructor(state cbor.Marshaler) error {
	caller, err := Caller()
	if err != nil {
		Abort(ferrors.USR_ILLEGAL_STATE, "unable to get caller")
	}

	if caller != 1 {
		Abort(ferrors.USR_ILLEGAL_STATE, "constructor invoked by non-init actor")
	}

	_ = SaveState(state)
	return nil
}

func LoadState(state cbor.Unmarshaler) {
	root, err := Root()
	if err != nil {
		Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get root: %v", err))
	}

	data, err := Get(root)
	if err != nil {
		Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get data: %v", err))
	}
	err = state.UnmarshalCBOR(bytes.NewReader(data))
	if err != nil {
		Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get data: %v", err))
	}
}

func LoadStateFromCid(cid cid.Cid, state cbor.Unmarshaler) { //nolint
	data, err := Get(cid)
	if err != nil {
		Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get data: %v", err))
	}
	err = state.UnmarshalCBOR(bytes.NewReader(data))
	if err != nil {
		Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get data: %v", err))
	}
}

//this code was from https://github.com/modern-go/reflect2/blob/2b33151c9bbc5231aea69b8861c540102b087070/reflect2.go#L238, and unable to use this package directly for now
type eface struct {
	_    unsafe.Pointer
	data unsafe.Pointer
}

func unpackEFace(obj interface{}) *eface {
	return (*eface)(unsafe.Pointer(&obj))
}

func IsNil(obj interface{}) bool {
	if obj == nil {
		return true
	}
	return unpackEFace(obj).data == nil
}

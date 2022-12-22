package sdk

import (
	"bytes"
	"context"
	"fmt"
	"unsafe"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-state-types/cbor"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

// SaveState save actor state
func SaveState(ctx context.Context, state cbor.Marshaler) cid.Cid {
	buf := bytes.NewBuffer([]byte{})
	err := state.MarshalCBOR(buf)
	if err != nil {
		Abort(ctx, ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get root: %v", err))
	}
	stBytes := buf.Bytes()
	stCid, err := Put(ctx, 0xb220, 32, types.DAGCbor, stBytes)
	if err != nil {
		Abort(ctx, ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get root: %v", err))
	}

	err = SetRoot(ctx, stCid)
	if err != nil {
		Abort(ctx, ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get root: %v", err))
	}
	return stCid
}

// Constructor construct a acor with initialize state
func Constructor(ctx context.Context, state cbor.Marshaler) error {
	caller, err := Caller(ctx)
	if err != nil {
		Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to get caller")
	}

	if caller != 1 {
		Abort(ctx, ferrors.USR_ILLEGAL_STATE, "constructor invoked by non-init actor")
	}

	_ = SaveState(ctx, state)
	return nil
}

// LoadState loads actors current state
func LoadState(ctx context.Context, state cbor.Unmarshaler) {
	root, err := Root(ctx)
	if err != nil {
		Abort(ctx, ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get root: %v", err))
	}

	data, err := Get(ctx, root)
	if err != nil {
		Abort(ctx, ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get data: %v", err))
	}
	err = state.UnmarshalCBOR(bytes.NewReader(data))
	if err != nil {
		Abort(ctx, ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get data: %v", err))
	}
}

// LoadStateFromCid load actor state by message cid
func LoadStateFromCid(ctx context.Context, cid cid.Cid, state cbor.Unmarshaler) { // nolint
	data, err := Get(ctx, cid)
	if err != nil {
		Abort(ctx, ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get data: %v", err))
	}
	err = state.UnmarshalCBOR(bytes.NewReader(data))
	if err != nil {
		Abort(ctx, ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get data: %v", err))
	}
}

// this code was from https://github.com/modern-go/reflect2/blob/2b33151c9bbc5231aea69b8861c540102b087070/reflect2.go#L238, and unable to use this package directly for now
type eface struct {
	_    unsafe.Pointer
	data unsafe.Pointer
}

func unpackEFace(obj interface{}) *eface {
	return (*eface)(unsafe.Pointer(&obj))
}

// IsNil check whether interface is nil
func IsNil(obj interface{}) bool { // nolint
	if obj == nil {
		return true
	}
	return unpackEFace(obj).data == nil
}

// MethodInfo used to mark actor export function.
type MethodInfo struct {
	// use alias name instead of function name
	Alias string
	// function gen tool get method params and return  by this field
	Func interface{}
	// indicate whether this method is a readonly function,  not need to send message when invoke this query state
	Readonly bool
}

package contract

import (
	"bytes"
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	cid "github.com/ipfs/go-cid"
)

type State struct {
	Count uint64
}

func (s *State) Save() cid.Cid {
	buf := bytes.NewBuffer([]byte{})
	err := s.MarshalCBOR(buf)
	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get root: %v", err))
	}
	stBytes := buf.Bytes()
	stCid, err := sdk.Put(0xb220, 32, types.DAG_CBOR, stBytes)
	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get root: %v", err))
	}

	err = sdk.SetRoot(stCid)
	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get root: %v", err))
	}
	return stCid
}

func NewState() *State {
	root, err := sdk.Root()
	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get root: %v", err))
	}

	data, err := sdk.Get(root)
	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get data: %v", err))
	}
	st := new(State)
	err = st.UnmarshalCBOR(bytes.NewReader(data))
	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get data: %v", err))
	}
	return st
}

/// The constructor populates the initial state.
///
/// Method num 1. This is part of the Filecoin calling convention.
/// InitActor#Exec will call the constructor on method_num = 1.
func Constructor() []byte {
	// This constant should be part of the SDK.
	// var  ActorID = 1;

	caller, err := sdk.Caller()
	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unbale to get caller")
	}

	if caller != 1 {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, "constructor invoked by non-init actor")
	}
	state := State{}
	_ = state.Save()
	return nil
}

/// Method num 2.
func SayHello() []byte {
	state := NewState()
	state.Count += 1
	state.Save()

	ret := fmt.Sprintf("Hello world %d!", state.Count)
	return []byte(ret)
}

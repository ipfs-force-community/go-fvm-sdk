package contract

import (
	"bytes"
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

type State struct {
	Count uint64
}

func (e *State) Export() map[int]interface{} {
	return map[int]interface{}{
		1: Constructor,
		2: e.SayHello,
	}
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
func Constructor() error {
	// This constant should be part of the SDK.
	// var  ActorID = 1;
	_ = sdk.Constructor(&State{})
	return nil
}

/// Method num 2.
func (st *State) SayHello() types.CBORBytes {
	state := NewState()
	state.Count += 1
	ret := fmt.Sprintf("Hello world %d!", state.Count)
	_ = sdk.SaveState(&State{})
	return []byte(ret)
}

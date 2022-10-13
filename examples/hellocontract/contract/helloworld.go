package contract

import (
	"context"
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
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

// The constructor populates the initial state.
// /
// Method num 1. This is part of the Filecoin calling convention.
// InitActor#Exec will call the constructor on method_num = 1.
func Constructor(ctx context.Context) error {
	// This constant should be part of the SDK.
	// var  ActorID = 1;
	s := &State{}
	_ = sdk.Constructor(ctx, s)
	return nil
}

// Method num 2.
func (st *State) SayHello(ctx context.Context) types.CBORBytes {
	st.Count += 1
	ret := fmt.Sprintf("%d", st.Count)
	_ = sdk.SaveState(ctx, &State{})
	return []byte(ret)
}

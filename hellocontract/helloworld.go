package main

import (
	"fmt"

	"hellocontract/contract"

	addr "github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

//not support non-main wasm in tinygo at present
func main() {
	_ = addr.Undef
}

/// The actor's WASM entrypoint. It takes the ID of the parameters block,
/// and returns the ID of the return value block, or NO_DATA_BLOCK_ID if no
/// return value.
///
/// Should probably have macros similar to the ones on fvm.filecoin.io snippets.
/// Put all methods inside an impl struct and annotate it with a derive macro
/// that handles state serde and dispatch.
//go:export invoke
func Invoke(_ uint32) uint32 {
	method, err := sdk.MethodNumber()
	if err != nil {
		sdk.Abort(20, "unable to get method number")
	}

	var rawBytes []byte
	switch method {
	case 1:
		rawBytes = constructor()
	case 2:
		rawBytes = say_hello()
	default:
		sdk.Abort(20, "unsupport method")
	}

	if rawBytes != nil {
		id, err := sdk.PutBlock(types.DAG_CBOR, rawBytes)
		if err != nil {
			sdk.Abort(20, fmt.Sprintf("failed to store return value: %v", err))
		}
		return id
	} else {
		return types.NO_DATA_BLOCK_ID
	}
}

/// The constructor populates the initial state.
///
/// Method num 1. This is part of the Filecoin calling convention.
/// InitActor#Exec will call the constructor on method_num = 1.
func constructor() []byte {
	// This constant should be part of the SDK.
	// var  ActorID = 1;

	caller, err := sdk.Caller()
	if err != nil {
		sdk.Abort(20, "unbale to get caller")
	}

	if caller != 1 {
		sdk.Abort(20, "constructor invoked by non-init actor")
	}
	state := contract.State{}
	_ = state.Save()
	return nil
}

/// Method num 2.
func say_hello() []byte {
	state := contract.NewState()
	state.Count += 1
	state.Save()

	ret := fmt.Sprintf("Hello world %d!", state.Count)
	return []byte(ret)
}

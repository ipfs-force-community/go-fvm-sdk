// Code generated by github.com/ipfs-force-community/go-fvm-sdk. DO NOT EDIT.
package main

import (
	"bytes"
	"fmt"

	cbor "github.com/filecoin-project/go-state-types/cbor"

	sdk "github.com/ipfs-force-community/go-fvm-sdk/sdk"

	ferrors "github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

	sdkTypes "github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	typegen "github.com/whyrusleeping/cbor-gen"

	contract "hellocontract/contract"
)

//not support non-main wasm in tinygo at present
func main() {}

/// The actor's WASM entrypoint. It takes the ID of the parameters block,
/// and returns the ID of the return value block, or NO_DATA_BLOCK_ID if no
/// return value.
///
/// Should probably have macros similar to the ones on fvm.filecoin.io snippets.
/// Put all methods inside an impl struct and annotate it with a derive macro
/// that handles state serde and dispatch.
//go:export invoke
func Invoke(blockId uint32) uint32 {
	method, err := sdk.MethodNumber()
	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to get method number")
	}

	var callResult cbor.Marshaler

	switch method {
	case 1:
		//Constuctor
		err = contract.Constructor()
		callResult = typegen.CborBool(true)

	case 2:

		//no params no error but have return value
		state := new(contract.State)
		sdk.LoadState(state)
		callResult = state.SayHello()

	default:
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unsupport method")
	}

	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("call error %s", err))
	}

	if !sdk.IsNil(callResult) {
		buf := bytes.NewBufferString("")
		err = callResult.MarshalCBOR(buf)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("marshal resp fail %s", err))
		}
		id, err := sdk.PutBlock(sdkTypes.DAGCbor, buf.Bytes())
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to store return value: %v", err))
		}
		return id
	} else {
		return sdkTypes.NoDataBlockID
	}
}

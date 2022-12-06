// Code generated by github.com/ipfs-force-community/go-fvm-sdk. DO NOT EDIT.
package main

import (
	"bytes"
	"errors"
	"fmt"

	context "context"

	cbor "github.com/filecoin-project/go-state-types/cbor"

	sdk "github.com/ipfs-force-community/go-fvm-sdk/sdk"

	ferrors "github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

	sdkTypes "github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	typegen "github.com/whyrusleeping/cbor-gen"

	contract "hellocontract/contract"
)

// not support non-main wasm in tinygo at present
func main() {}

// Invoke The actor's WASM entrypoint. It takes the ID of the parameters block,
// and returns the ID of the return value block, or NO_DATA_BLOCK_ID if no
// return value.
//
// Should probably have macros similar to the ones on fvm.filecoin.io snippets.
// Put all methods inside an impl struct and annotate it with a derive macro
// that handles state serde and dispatch.
//
//go:export invoke
func Invoke(blockId uint32) uint32 {
	ctx := context.Background()
	method, err := sdk.MethodNumber(ctx)
	if err != nil {
		sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to get method number")
	}

	var callResult cbor.Marshaler

	switch method {
	case 0x1:
		// Constuctor
		err = contract.Constructor(ctx)
		callResult = typegen.CborBool(true)

	case 0xc551429c:

		// no params no error but have return value
		state := new(contract.State)
		sdk.LoadState(ctx, state)
		callResult = state.SayHello(ctx)

	default:
		sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unsupport method")
	}

	if err != nil {
		exitCode := ferrors.USR_ILLEGAL_STATE
		errors.As(err, &exitCode)
		sdk.Abort(ctx, exitCode, fmt.Sprintf("call error %s", err))
	}

	if !sdk.IsNil(callResult) {
		buf := bytes.NewBufferString("")
		err = callResult.MarshalCBOR(buf)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("marshal resp fail %s", err))
		}
		id, err := sdk.PutBlock(ctx, sdkTypes.DAGCbor, buf.Bytes())
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to store return value: %v", err))
		}
		return id
	} else {
		return sdkTypes.NoDataBlockID
	}
}

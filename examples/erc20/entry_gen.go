// Code generated by github.com/ipfs-force-community/go-fvm-sdk. DO NOT EDIT.
package main

import (
	"bytes"
	"context"
	"fmt"

	contract "erc20/contract"

	address "github.com/filecoin-project/go-address"

	cbor "github.com/filecoin-project/go-state-types/cbor"

	sdk "github.com/ipfs-force-community/go-fvm-sdk/sdk"

	ferrors "github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

	sdkTypes "github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	typegen "github.com/whyrusleeping/cbor-gen"
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
	var raw *sdkTypes.ParamsRaw
	switch method {
	case 1:
		// Constuctor
		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.ConstructorReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}
		err = contract.Constructor(ctx, &req)
		callResult = typegen.CborBool(true)

	case 2:

		// no params no error but have return value
		state := new(contract.Erc20Token)
		sdk.LoadState(ctx, state)
		callResult = state.GetName()

	case 3:

		// no params no error but have return value
		state := new(contract.Erc20Token)
		sdk.LoadState(ctx, state)
		callResult = state.GetSymbol()

	case 4:

		// no params no error but have return value
		state := new(contract.Erc20Token)
		sdk.LoadState(ctx, state)
		callResult = state.GetDecimal()

	case 5:

		// no params no error but have return value
		state := new(contract.Erc20Token)
		sdk.LoadState(ctx, state)
		callResult = state.GetTotalSupply()

	case 6:

		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req address.Address
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		// have params/return/error
		state := new(contract.Erc20Token)
		sdk.LoadState(ctx, state)
		callResult, err = state.GetBalanceOf(ctx, &req)

	case 7:

		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.TransferReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		// have params/error but no return val
		state := new(contract.Erc20Token)
		sdk.LoadState(ctx, state)
		if err = state.Transfer(ctx, &req); err == nil {
			callResult = typegen.CborBool(true)
		}

	case 8:

		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.TransferFromReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		// have params/error but no return val
		state := new(contract.Erc20Token)
		sdk.LoadState(ctx, state)
		if err = state.TransferFrom(ctx, &req); err == nil {
			callResult = typegen.CborBool(true)
		}

	case 9:

		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.ApprovalReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		// have params/error but no return val
		state := new(contract.Erc20Token)
		sdk.LoadState(ctx, state)
		if err = state.Approval(ctx, &req); err == nil {
			callResult = typegen.CborBool(true)
		}

	case 10:

		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.AllowanceReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		// have params/return/error
		state := new(contract.Erc20Token)
		sdk.LoadState(ctx, state)
		callResult, err = state.Allowance(ctx, &req)

	default:
		sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unsupport method")
	}

	if err != nil {
		sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("call error %s", err))
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
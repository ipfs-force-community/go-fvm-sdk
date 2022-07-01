// Code generated by github.com/ipfs-force-community/go-fvm-sdk. DO NOT EDIT.
package main

import (
	"bytes"
	"fmt"

	contract "erc20/contract"

	address "github.com/filecoin-project/go-address"

	cbor "github.com/filecoin-project/go-state-types/cbor"

	sdk "github.com/ipfs-force-community/go-fvm-sdk/sdk"

	ferrors "github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	sdkTypes "github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	typegen "github.com/whyrusleeping/cbor-gen"
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
	var raw *sdkTypes.ParamsRaw
	switch method {
	case 1:
		//Constuctor
		raw, err = sdk.ParamsRaw(blockId)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.ConstructorReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}
		err = contract.Constructor(&req)
		callResult = typegen.CborBool(true)

	case 2:

		//no params no error but have return value
		state := new(contract.Erc20Token)
		sdk.LoadState(state)
		callResult = state.GetName()

	case 3:

		//no params no error but have return value
		state := new(contract.Erc20Token)
		sdk.LoadState(state)
		callResult = state.GetSymbol()

	case 4:

		//no params no error but have return value
		state := new(contract.Erc20Token)
		sdk.LoadState(state)
		callResult = state.GetDecimal()

	case 5:

		//no params no error but have return value
		state := new(contract.Erc20Token)
		sdk.LoadState(state)
		callResult = state.GetTotalSupply()

	case 6:

		raw, err = sdk.ParamsRaw(blockId)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req address.Address
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		//have params/return/error
		state := new(contract.Erc20Token)
		sdk.LoadState(state)
		callResult, err = state.GetBalanceOf(&req)

	case 7:

		raw, err = sdk.ParamsRaw(blockId)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.TransferReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		//have params/error but no return val
		state := new(contract.Erc20Token)
		sdk.LoadState(state)
		if err = state.Transfer(&req); err == nil {
			callResult = typegen.CborBool(true)
		}

	case 8:

		raw, err = sdk.ParamsRaw(blockId)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.TransferFromReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		//have params/error but no return val
		state := new(contract.Erc20Token)
		sdk.LoadState(state)
		if err = state.TransferFrom(&req); err == nil {
			callResult = typegen.CborBool(true)
		}

	case 9:

		raw, err = sdk.ParamsRaw(blockId)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.ApprovalReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		//have params/error but no return val
		state := new(contract.Erc20Token)
		sdk.LoadState(state)
		if err = state.Approval(&req); err == nil {
			callResult = typegen.CborBool(true)
		}

	case 10:

		raw, err = sdk.ParamsRaw(blockId)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.AllowanceReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		//have params/return/error
		state := new(contract.Erc20Token)
		sdk.LoadState(state)
		callResult, err = state.Allowance(&req)

	case 11:

		raw, err = sdk.ParamsRaw(blockId)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.FakeSetBalance
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		//have params/error but no return val
		state := new(contract.Erc20Token)
		sdk.LoadState(state)
		if err = state.FakeSetBalance(&req); err == nil {
			callResult = typegen.CborBool(true)
		}

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
		id, err := sdk.PutBlock(types.DAGCbor, buf.Bytes())
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to store return value: %v", err))
		}
		return id
	} else {
		return types.NoDataBlockID
	}
}

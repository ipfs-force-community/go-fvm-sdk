// Code generated by github.com/ipfs-force-community/go-fvm-sdk. DO NOT EDIT.
package main

import (
	"bytes"
	"errors"
	"fmt"

	context "context"

	contract "frc46token/contract"

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
	case 0x1:
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

	case 0x2ea015c:

		// no params no error but have return value
		state := new(contract.Frc46Token)
		sdk.LoadState(ctx, state)
		callResult = state.GetName(ctx)

	case 0x6d0c41e0:

		// no params no error but have return value
		state := new(contract.Frc46Token)
		sdk.LoadState(ctx, state)
		callResult = state.GetSymbol(ctx)

	case 0x9ed113b6:

		// no params no error but have return value
		state := new(contract.Frc46Token)
		sdk.LoadState(ctx, state)
		callResult = state.GetGranularity(ctx)

	case 0xe435da43:

		// no params no error but have return value
		state := new(contract.Frc46Token)
		sdk.LoadState(ctx, state)
		callResult = state.GetTotalSupply(ctx)

	case 0x6f84ab2:

		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.MintParams
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		// have params/return/error
		state := new(contract.Frc46Token)
		sdk.LoadState(ctx, state)
		callResult, err = state.Mint(ctx, &req)

	case 0x8710e1ac:

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
		state := new(contract.Frc46Token)
		sdk.LoadState(ctx, state)
		callResult, err = state.BalanceOf(ctx, &req)

	case 0xfaa45236:

		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.GetAllowanceParams
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		// have params/return/error
		state := new(contract.Frc46Token)
		sdk.LoadState(ctx, state)
		callResult, err = state.Allowance(ctx, &req)

	case 0x4cbf732:

		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.TransferParams
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		// have params/return/error
		state := new(contract.Frc46Token)
		sdk.LoadState(ctx, state)
		callResult, err = state.Transfer(ctx, &req)

	case 0xd7d4deed:

		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.TransferFromParams
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		// have params/return/error
		state := new(contract.Frc46Token)
		sdk.LoadState(ctx, state)
		callResult, err = state.TransferFrom(ctx, &req)

	case 0x69ecb918:

		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.IncreaseAllowanceParams
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		// have params/return/error
		state := new(contract.Frc46Token)
		sdk.LoadState(ctx, state)
		callResult, err = state.IncreaseAllowance(ctx, &req)

	case 0x5b286f21:

		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.DecreaseAllowanceParams
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		// have params/return/error
		state := new(contract.Frc46Token)
		sdk.LoadState(ctx, state)
		callResult, err = state.DecreaseAllowance(ctx, &req)

	case 0xa4d840b1:

		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.RevokeAllowanceParams
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		// have params/return/error
		state := new(contract.Frc46Token)
		sdk.LoadState(ctx, state)
		callResult, err = state.RevokeAllowance(ctx, &req)

	case 0x5584159a:

		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.BurnParams
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		// have params/return/error
		state := new(contract.Frc46Token)
		sdk.LoadState(ctx, state)
		callResult, err = state.Burn(ctx, &req)

	case 0xb19a37a2:

		raw, err = sdk.ParamsRaw(ctx, blockId)
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.BurnFromParams
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
		}

		// have params/return/error
		state := new(contract.Frc46Token)
		sdk.LoadState(ctx, state)
		callResult, err = state.BurnFrom(ctx, &req)

	default:
		sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unsupport method")
	}

	if err != nil {
		exitCode := ferrors.USR_ILLEGAL_STATE
		errors.As(err, exitCode)
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

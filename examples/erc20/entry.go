package main

import (
	"bytes"
	"erc20/contract"
	"fmt"
	"strconv"

	"github.com/filecoin-project/go-address"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

//not support non-main wasm in tinygo at present
func main() {}

var Success = []byte("success")

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

	var rawBytes []byte
	switch method {
	case 1:
		raw, err := sdk.ParamsRaw(blockId)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.ConstructorReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		contract.Constructor(req.Name, req.Symbol, req.Decimals, req.TotalSupply)
		rawBytes = Success
	case 2: //GetName
		rawBytes = []byte(contract.LoadToken().GetName())
	case 3: //GetSymbol
		rawBytes = []byte(contract.LoadToken().GetSymbol())
	case 4: //GetDecimal
		decimal := contract.LoadToken().GetDecimal()
		rawBytes = []byte(strconv.Itoa(int(decimal)))
	case 5: //GetTotalSupply
		supply := contract.LoadToken().GetTotalSupply()
		rawBytes = []byte(supply.String())
	case 6: //GetBalanceOf
		raw, err := sdk.ParamsRaw(blockId)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		addr, err := address.NewFromString(string(raw.Raw))
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to parser address")
		}
		balance, err := contract.LoadToken().GetBalanceOf(addr)
		rawBytes = []byte(balance.String())
	case 7: //Transfer
		raw, err := sdk.ParamsRaw(blockId)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.TransferReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		err = contract.LoadToken().Transfer(req.ReceiverAddr, req.TransferAmount)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to transfer")
		}
		rawBytes = Success
	case 8: //Allowance
		raw, err := sdk.ParamsRaw(blockId)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.AllowanceReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		allowBalance, err := contract.LoadToken().Allowance(req.OwnerAddr, req.SpenderAddr)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to allowance")
		}
		rawBytes = []byte(allowBalance.String())
	case 9: //TransferFrom
		raw, err := sdk.ParamsRaw(blockId)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.TransferFromReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		err = contract.LoadToken().TransferFrom(req.OwnerAddr, req.SpenderAddr, req.TransferAmount)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to transfer from")
		}
		rawBytes = Success
	case 10: //Approval
		raw, err := sdk.ParamsRaw(blockId)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		var req contract.ApprovalReq
		err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
		}
		err = contract.LoadToken().Approval(req.SpenderAddr, req.NewAllowance)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to approval")
		}
		rawBytes = Success
	default:
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unsupport method")
	}

	if rawBytes != nil {
		id, err := sdk.PutBlock(types.DAG_CBOR, rawBytes)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to store return value: %v", err))
		}
		return id
	} else {
		return types.NO_DATA_BLOCK_ID
	}
}

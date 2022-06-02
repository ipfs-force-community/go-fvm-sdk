package sdk

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

var InvocationCtx *types.InvocationContext

func Caller() (abi.ActorID, error) {

	if InvocationCtx == nil {
		var err error
		InvocationCtx, err = sys.VmContext()
		if err != nil {
			return 0, err
		}
	}
	return InvocationCtx.Caller, nil
}

func Receiver() (abi.ActorID, error) {
	if InvocationCtx == nil {
		var err error
		InvocationCtx, err = sys.VmContext()
		if err != nil {
			return 0, err
		}
	}
	return InvocationCtx.Receiver, nil
}

func MethodNumber() (abi.MethodNum, error) {
	if InvocationCtx == nil {
		var err error
		InvocationCtx, err = sys.VmContext()
		if err != nil {
			return 0, err
		}
	}
	return InvocationCtx.MethodNumber, nil
}

func ValueReceived() (*types.TokenAmount, error) {
	if InvocationCtx == nil {
		var err error
		InvocationCtx, err = sys.VmContext()
		if err != nil {
			return nil, err
		}
	}

	return &types.TokenAmount{ //avoud change
		Lo: InvocationCtx.ValueReceived.Lo,
		Hi: InvocationCtx.ValueReceived.Hi,
	}, nil
}

func CurrEpoch() (abi.ChainEpoch, error) {
	if InvocationCtx == nil {
		var err error
		InvocationCtx, err = sys.VmContext()
		if err != nil {
			return 0, err
		}
	}
	return InvocationCtx.NetworkCurrEpoch, nil
}

func Version() (network.Version, error) {
	if InvocationCtx == nil {
		var err error
		InvocationCtx, err = sys.VmContext()
		if err != nil {
			return 0, err
		}
	}
	return network.Version(InvocationCtx.NetworkVersion), nil
}

package sdk

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

// Caller get caller, from address of message
// todo cache invocation in context
func Caller(ctx context.Context) (abi.ActorID, error) {
	msgCtx, err := sys.VMMessageContext(ctx)
	if err != nil {
		return 0, err
	}
	return msgCtx.Caller, nil
}

// Origin get caller, from address of message
func Origin(ctx context.Context) (abi.ActorID, error) {
	msgCtx, err := sys.VMMessageContext(ctx)
	if err != nil {
		return 0, err
	}
	return msgCtx.Origin, nil
}

// Receiver get recevier, to address of message
func Receiver(ctx context.Context) (abi.ActorID, error) {
	msgCtx, err := sys.VMMessageContext(ctx)
	if err != nil {
		return 0, err
	}
	return msgCtx.Receiver, nil
}

// MethodNumber method number
func MethodNumber(ctx context.Context) (abi.MethodNum, error) {
	msgCtx, err := sys.VMMessageContext(ctx)
	if err != nil {
		return 0, err
	}
	return msgCtx.MethodNumber, nil
}

// ValueReceived the amount was transferred in message
func ValueReceived(ctx context.Context) (abi.TokenAmount, error) {
	msgCtx, err := sys.VMMessageContext(ctx)
	if err != nil {
		return abi.TokenAmount{}, err
	}
	return msgCtx.ValueReceived, nil
}

// ValueReceived the amount was transferred in message
func MsgGasPremium(ctx context.Context) (abi.TokenAmount, error) {
	msgCtx, err := sys.VMMessageContext(ctx)
	if err != nil {
		return abi.TokenAmount{}, err
	}
	return msgCtx.GasPremium, nil
}

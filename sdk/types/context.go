package types

import (
	"github.com/filecoin-project/go-state-types/abi"
)

type SendFlags = uint64

const ReadonlyFlag = 0b00000001

type NetworkContext struct {
	/// The current epoch.
	Epoch abi.ChainEpoch
	/// The current time (seconds since the unix epoch).
	Timestamp uint64
	/// The current base-fee.
	BaseFee abi.TokenAmount
	/// The network version.
	NetworkVersion uint32
}

type MessageContext struct {
	/// The current call's origin actor ID.
	Origin abi.ActorID
	/// The caller's actor ID.
	Caller abi.ActorID
	/// The receiver's actor ID (i.e. ourselves).
	Receiver abi.ActorID
	/// The method number from the message.
	MethodNumber abi.MethodNum
	/// The value that was received.
	ValueReceived abi.TokenAmount
	/// The current gas premium
	GasPremium abi.TokenAmount
	/// Flags pertaining to the currently executing actor's invocation context.
	Flags SendFlags
}

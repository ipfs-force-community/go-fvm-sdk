package types

import (
	"github.com/filecoin-project/go-state-types/abi"
)

type InvocationContext struct {
	/// The value that was received.
	ValueReceived TokenAmount
	/// The caller's actor ID.
	Caller abi.ActorID
	/// The receiver's actor ID (i.e. ourselves).
	Receiver abi.ActorID

	/// The method number from the message.
	MethodNumber abi.MethodNum
	/// The current epoch.
	NetworkCurrEpoch abi.ChainEpoch
	/// The network version.
	NetworkVersion uint32
}

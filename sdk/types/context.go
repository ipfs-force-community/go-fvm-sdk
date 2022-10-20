package types

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/network"
)

type InvocationContext struct {
	/// The value that was received.
	ValueReceived abi.TokenAmount
	/// The caller's actor ID.
	Caller abi.ActorID
	/// The receiver's actor ID (i.e. ourselves).
	Receiver abi.ActorID

	/// The method number from the message.
	MethodNumber abi.MethodNum
	/// The current epoch.
	NetworkCurrEpoch abi.ChainEpoch
	/// The network version.
	NetworkVersion network.Version
}

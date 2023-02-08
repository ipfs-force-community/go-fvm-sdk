// Package types go-fvm-sdk types
package types

type ActorType int32

const (
	System           ActorType = 1
	Init             ActorType = 2
	Cron             ActorType = 3
	Account          ActorType = 4
	Power            ActorType = 5
	Miner            ActorType = 6
	Market           ActorType = 7
	PaymentChannel   ActorType = 8
	Multisig         ActorType = 9
	Reward           ActorType = 10
	VerifiedRegistry ActorType = 11

	PlaceHolder ActorType = 12
	Evm         ActorType = 13
	Eam         ActorType = 14
	EthAccount  ActorType = 15
)

func (t ActorType) IsSingletonActor() bool {
	switch t {
	case Init:
		fallthrough
	case Reward:
		fallthrough
	case Cron:
		fallthrough
	case Power:
		fallthrough
	case Market:
		fallthrough
	case VerifiedRegistry:
		return true
	case Eam:
		return true
	default:
		return false
	}
}

func (t ActorType) IsAccount() bool {
	return t == Account
}

func (t ActorType) IsPrincipal() bool {
	return t == Account || t == Multisig
}

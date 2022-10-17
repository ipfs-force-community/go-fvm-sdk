package types

import (
	"github.com/filecoin-project/go-state-types/abi"
)

// ActorKey adapts an actor id as a mapping key.
type ActorKey abi.ActorID

// Key get actor id string as key
func (k ActorKey) Key() string {
	return abi.ActorID(k).String()
}

// StringKey Adapts an string as a mapping key.
type StringKey string

// Key return string as key
func (k StringKey) Key() string {
	return string(k)
}

type emptyKeyType struct{}

var SimulatedEnvkey emptyKeyType

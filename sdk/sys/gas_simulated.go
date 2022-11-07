//go:build simulate
// +build simulate

package sys

import (
	"context"
)

// Charge charge gas for the operation identified by name.
func Charge(ctx context.Context, name string, compute uint64) error {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.GasCharge(name, compute)
	}
	panic(ErrorEnvValid)
}

// Returns the amount of gas remaining.
func Available(_ context.Context) (uint64, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.GasAvailable()
	}
	panic(ErrorEnvValid)
}

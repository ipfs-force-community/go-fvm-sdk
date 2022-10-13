//go:build simulated
// +build simulated

package sys

import (
	"context"
)

// Charge charge gas for the operation identified by name.
func Charge(ctx context.Context, name string, compute uint64) error {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.Charge(name, compute)
	}
	return nil
}

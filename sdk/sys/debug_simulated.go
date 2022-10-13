//go:build simulated
// +build simulated

package sys

import (
	"context"
)

func Enabled(ctx context.Context) (bool, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.Enabled()
	}
	return false, nil
}

func Log(ctx context.Context, msg string) error {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.Log(msg)
	}
	return nil
}

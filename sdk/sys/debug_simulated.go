//go:build simulate
// +build simulate

package sys

import (
	"context"
)

func Enabled(ctx context.Context) (bool, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.Enabled()
	}
	panic(ErrorEnvValid)
}

func Log(ctx context.Context, msg string) error {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.Log(msg)
	}
	panic(ErrorEnvValid)
}
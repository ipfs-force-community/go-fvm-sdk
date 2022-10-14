//go:build simulate
// +build simulate

package sys

import (
	"context"
)

func Abort(ctx context.Context, code uint32, msg string) {
	if env, ok := tryGetSimulator(ctx); ok {
		env.Abort(code, msg)
		return
	}
	panic(ErrorEnvValid)
}

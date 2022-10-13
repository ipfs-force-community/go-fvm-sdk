//go:build simulated
// +build simulated

package sys

import (
	"context"
)

func Abort(ctx context.Context, code uint32, msg string) {
	if env, ok := isSimulatedEnv(ctx); ok {
		env.Abort(code, msg)
		return
	}
	
}

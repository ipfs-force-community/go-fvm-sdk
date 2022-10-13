package sdk

import (
	"context"
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

// Logger is a debug-only logger that uses the FVM syscalls.
type Logger interface {
	Enabled(ctx context.Context) bool
	Log(ctx context.Context, args ...interface{}) error
	Logf(ctx context.Context, format string, a ...interface{}) error
}

var _ Logger = (*logger)(nil)

// NewLogger create a logging if debugging is enabled.
func NewLogger(ctx context.Context) (Logger, error) {
	return &logger{}, nil
}

type logger struct {
	enable *bool
}

// inline
func (l *logger) Enabled(ctx context.Context) bool {
	if l.enable != nil {
		return *l.enable
	}
	logEnable, _ := sys.Enabled(ctx)
	l.enable = &logEnable
	return logEnable
}

func (l *logger) Log(ctx context.Context, a ...interface{}) error {
	if l.Enabled(ctx) {
		return sys.Log(ctx, fmt.Sprint(a...))
	}
	return nil
}

func (l *logger) Logf(ctx context.Context, format string, a ...interface{}) error {
	if l.Enabled(ctx) {
		return sys.Log(ctx, fmt.Sprintf(format, a...))
	}
	return nil
}

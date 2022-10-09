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
	debugEnabled, err := sys.Enabled(ctx)
	if err != nil {
		return nil, err
	}
	return &logger{
		enable: debugEnabled,
	}, nil
}

type logger struct {
	enable bool
}

func (l *logger) Enabled(ctx context.Context) bool {
	return l.enable
}

func (l *logger) Log(ctx context.Context, a ...interface{}) error {
	if l.enable {
		return sys.Log(ctx, fmt.Sprint(a...))
	}
	return nil
}

func (l *logger) Logf(ctx context.Context, format string, a ...interface{}) error {
	if l.enable {
		return sys.Log(ctx, fmt.Sprintf(format, a...))
	}
	return nil
}

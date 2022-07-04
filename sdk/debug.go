package sdk

import (
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

// Logger is a debug-only logger that uses the FVM syscalls.
type Logger interface {
	Enabled() bool
	Log(string) error
	Logf(format string, a ...interface{}) error
}

var _ Logger = (*logger)(nil)

// NewLogger create a logging if debugging is enabled.
func NewLogger() (Logger, error) {
	debugEnabled, err := sys.Enabled()
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

func (l *logger) Enabled() bool {
	return l.enable
}

func (l *logger) Log(msg string) error {
	if l.enable {
		return sys.Log(msg)
	}
	return nil
}

func (l *logger) Logf(format string, a ...interface{}) error {
	if l.enable {
		return sys.Log(fmt.Sprintf(format, a...))
	}
	return nil
}

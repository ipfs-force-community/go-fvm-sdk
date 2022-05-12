package sdk

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

var DebugEnabled bool

type Logger interface {
	Enabled() bool
	Log(string) error
}

var _ Logger = (*logger)(nil)

func init() {
	DebugEnabled, _ = sys.Enabled()
}

func NewLogger() (Logger, error) {
	return &logger{
		enable: DebugEnabled,
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

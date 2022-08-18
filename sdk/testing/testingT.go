package testing

import (
	"bytes"
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

type TestingT struct {
	errBuf *bytes.Buffer
	logger sdk.Logger
	failed bool
}

func NewTestingT() *TestingT {
	logger, _ := sdk.NewLogger()
	return &TestingT{
		errBuf: bytes.NewBufferString(""),
		logger: logger,
	}
}

func (t *TestingT) Errorf(format string, args ...interface{}) {
	t.errBuf.WriteString(fmt.Sprintf(format, args...))
	t.failed = true
}

func (t *TestingT) Error(args ...interface{}) {
	t.errBuf.WriteString(fmt.Sprint(args...))
	t.failed = true
}

func (t *TestingT) Infof(format string, args ...interface{}) {
	t.logger.Logf(format, args...)
}

func (t *TestingT) Info(args ...interface{}) {
	t.logger.Log(args...)
}

func (t *TestingT) FailNow() {
	t.failed = true
	sdk.Abort(ferrors.SYS_ASSERTION_FAILED, t.errBuf.String())
}

func (t *TestingT) CheckResult() {

	if t.failed {
		sdk.Abort(ferrors.SYS_ASSERTION_FAILED, fmt.Sprintf("assert fail:\n"+t.errBuf.String()))
	}
}

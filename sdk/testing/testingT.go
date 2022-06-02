package testing

import (
	"bytes"
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

type TestingT struct {
	buf    *bytes.Buffer
	failed bool
}

func NewTestingT() *TestingT {
	return &TestingT{
		buf: bytes.NewBufferString(""),
	}
}

func (t *TestingT) Errorf(format string, args ...interface{}) {
	t.buf.WriteString(fmt.Sprintf(format, args...))
	t.failed = true
}
func (t *TestingT) Error(args ...interface{}) {
	t.buf.WriteString(fmt.Sprint(args...))
	t.failed = true
}

func (t *TestingT) FailNow() {
	t.failed = true
	sdk.Abort(ferrors.SYS_ASSERTION_FAILED, t.buf.String())
}

func (t *TestingT) CheckResult() {
	if t.failed {
		sdk.Abort(ferrors.SYS_ASSERTION_FAILED, fmt.Sprintf("assert fail:\n"+t.buf.String()))
	}
}

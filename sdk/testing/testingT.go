package testing

import (
	"bytes"
	"fmt"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

type TestingT struct {
	buf *bytes.Buffer
}

func NewTestingT() *TestingT {
	return &TestingT{
		buf: bytes.NewBufferString(""),
	}
}

func(t *TestingT) Errorf(format string, args ...interface{}) {
	t.buf.WriteString(fmt.Sprintf(format, args))
}
func(t *TestingT) Error( args ...interface{}) {
	t.buf.WriteString(fmt.Sprint(args))
}

func(t *TestingT)  FailNow() {
	sdk.Abort(ferrors.USR_ILLEGAL_STATE, t.buf.String())
}

func(t *TestingT)  CheckResult() {
	errStr := t.buf.String()
	if len(errStr) != 0 {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE,fmt.Sprintf("assert fail:\n"+errStr))
	}
}
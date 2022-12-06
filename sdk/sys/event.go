//go:build !simulate
// +build !simulate

package sys

import (
	"bytes"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func EmitEvent(evt types.ActorEvent) error {
	buf := bytes.NewBuffer(nil)
	err := evt.MarshalCBOR(buf)
	if err != nil {
		return err
	}

	bufPtr, bufLen := GetSlicePointerAndLen(buf)
	code := emitEvent(bufPtr, bufLen)
	if code != 0 {
		return ferrors.NewSysCallError(ferrors.ErrorNumber(code), "failed to get debug-enabled")
	}

	return nil
}

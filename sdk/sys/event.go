//go:build !simulate
// +build !simulate

package sys

import (
	"bytes"
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/internal"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

// EmitEvent an actor event. It takes an DAG-CBOR encoded ActorEvent that has been
// written to Wasm memory, as an offset and length tuple.
//
// The FVM validates the structural, syntatic, and semantic correctness of the
// supplied event, and errors with `IllegalArgument` if the payload was invalid.
//
// Calling this syscall may immediately halt execution with an out of gas error,
// if such condition arises.
func EmitEvent(_ context.Context, evt types.ActorEvent) error {
	buf := bytes.NewBuffer(nil)
	err := internal.WriteCborArray(buf, evt.Entries)
	if err != nil {
		return err
	}
	bufPtr, bufLen := GetSlicePointerAndLen(buf.Bytes())
	code := emitEvent(bufPtr, bufLen)
	if code != 0 {
		return ferrors.NewSysCallError(ferrors.ErrorNumber(code), "failed to call emit event")
	}
	return nil
}

package types

import (
	"fmt"
	"io"

	typegen "github.com/whyrusleeping/cbor-gen"
)

type CborString string

func (cb CborString) MarshalCBOR(w io.Writer) error {
	if len(cb) > typegen.MaxLength {
		return fmt.Errorf("cborstring exceed max length")
	}

	if err := typegen.WriteMajorTypeHeader(w, typegen.MajTextString, uint64(len(cb))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(cb)); err != nil {
		return err
	}
	return nil
}

func (cb *CborString) UnmarshalCBOR(r io.Reader) error {
	str, err := typegen.ReadString(r)
	if err != nil {
		return err
	}
	*cb = CborString(str)
	return nil
}

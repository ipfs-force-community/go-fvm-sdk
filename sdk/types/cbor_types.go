package types

import (
	"bytes"
	"fmt"
	"io"

	cbg "github.com/whyrusleeping/cbor-gen"
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

type CborInt = typegen.CborInt
type CborUint uint64

func (cb CborUint) MarshalCBOR(w io.Writer) error {
	cw := typegen.NewCborWriter(w)
	return cw.WriteMajorTypeHeader(typegen.MajUnsignedInt, uint64(cb))
}

func (cb *CborUint) UnmarshalCBOR(r io.Reader) error {
	cr := typegen.NewCborReader(r)
	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for uint64 field")
	}
	*cb = CborUint(extra)
	return nil
}

// CBORBytes Wraps already-serialized bytes as CBOR-marshalable.
type CBORBytes []byte

func (b CBORBytes) MarshalCBOR(w io.Writer) error {
	_, err := w.Write(b)
	return err
}

func (b *CBORBytes) UnmarshalCBOR(r io.Reader) error {
	var c bytes.Buffer
	_, err := c.ReadFrom(r)
	*b = c.Bytes()
	return err
}

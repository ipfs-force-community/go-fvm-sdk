package internal

import (
	"errors"
	"fmt"
	"io"

	cbg "github.com/whyrusleeping/cbor-gen"
)

// WriteCborArray marshal cbor array to bytes
func WriteCborArray[T cbg.CBORMarshaler](w io.Writer, arr []T) error {
	if len(arr) > cbg.MaxLength {
		return errors.New("slice was too long")
	}
	cw := cbg.NewCborWriter(w)

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(arr))); err != nil {
		return err
	}
	for _, v := range arr {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}
	}
	return nil
}

// ReadCborArray reader cbor array from reader
func ReadCborArray[T any, PT interface {
	cbg.CBORUnmarshaler
	*T
}](w io.Reader) ([]*T, error) {
	cr := cbg.NewCborReader(w)
	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if extra > cbg.MaxLength {
		return nil, fmt.Errorf("t.Entries: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return nil, fmt.Errorf("expected cbor array")
	}

	vals := make([]*T, extra)
	for i := 0; i < int(extra); i++ {
		v := PT(new(T))
		if err := v.UnmarshalCBOR(cr); err != nil {
			return nil, err
		}
		vals[i] = (*T)(v)
	}

	return vals, nil
}

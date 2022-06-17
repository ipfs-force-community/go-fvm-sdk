// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package contract

import (
	"fmt"
	"io"
	"math"
	"sort"

	big "github.com/filecoin-project/go-state-types/big"
	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = math.E
var _ = sort.Sort

var lengthBufErc20Token = []byte{134}

func (t *Erc20Token) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufErc20Token); err != nil {
		return err
	}

	// t.Name (string) (string)
	if len(t.Name) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Name was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Name))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Name)); err != nil {
		return err
	}

	// t.Symbol (string) (string)
	if len(t.Symbol) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Symbol was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Symbol))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Symbol)); err != nil {
		return err
	}

	// t.Decimals (uint8) (uint8)
	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Decimals)); err != nil {
		return err
	}

	// t.TotalSupply (big.Int) (struct)
	if err := t.TotalSupply.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Balances (map[string]*big.Int) (map)
	{
		if len(t.Balances) > 4096 {
			return xerrors.Errorf("cannot marshal t.Balances map too large")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajMap, uint64(len(t.Balances))); err != nil {
			return err
		}

		keys := make([]string, 0, len(t.Balances))
		for k := range t.Balances {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := t.Balances[k]

			if len(k) > cbg.MaxLength {
				return xerrors.Errorf("Value in field k was too long")
			}

			if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(k))); err != nil {
				return err
			}
			if _, err := io.WriteString(w, string(k)); err != nil {
				return err
			}

			if err := v.MarshalCBOR(cw); err != nil {
				return err
			}

		}
	}

	// t.Allowed (map[string]*big.Int) (map)
	{
		if len(t.Allowed) > 4096 {
			return xerrors.Errorf("cannot marshal t.Allowed map too large")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajMap, uint64(len(t.Allowed))); err != nil {
			return err
		}

		keys := make([]string, 0, len(t.Allowed))
		for k := range t.Allowed {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := t.Allowed[k]

			if len(k) > cbg.MaxLength {
				return xerrors.Errorf("Value in field k was too long")
			}

			if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(k))); err != nil {
				return err
			}
			if _, err := io.WriteString(w, string(k)); err != nil {
				return err
			}

			if err := v.MarshalCBOR(cw); err != nil {
				return err
			}

		}
	}
	return nil
}

func (t *Erc20Token) UnmarshalCBOR(r io.Reader) (err error) {
	*t = Erc20Token{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 6 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Name (string) (string)

	{
		sval, err := cbg.ReadString(cr)
		if err != nil {
			return err
		}

		t.Name = string(sval)
	}
	// t.Symbol (string) (string)

	{
		sval, err := cbg.ReadString(cr)
		if err != nil {
			return err
		}

		t.Symbol = string(sval)
	}
	// t.Decimals (uint8) (uint8)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for uint8 field")
	}
	if extra > math.MaxUint8 {
		return fmt.Errorf("integer in input was too large for uint8 field")
	}
	t.Decimals = uint8(extra)
	// t.TotalSupply (big.Int) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			t.TotalSupply = new(big.Int)
			if err := t.TotalSupply.UnmarshalCBOR(cr); err != nil {
				return xerrors.Errorf("unmarshaling t.TotalSupply pointer: %w", err)
			}
		}

	}
	// t.Balances (map[string]*big.Int) (map)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("expected a map (major type 5)")
	}
	if extra > 4096 {
		return fmt.Errorf("t.Balances: map too large")
	}

	t.Balances = make(map[string]*big.Int, extra)

	for i, l := 0, int(extra); i < l; i++ {

		var k string

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			k = string(sval)
		}

		var v *big.Int

		{

			b, err := cr.ReadByte()
			if err != nil {
				return err
			}
			if b != cbg.CborNull[0] {
				if err := cr.UnreadByte(); err != nil {
					return err
				}
				v = new(big.Int)
				if err := v.UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling v pointer: %w", err)
				}
			}

		}

		t.Balances[k] = v

	}
	// t.Allowed (map[string]*big.Int) (map)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("expected a map (major type 5)")
	}
	if extra > 4096 {
		return fmt.Errorf("t.Allowed: map too large")
	}

	t.Allowed = make(map[string]*big.Int, extra)

	for i, l := 0, int(extra); i < l; i++ {

		var k string

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			k = string(sval)
		}

		var v *big.Int

		{

			b, err := cr.ReadByte()
			if err != nil {
				return err
			}
			if b != cbg.CborNull[0] {
				if err := cr.UnreadByte(); err != nil {
					return err
				}
				v = new(big.Int)
				if err := v.UnmarshalCBOR(cr); err != nil {
					return xerrors.Errorf("unmarshaling v pointer: %w", err)
				}
			}

		}

		t.Allowed[k] = v

	}
	return nil
}

var lengthBufConstructorReq = []byte{132}

func (t *ConstructorReq) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufConstructorReq); err != nil {
		return err
	}

	// t.Name (string) (string)
	if len(t.Name) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Name was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Name))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Name)); err != nil {
		return err
	}

	// t.Symbol (string) (string)
	if len(t.Symbol) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Symbol was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Symbol))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Symbol)); err != nil {
		return err
	}

	// t.Decimals (uint8) (uint8)
	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Decimals)); err != nil {
		return err
	}

	// t.TotalSupply (big.Int) (struct)
	if err := t.TotalSupply.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *ConstructorReq) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ConstructorReq{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 4 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Name (string) (string)

	{
		sval, err := cbg.ReadString(cr)
		if err != nil {
			return err
		}

		t.Name = string(sval)
	}
	// t.Symbol (string) (string)

	{
		sval, err := cbg.ReadString(cr)
		if err != nil {
			return err
		}

		t.Symbol = string(sval)
	}
	// t.Decimals (uint8) (uint8)

	maj, extra, err = cr.ReadHeader()
	if err != nil {
		return err
	}
	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for uint8 field")
	}
	if extra > math.MaxUint8 {
		return fmt.Errorf("integer in input was too large for uint8 field")
	}
	t.Decimals = uint8(extra)
	// t.TotalSupply (big.Int) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			t.TotalSupply = new(big.Int)
			if err := t.TotalSupply.UnmarshalCBOR(cr); err != nil {
				return xerrors.Errorf("unmarshaling t.TotalSupply pointer: %w", err)
			}
		}

	}
	return nil
}

var lengthBufTransferReq = []byte{130}

func (t *TransferReq) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufTransferReq); err != nil {
		return err
	}

	// t.ReceiverAddr (address.Address) (struct)
	if err := t.ReceiverAddr.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.TransferAmount (big.Int) (struct)
	if err := t.TransferAmount.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *TransferReq) UnmarshalCBOR(r io.Reader) (err error) {
	*t = TransferReq{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.ReceiverAddr (address.Address) (struct)

	{

		if err := t.ReceiverAddr.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.ReceiverAddr: %w", err)
		}

	}
	// t.TransferAmount (big.Int) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			t.TransferAmount = new(big.Int)
			if err := t.TransferAmount.UnmarshalCBOR(cr); err != nil {
				return xerrors.Errorf("unmarshaling t.TransferAmount pointer: %w", err)
			}
		}

	}
	return nil
}

var lengthBufAllowanceReq = []byte{130}

func (t *AllowanceReq) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufAllowanceReq); err != nil {
		return err
	}

	// t.OwnerAddr (address.Address) (struct)
	if err := t.OwnerAddr.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.SpenderAddr (address.Address) (struct)
	if err := t.SpenderAddr.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *AllowanceReq) UnmarshalCBOR(r io.Reader) (err error) {
	*t = AllowanceReq{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.OwnerAddr (address.Address) (struct)

	{

		if err := t.OwnerAddr.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.OwnerAddr: %w", err)
		}

	}
	// t.SpenderAddr (address.Address) (struct)

	{

		if err := t.SpenderAddr.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.SpenderAddr: %w", err)
		}

	}
	return nil
}

var lengthBufTransferFromReq = []byte{131}

func (t *TransferFromReq) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufTransferFromReq); err != nil {
		return err
	}

	// t.OwnerAddr (address.Address) (struct)
	if err := t.OwnerAddr.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.ReceiverAddr (address.Address) (struct)
	if err := t.ReceiverAddr.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.TransferAmount (big.Int) (struct)
	if err := t.TransferAmount.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *TransferFromReq) UnmarshalCBOR(r io.Reader) (err error) {
	*t = TransferFromReq{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 3 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.OwnerAddr (address.Address) (struct)

	{

		if err := t.OwnerAddr.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.OwnerAddr: %w", err)
		}

	}
	// t.ReceiverAddr (address.Address) (struct)

	{

		if err := t.ReceiverAddr.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.ReceiverAddr: %w", err)
		}

	}
	// t.TransferAmount (big.Int) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			t.TransferAmount = new(big.Int)
			if err := t.TransferAmount.UnmarshalCBOR(cr); err != nil {
				return xerrors.Errorf("unmarshaling t.TransferAmount pointer: %w", err)
			}
		}

	}
	return nil
}

var lengthBufApprovalReq = []byte{130}

func (t *ApprovalReq) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufApprovalReq); err != nil {
		return err
	}

	// t.SpenderAddr (address.Address) (struct)
	if err := t.SpenderAddr.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.NewAllowance (big.Int) (struct)
	if err := t.NewAllowance.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *ApprovalReq) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ApprovalReq{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.SpenderAddr (address.Address) (struct)

	{

		if err := t.SpenderAddr.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.SpenderAddr: %w", err)
		}

	}
	// t.NewAllowance (big.Int) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			t.NewAllowance = new(big.Int)
			if err := t.NewAllowance.UnmarshalCBOR(cr); err != nil {
				return xerrors.Errorf("unmarshaling t.NewAllowance pointer: %w", err)
			}
		}

	}
	return nil
}

var lengthBufFakeSetBalance = []byte{130}

func (t *FakeSetBalance) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufFakeSetBalance); err != nil {
		return err
	}

	// t.Addr (address.Address) (struct)
	if err := t.Addr.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Balance (big.Int) (struct)
	if err := t.Balance.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *FakeSetBalance) UnmarshalCBOR(r io.Reader) (err error) {
	*t = FakeSetBalance{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Addr (address.Address) (struct)

	{

		if err := t.Addr.UnmarshalCBOR(cr); err != nil {
			return xerrors.Errorf("unmarshaling t.Addr: %w", err)
		}

	}
	// t.Balance (big.Int) (struct)

	{

		b, err := cr.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := cr.UnreadByte(); err != nil {
				return err
			}
			t.Balance = new(big.Int)
			if err := t.Balance.UnmarshalCBOR(cr); err != nil {
				return xerrors.Errorf("unmarshaling t.Balance pointer: %w", err)
			}
		}

	}
	return nil
}

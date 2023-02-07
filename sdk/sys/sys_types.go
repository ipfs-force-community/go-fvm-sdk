package sys

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"

	"github.com/filecoin-project/go-state-types/abi"

	stdBig "math/big"

	"github.com/filecoin-project/go-state-types/big"
)

type networkContext_ struct {
	// The current epoch.
	Epoch abi.ChainEpoch
	// The current time (seconds since the unix epoch).
	Timestamp uint64
	// The current base-fee.
	BaseFee fvmTokenAmount
	// The Chain ID of the network.
	ChainId uint64
	// The network version.
	NetworkVersion uint32
}

type messageContext struct {
	// The current call's origin actor ID.
	Origin abi.ActorID
	// The nonce from the explicit message.
	Nonce uint64
	// The caller's actor ID.
	Caller abi.ActorID
	// The receiver's actor ID (i.e. ourselves).
	Receiver abi.ActorID
	// The method number from the message.
	MethodNumber abi.MethodNum
	// The value that was received.
	ValueReceived fvmTokenAmount
	// The current gas premium
	GasPremium fvmTokenAmount
	// Flags pertaining to the currently executing actor's invocation context.
	Flags uint64
}

// fvmTokenAmount use this amount to receive value from fvm
type fvmTokenAmount struct {
	Lo uint64
	Hi uint64
}

// IsZero returns true if u == 0.
func (u fvmTokenAmount) IsZero() bool {
	// NOTE: we do not compare against Zero, because that is a global variable
	// that could be modified.
	return u == fvmTokenAmount{}
}

// Equals returns true if u == v.
//
// Uint128 values can be compared directly with ==, but use of the Equals method
// is preferred for consistency.
func (u fvmTokenAmount) Equals(v fvmTokenAmount) bool {
	return u == v
}

// Equals64 returns true if u == v.
func (u fvmTokenAmount) Equals64(v uint64) bool {
	return u.Lo == v && u.Hi == 0
}

// Cmp compares u and v and returns:
//
//	-1 if u <  v
//	 0 if u == v
//	+1 if u >  v
func (u fvmTokenAmount) Cmp(v fvmTokenAmount) int {
	if u == v {
		return 0
	} else if u.Hi < v.Hi || (u.Hi == v.Hi && u.Lo < v.Lo) {
		return -1
	} else {
		return 1
	}
}

// Cmp64 compares u and v and returns:
//
//	-1 if u <  v
//	 0 if u == v
//	+1 if u >  v
func (u fvmTokenAmount) Cmp64(v uint64) int {
	if u.Hi == 0 && u.Lo == v {
		return 0
	} else if u.Hi == 0 && u.Lo < v {
		return -1
	} else {
		return 1
	}
}

// And returns u&v.
func (u fvmTokenAmount) And(v fvmTokenAmount) fvmTokenAmount {
	return fvmTokenAmount{u.Lo & v.Lo, u.Hi & v.Hi}
}

// And64 returns u&v.
func (u fvmTokenAmount) And64(v uint64) fvmTokenAmount {
	return fvmTokenAmount{u.Lo & v, 0}
}

// Or returns u|v.
func (u fvmTokenAmount) Or(v fvmTokenAmount) fvmTokenAmount {
	return fvmTokenAmount{u.Lo | v.Lo, u.Hi | v.Hi}
}

// Or64 returns u|v.
func (u fvmTokenAmount) Or64(v uint64) fvmTokenAmount {
	return fvmTokenAmount{u.Lo | v, 0}
}

// Xor returns u^v.
func (u fvmTokenAmount) Xor(v fvmTokenAmount) fvmTokenAmount {
	return fvmTokenAmount{u.Lo ^ v.Lo, u.Hi ^ v.Hi}
}

// Xor64 returns u^v.
func (u fvmTokenAmount) Xor64(v uint64) fvmTokenAmount {
	return fvmTokenAmount{u.Lo ^ v, 0}
}

// Add returns u+v.
func (u fvmTokenAmount) Add(v fvmTokenAmount) fvmTokenAmount {
	lo, carry := bits.Add64(u.Lo, v.Lo, 0)
	hi, carry := bits.Add64(u.Hi, v.Hi, carry)
	if carry != 0 {
		panic("overflow")
	}
	return fvmTokenAmount{lo, hi}
}

// AddWrap returns u+v with wraparound semantics; for example,
// Max.AddWrap(From64(1)) == Zero.
func (u fvmTokenAmount) AddWrap(v fvmTokenAmount) fvmTokenAmount {
	lo, carry := bits.Add64(u.Lo, v.Lo, 0)
	hi, _ := bits.Add64(u.Hi, v.Hi, carry)
	return fvmTokenAmount{lo, hi}
}

// Add64 returns u+v.
func (u fvmTokenAmount) Add64(v uint64) fvmTokenAmount {
	lo, carry := bits.Add64(u.Lo, v, 0)
	hi, carry := bits.Add64(u.Hi, 0, carry)
	if carry != 0 {
		panic("overflow")
	}
	return fvmTokenAmount{lo, hi}
}

// AddWrap64 returns u+v with wraparound semantics; for example,
// Max.AddWrap64(1) == Zero.
func (u fvmTokenAmount) AddWrap64(v uint64) fvmTokenAmount {
	lo, carry := bits.Add64(u.Lo, v, 0)
	hi := u.Hi + carry
	return fvmTokenAmount{lo, hi}
}

// Sub returns u-v.
func (u fvmTokenAmount) Sub(v fvmTokenAmount) fvmTokenAmount {
	lo, borrow := bits.Sub64(u.Lo, v.Lo, 0)
	hi, borrow := bits.Sub64(u.Hi, v.Hi, borrow)
	if borrow != 0 {
		panic("underflow")
	}
	return fvmTokenAmount{lo, hi}
}

// SubWrap returns u-v with wraparound semantics; for example,
// Zero.SubWrap(From64(1)) == Max.
func (u fvmTokenAmount) SubWrap(v fvmTokenAmount) fvmTokenAmount {
	lo, borrow := bits.Sub64(u.Lo, v.Lo, 0)
	hi, _ := bits.Sub64(u.Hi, v.Hi, borrow)
	return fvmTokenAmount{lo, hi}
}

// Sub64 returns u-v.
func (u fvmTokenAmount) Sub64(v uint64) fvmTokenAmount {
	lo, borrow := bits.Sub64(u.Lo, v, 0)
	hi, borrow := bits.Sub64(u.Hi, 0, borrow)
	if borrow != 0 {
		panic("underflow")
	}
	return fvmTokenAmount{lo, hi}
}

// SubWrap64 returns u-v with wraparound semantics; for example,
// Zero.SubWrap64(1) == Max.
func (u fvmTokenAmount) SubWrap64(v uint64) fvmTokenAmount {
	lo, borrow := bits.Sub64(u.Lo, v, 0)
	hi := u.Hi - borrow
	return fvmTokenAmount{lo, hi}
}

// Mul returns u*v, panicking on overflow.
func (u fvmTokenAmount) Mul(v fvmTokenAmount) fvmTokenAmount {
	hi, lo := bits.Mul64(u.Lo, v.Lo)
	p0, p1 := bits.Mul64(u.Hi, v.Lo)
	p2, p3 := bits.Mul64(u.Lo, v.Hi)
	hi, c0 := bits.Add64(hi, p1, 0)
	hi, c1 := bits.Add64(hi, p3, c0)
	if (u.Hi != 0 && v.Hi != 0) || p0 != 0 || p2 != 0 || c1 != 0 {
		panic("overflow")
	}
	return fvmTokenAmount{lo, hi}
}

// MulWrap returns u*v with wraparound semantics; for example,
// Max.MulWrap(Max) == 1.
func (u fvmTokenAmount) MulWrap(v fvmTokenAmount) fvmTokenAmount {
	hi, lo := bits.Mul64(u.Lo, v.Lo)
	hi += u.Hi*v.Lo + u.Lo*v.Hi
	return fvmTokenAmount{lo, hi}
}

// Mul64 returns u*v, panicking on overflow.
func (u fvmTokenAmount) Mul64(v uint64) fvmTokenAmount {
	hi, lo := bits.Mul64(u.Lo, v)
	p0, p1 := bits.Mul64(u.Hi, v)
	hi, c0 := bits.Add64(hi, p1, 0)
	if p0 != 0 || c0 != 0 {
		panic("overflow")
	}
	return fvmTokenAmount{lo, hi}
}

// MulWrap64 returns u*v with wraparound semantics; for example,
// Max.MulWrap64(2) == Max.Sub64(1).
func (u fvmTokenAmount) MulWrap64(v uint64) fvmTokenAmount {
	hi, lo := bits.Mul64(u.Lo, v)
	hi += u.Hi * v
	return fvmTokenAmount{lo, hi}
}

// Div returns u/v.
func (u fvmTokenAmount) Div(v fvmTokenAmount) fvmTokenAmount {
	q, _ := u.QuoRem(v)
	return q
}

// Div64 returns u/v.
func (u fvmTokenAmount) Div64(v uint64) fvmTokenAmount {
	q, _ := u.QuoRem64(v)
	return q
}

// QuoRem returns q = u/v and r = u%v.
func (u fvmTokenAmount) QuoRem(v fvmTokenAmount) (q, r fvmTokenAmount) {
	if v.Hi == 0 {
		var r64 uint64
		q, r64 = u.QuoRem64(v.Lo)
		r = From64(r64)
	} else {
		// generate a "trial quotient," guaranteed to be within 1 of the actual
		// quotient, then adjust.
		n := uint(bits.LeadingZeros64(v.Hi))
		v1 := v.Lsh(n)
		u1 := u.Rsh(1)
		tq, _ := bits.Div64(u1.Hi, u1.Lo, v1.Hi)
		tq >>= 63 - n
		if tq != 0 {
			tq--
		}
		q = From64(tq)
		// calculate remainder using trial quotient, then adjust if remainder is
		// greater than divisor
		r = u.Sub(v.Mul64(tq))
		if r.Cmp(v) >= 0 {
			q = q.Add64(1)
			r = r.Sub(v)
		}
	}
	return
}

// QuoRem64 returns q = u/v and r = u%v.
func (u fvmTokenAmount) QuoRem64(v uint64) (q fvmTokenAmount, r uint64) {
	if u.Hi < v {
		q.Lo, r = bits.Div64(u.Hi, u.Lo, v)
	} else {
		q.Hi, r = bits.Div64(0, u.Hi, v)
		q.Lo, r = bits.Div64(r, u.Lo, v)
	}
	return
}

// Mod returns r = u%v.
func (u fvmTokenAmount) Mod(v fvmTokenAmount) (r fvmTokenAmount) {
	_, r = u.QuoRem(v)
	return
}

// Mod64 returns r = u%v.
func (u fvmTokenAmount) Mod64(v uint64) (r uint64) {
	_, r = u.QuoRem64(v)
	return
}

// Lsh returns u<<n.
func (u fvmTokenAmount) Lsh(n uint) (s fvmTokenAmount) {
	if n > 64 {
		s.Lo = 0
		s.Hi = u.Lo << (n - 64)
	} else {
		s.Lo = u.Lo << n
		s.Hi = u.Hi<<n | u.Lo>>(64-n)
	}
	return
}

// Rsh returns u>>n.
func (u fvmTokenAmount) Rsh(n uint) (s fvmTokenAmount) {
	if n > 64 {
		s.Lo = u.Hi >> (n - 64)
		s.Hi = 0
	} else {
		s.Lo = u.Lo>>n | u.Hi<<(64-n)
		s.Hi = u.Hi >> n
	}
	return
}

// LeadingZeros returns the number of leading zero bits in u; the result is 128
// for u == 0.
func (u fvmTokenAmount) LeadingZeros() int {
	if u.Hi > 0 {
		return bits.LeadingZeros64(u.Hi)
	}
	return 64 + bits.LeadingZeros64(u.Lo)
}

// TrailingZeros returns the number of trailing zero bits in u; the result is
// 128 for u == 0.
func (u fvmTokenAmount) TrailingZeros() int {
	if u.Lo > 0 {
		return bits.TrailingZeros64(u.Lo)
	}
	return 64 + bits.TrailingZeros64(u.Hi)
}

// OnesCount returns the number of one bits ("population count") in u.
func (u fvmTokenAmount) OnesCount() int {
	return bits.OnesCount64(u.Hi) + bits.OnesCount64(u.Lo)
}

// RotateLeft returns the value of u rotated left by (k mod 128) bits.
func (u fvmTokenAmount) RotateLeft(k int) fvmTokenAmount {
	const n = 128
	s := uint(k) & (n - 1)
	return u.Lsh(s).Or(u.Rsh(n - s))
}

// RotateRight returns the value of u rotated left by (k mod 128) bits.
func (u fvmTokenAmount) RotateRight(k int) fvmTokenAmount {
	return u.RotateLeft(-k)
}

// Reverse returns the value of u with its bits in reversed order.
func (u fvmTokenAmount) Reverse() fvmTokenAmount {
	return fvmTokenAmount{bits.Reverse64(u.Hi), bits.Reverse64(u.Lo)}
}

// ReverseBytes returns the value of u with its bytes in reversed order.
func (u fvmTokenAmount) ReverseBytes() fvmTokenAmount {
	return fvmTokenAmount{bits.ReverseBytes64(u.Hi), bits.ReverseBytes64(u.Lo)}
}

// Len returns the minimum number of bits required to represent u; the result is
// 0 for u == 0.
func (u fvmTokenAmount) Len() int {
	return 128 - u.LeadingZeros()
}

// String returns the base-10 representation of u as a string.
func (u fvmTokenAmount) String() string {
	if u.IsZero() {
		return "0"
	}
	buf := []byte("0000000000000000000000000000000000000000") // log10(2^128) < 40
	for i := len(buf); ; i -= 19 {
		q, r := u.QuoRem64(1e19) // largest power of 10 that fits in a uint64
		var n int
		for ; r != 0; r /= 10 {
			n++
			buf[i-n] += byte(r % 10)
		}
		if q.IsZero() {
			return string(buf[i-n:])
		}
		u = q
	}
}

// PutBytes stores u in b in little-endian order. It panics if len(b) < 16.
func (u fvmTokenAmount) PutBytes(b []byte) {
	binary.LittleEndian.PutUint64(b[:8], u.Lo)
	binary.LittleEndian.PutUint64(b[8:], u.Hi)
}

// Scan implements fmt.Scanner.
func (u *fvmTokenAmount) Scan(s fmt.ScanState, ch rune) error {
	i := new(big.Int)
	if err := i.Scan(s, ch); err != nil {
		return err
	} else if i.Sign() < 0 {
		return errors.New("value cannot be negative")
	} else if i.BitLen() > 128 {
		return errors.New("value overflows fvmTokenAmount")
	}
	u.Lo = i.Uint64()
	u.Hi = i.Rsh(i.Int, 64).Uint64()
	return nil
}

// New returns the fvmTokenAmount value (lo,hi).
func New(lo, hi uint64) fvmTokenAmount {
	return fvmTokenAmount{lo, hi}
}

// From64 converts v to a fvmTokenAmount value.
func From64(v uint64) fvmTokenAmount {
	return New(v, 0)
}

// FromBytes converts b to a fvmTokenAmount value.
func FromBytes(b []byte) fvmTokenAmount {
	return New(
		binary.LittleEndian.Uint64(b[:8]),
		binary.LittleEndian.Uint64(b[8:]),
	)
}

// TokenAmount returns u as a *ab.
func (u fvmTokenAmount) TokenAmount() *abi.TokenAmount {
	i := new(stdBig.Int).SetUint64(u.Hi)
	i = i.Lsh(i, 64)
	i = i.Xor(i, new(stdBig.Int).SetUint64(u.Lo))
	return &abi.TokenAmount{
		Int: i,
	}
}

// FromBig converts i to a fvmTokenAmount value. It panics if i is negative or
// overflows 128 bits.
func FromBig(i *big.Int) (u fvmTokenAmount) {
	if i.Sign() < 0 {
		panic("value cannot be negative")
	} else if i.BitLen() > 128 {
		panic("value overflows fvmTokenAmount")
	}
	u.Lo = i.Uint64()
	u.Hi = i.Rsh(i.Int, 64).Uint64()
	return u
}

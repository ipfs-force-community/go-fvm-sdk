package rlepluslazy

import (
	"math"
	"math/bits"

	"golang.org/x/xerrors"
)

type decodeInfo struct {
	length byte // length of the run
	i      byte // i+1 is number of repeats of above run lengths
	n      byte // number of bits to read
	varint bool // varint signifies that futher bits need to be processed as a varint
}

func init() {
	buildDecodeTable()
}

// This is a LUT for all possible 6 bit codes and what they decode into
// possible combinations are:
// 0bxxxxxx1 - 1 run of 1
// 0bxxxxx11 - 2 runs of 1
// up to 0b111111 - 6 runs of 1
// 0bAAAA10 - 1 run of length 0bAAAA
// 0bxxxx00 - varint run, the decode value not defined in LUT
var decodeTable = [1 << 6]decodeInfo{}

func buildDecodeTable() {
	for idx := uint8(0); int(idx) < len(decodeTable); idx++ {
		switch {
		case bits.TrailingZeros8(^idx) > 0:
			i := uint8(bits.TrailingZeros8(^idx))
			decodeTable[idx] = decodeInfo{
				length: 1,
				i:      i - 1,
				n:      i,
			}
		case idx&0b11 == 0b10:
			// 01 + 4bit : run of 0 to 15
			decodeTable[idx] = decodeInfo{
				length: byte(idx >> 2),
				i:      0,
				n:      6,
			}
		case idx&0b11 == 0b00:
			decodeTable[idx] = decodeInfo{
				n:      2,
				varint: true,
			}
		}
	}
}

func DecodeRLE(buf []byte) (RunIterator, error) {
	if len(buf) > 0 && buf[len(buf)-1] == 0 {
		// trailing zeros bytes not allowed.
		return nil, xerrors.Errorf("not minimally encoded: %w", ErrDecode)
	}

	bv := readBitvec(buf)

	ver := bv.Get(2) // Read version
	if ver != Version {
		return nil, ErrWrongVersion
	}

	it := &rleIterator{bv: bv}

	// next run is previous in relation to prep
	// so we invert the value
	it.lastVal = bv.Get(1) != 1
	if err := it.prep(); err != nil {
		return nil, err
	}
	return it, nil
}

// ValidateRLE validates the RLE+ in buf does not overflow Uint64
func ValidateRLE(buf []byte) error {
	if len(buf) > 0 && buf[len(buf)-1] == 0 {
		// trailing zeros bytes not allowed.
		return xerrors.Errorf("not minimally encoded: %w", ErrDecode)
	}
	bv := readBitvec(buf)

	ver := bv.Get(2) // Read version
	if ver != Version {
		return ErrWrongVersion
	}

	// this is run value bit, as we are validating lengths we don't care about it
	bv.Get(1)

	totalLen := uint64(0)
	for {
		idx := bv.Peek6()
		decode := decodeTable[idx]
		_ = bv.Get(decode.n)

		var runLen uint64
		if decode.varint {
			x, err := decodeBFVarint(bv)
			if err != nil {
				return err
			}
			runLen = x
		} else {
			runLen = uint64(decode.i+1) * uint64(decode.length)
		}

		if math.MaxUint64-runLen < totalLen {
			return xerrors.Errorf("RLE+ overflow")
		}
		totalLen += runLen
		if runLen == 0 {
			break
		}
	}
	return nil
}

type rleIterator struct {
	bv     *rbitvec
	length uint64

	lastVal bool
	i       uint8
}

func (it *rleIterator) HasNext() bool {
	return it.length != 0
}

func (it *rleIterator) NextRun() (r Run, err error) {
	ret := Run{Len: it.length, Val: !it.lastVal}
	it.lastVal = ret.Val

	if it.i == 0 {
		err = it.prep()
	} else {
		it.i--
	}
	return ret, err
}

func decodeBFVarint(bv *rbitvec) (uint64, error) {
	// Modified from the go standard library. Copyright the Go Authors and
	// released under the BSD License.
	var x uint64
	var s uint
	for i := 0; ; i++ {
		if i == 10 {
			return 0, xerrors.Errorf("run too long: %w", ErrDecode)
		}
		b := bv.GetByte()
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, xerrors.Errorf("run too long: %w", ErrDecode)
			} else if b == 0 && s > 0 {
				return 0, xerrors.Errorf("invalid run: %w", ErrDecode)
			}
			x |= uint64(b) << s
			break
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return x, nil
}

func (it *rleIterator) prep() error {
	idx := it.bv.Peek6()
	decode := decodeTable[idx]
	_ = it.bv.Get(decode.n)

	it.i = decode.i
	it.length = uint64(decode.length)
	if decode.varint {
		x, err := decodeBFVarint(it.bv)
		if err != nil {
			return err
		}
		it.length = x
	}
	return nil
}

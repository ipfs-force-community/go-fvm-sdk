package nonffi

import (
	"errors"
	"fmt"
	"math/bits"

	"github.com/filecoin-project/go-commp-utils/zerocomm"
	commcid "github.com/filecoin-project/go-fil-commcid"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	sha256simd "github.com/minio/sha256-simd"
)

type stackFrame struct {
	size  uint64
	commP []byte
}

func GenerateUnsealedCID(proofType abi.RegisteredSealProof, pieceInfos []abi.PieceInfo) (cid.Cid, error) {
	spi, found := abi.SealProofInfos[proofType]
	if !found {
		return cid.Undef, fmt.Errorf("unknown seal proof type %d", proofType)
	}
	if len(pieceInfos) == 0 {
		return cid.Undef, errors.New("no pieces provided")
	}

	maxSize := uint64(spi.SectorSize)

	todo := make([]stackFrame, len(pieceInfos))

	// sancheck everything
	for i, p := range pieceInfos {
		if p.Size < 128 {
			return cid.Undef, fmt.Errorf("invalid Size of PieceInfo %d: value %d is too small", i, p.Size)
		}
		if uint64(p.Size) > maxSize {
			return cid.Undef, fmt.Errorf("invalid Size of PieceInfo %d: value %d is larger than sector size of SealProofType %d", i, p.Size, proofType)
		}
		if bits.OnesCount64(uint64(p.Size)) != 1 {
			return cid.Undef, fmt.Errorf("invalid Size of PieceInfo %d: value %d is not a power of 2", i, p.Size)
		}

		cp, err := commcid.CIDToPieceCommitmentV1(p.PieceCID)
		if err != nil {
			return cid.Undef, fmt.Errorf("invalid PieceCid for PieceInfo %d: %w", i, err)
		}
		todo[i] = stackFrame{size: uint64(p.Size), commP: cp}
	}

	// reimplement https://github.com/filecoin-project/rust-fil-proofs/blob/380d6437c2/filecoin-proofs/src/pieces.rs#L85-L145
	stack := append(
		make(
			[]stackFrame,
			0,
			32,
		),
		todo[0],
	)

	for _, f := range todo[1:] {

		// pre-pad if needed to balance the left limb
		for stack[len(stack)-1].size < f.size {
			lastSize := stack[len(stack)-1].size

			stack = reduceStack(
				append(
					stack,
					stackFrame{
						size:  lastSize,
						commP: zeroCommForSize(lastSize),
					},
				),
			)
		}

		stack = reduceStack(
			append(
				stack,
				f,
			),
		)
	}

	for len(stack) > 1 {
		lastSize := stack[len(stack)-1].size
		stack = reduceStack(
			append(
				stack,
				stackFrame{
					size:  lastSize,
					commP: zeroCommForSize(lastSize),
				},
			),
		)
	}

	if stack[0].size > maxSize {
		return cid.Undef, fmt.Errorf("provided pieces sum up to %d bytes, which is larger than sector size of SealProofType %d", stack[0].size, proofType)
	}

	return commcid.PieceCommitmentV1ToCID(stack[0].commP)
}

var s256 = sha256simd.New()

func zeroCommForSize(s uint64) []byte { return zerocomm.PieceComms[bits.TrailingZeros64(s)-7][:] }

func reduceStack(s []stackFrame) []stackFrame {
	for len(s) > 1 && s[len(s)-2].size == s[len(s)-1].size {

		s256.Reset()
		s256.Write(s[len(s)-2].commP)
		s256.Write(s[len(s)-1].commP)
		d := s256.Sum(make([]byte, 0, 32))
		d[31] &= 0b00111111

		s[len(s)-2] = stackFrame{
			size:  2 * s[len(s)-2].size,
			commP: d,
		}

		s = s[:len(s)-1]
	}

	return s
}

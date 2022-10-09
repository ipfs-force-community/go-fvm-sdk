package rlepluslazy

import (
	"math"
	"math/rand"
)

func NewFromZipfDist(seed int64, size int) RunIterator {
	zipf := rand.NewZipf(rand.New(rand.NewSource(seed)), 1.6978377, 1, math.MaxUint64/(1<<16))
	return &zipfIterator{
		i:    size,
		zipf: zipf,
	}
}

type zipfIterator struct {
	i    int
	zipf *rand.Zipf
}

func (zi *zipfIterator) HasNext() bool {
	return zi.i != 0
}

func (zi *zipfIterator) NextRun() (Run, error) {
	zi.i--
	return Run{
		Len: zi.zipf.Uint64() + 1,
		Val: zi.i%2 == 0,
	}, nil
}

package builtin

import (
	"github.com/filecoin-project/go-state-types/big"
)

///// Code shared by multiple built-in actors. /////

// Default log2 of branching factor for HAMTs.
// This value has been empirically chosen, but the optimal value for maps with different mutation profiles may differ.
const DefaultHamtBitwidth = 5

type BigFrac struct {
	Numerator   big.Int
	Denominator big.Int
}

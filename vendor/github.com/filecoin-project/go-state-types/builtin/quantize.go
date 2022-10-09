package builtin

import "github.com/filecoin-project/go-state-types/abi"

// A spec for quantization.
type QuantSpec struct {
	unit   abi.ChainEpoch // The unit of quantization
	offset abi.ChainEpoch // The offset from zero from which to base the modulus
}

func NewQuantSpec(unit, offset abi.ChainEpoch) QuantSpec {
	return QuantSpec{unit: unit, offset: offset}
}

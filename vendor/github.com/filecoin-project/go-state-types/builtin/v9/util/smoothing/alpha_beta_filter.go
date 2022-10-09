package smoothing

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v9/util/math"
)

var (
	DefaultAlpha big.Int // Q.128 value of 9.25e-4
	DefaultBeta  big.Int // Q.128 value of 2.84e-7

	ExtrapolatedCumSumRatioEpsilon big.Int // Q.128 value of 2^-50
)

func init() {
	// Alpha Beta Filter constants
	constStrs := []string{
		"314760000000000000000000000000000000", // DefaultAlpha
		"96640100000000000000000000000000",     // DefaultBeta
		"302231454903657293676544",             // Epsilon

	}
	constBigs := math.Parse(constStrs)
	DefaultAlpha = big.NewFromGo(constBigs[0])
	DefaultBeta = big.NewFromGo(constBigs[1])
	ExtrapolatedCumSumRatioEpsilon = big.NewFromGo(constBigs[2])

}

//Alpha Beta Filter "position" (value) and "velocity" (rate of change of value) estimates
//Estimates are in Q.128 format
type FilterEstimate struct {
	PositionEstimate big.Int // Q.128
	VelocityEstimate big.Int // Q.128
}

// Returns the Q.0 position estimate of the filter
func Estimate(fe *FilterEstimate) big.Int {
	return big.Rsh(fe.PositionEstimate, math.Precision128) // Q.128 => Q.0
}

// Create a new filter estimate given two Q.0 format ints.
func NewEstimate(position, velocity big.Int) FilterEstimate {
	return FilterEstimate{
		PositionEstimate: big.Lsh(position, math.Precision128), // Q.0 => Q.128
		VelocityEstimate: big.Lsh(velocity, math.Precision128), // Q.0 => Q.128
	}
}

// Extrapolate the CumSumRatio given two filters.
// Output is in Q.128 format
func ExtrapolatedCumSumOfRatio(delta abi.ChainEpoch, relativeStart abi.ChainEpoch, estimateNum, estimateDenom FilterEstimate) big.Int {
	deltaT := big.Lsh(big.NewInt(int64(delta)), math.Precision128)     // Q.0 => Q.128
	t0 := big.Lsh(big.NewInt(int64(relativeStart)), math.Precision128) // Q.0 => Q.128
	// Renaming for ease of following spec and clarity
	position1 := estimateNum.PositionEstimate
	position2 := estimateDenom.PositionEstimate
	velocity1 := estimateNum.VelocityEstimate
	velocity2 := estimateDenom.VelocityEstimate

	squaredVelocity2 := big.Mul(velocity2, velocity2)               // Q.128 * Q.128 => Q.256
	squaredVelocity2 = big.Rsh(squaredVelocity2, math.Precision128) // Q.256 => Q.128

	if squaredVelocity2.GreaterThan(ExtrapolatedCumSumRatioEpsilon) {
		x2a := big.Mul(t0, velocity2)         // Q.128 * Q.128 => Q.256
		x2a = big.Rsh(x2a, math.Precision128) // Q.256 => Q.128
		x2a = big.Sum(position2, x2a)

		x2b := big.Mul(deltaT, velocity2)     // Q.128 * Q.128 => Q.256
		x2b = big.Rsh(x2b, math.Precision128) // Q.256 => Q.128
		x2b = big.Sum(x2a, x2b)

		x2a = math.Ln(x2a) // Q.128
		x2b = math.Ln(x2b) // Q.128

		m1 := big.Sub(x2b, x2a)
		m1 = big.Mul(velocity2, big.Mul(position1, m1)) // Q.128 * Q.128 * Q.128 => Q.384
		m1 = big.Rsh(m1, math.Precision128)             //Q.384 => Q.256

		m2L := big.Sub(x2a, x2b)
		m2L = big.Mul(position2, m2L)     // Q.128 * Q.128 => Q.256
		m2R := big.Mul(velocity2, deltaT) // Q.128 * Q.128 => Q.256
		m2 := big.Sum(m2L, m2R)
		m2 = big.Mul(velocity1, m2)         // Q.256 => Q.384
		m2 = big.Rsh(m2, math.Precision128) //Q.384 => Q.256

		return big.Div(big.Sum(m1, m2), squaredVelocity2) // Q.256 / Q.128 => Q.128

	}

	halfDeltaT := big.Rsh(deltaT, 1)                   // Q.128 / Q.0 => Q.128
	x1m := big.Mul(velocity1, big.Sum(t0, halfDeltaT)) // Q.128 * Q.128 => Q.256
	x1m = big.Rsh(x1m, math.Precision128)              // Q.256 => Q.128
	x1m = big.Add(position1, x1m)

	cumsumRatio := big.Mul(x1m, deltaT)           // Q.128 * Q.128 => Q.256
	cumsumRatio = big.Div(cumsumRatio, position2) // Q.256 / Q.128 => Q.128
	return cumsumRatio

}

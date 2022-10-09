package market

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
)

// The percentage of normalized cirulating
// supply that must be covered by provider collateral in a deal
var ProviderCollateralSupplyTarget = builtin.BigFrac{
	Numerator:   big.NewInt(1), // PARAM_SPEC
	Denominator: big.NewInt(100),
}

// Minimum deal duration.
var DealMinDuration = abi.ChainEpoch(180 * builtin.EpochsInDay) // PARAM_SPEC

// Maximum deal duration
var DealMaxDuration = abi.ChainEpoch(540 * builtin.EpochsInDay) // PARAM_SPEC

// Bounds (inclusive) on deal duration
func DealDurationBounds(_ abi.PaddedPieceSize) (min abi.ChainEpoch, max abi.ChainEpoch) {
	return DealMinDuration, DealMaxDuration
}

// DealMaxLabelSize is the maximum size of a deal label.
const DealMaxLabelSize = 256

func DealPricePerEpochBounds(_ abi.PaddedPieceSize, _ abi.ChainEpoch) (min abi.TokenAmount, max abi.TokenAmount) {
	return abi.NewTokenAmount(0), builtin.TotalFilecoin
}

func DealProviderCollateralBounds(pieceSize abi.PaddedPieceSize, verified bool, networkRawPower, networkQAPower, baselinePower abi.StoragePower,
	networkCirculatingSupply abi.TokenAmount) (min, max abi.TokenAmount) {
	// minimumProviderCollateral = ProviderCollateralSupplyTarget * normalizedCirculatingSupply
	// normalizedCirculatingSupply = networkCirculatingSupply * dealPowerShare
	// dealPowerShare = dealRawPower / max(BaselinePower(t), NetworkRawPower(t), dealRawPower)

	lockTargetNum := big.Mul(ProviderCollateralSupplyTarget.Numerator, networkCirculatingSupply)
	lockTargetDenom := ProviderCollateralSupplyTarget.Denominator
	powerShareNum := big.NewIntUnsigned(uint64(pieceSize))
	powerShareDenom := big.Max(big.Max(networkRawPower, baselinePower), powerShareNum)

	num := big.Mul(lockTargetNum, powerShareNum)
	denom := big.Mul(lockTargetDenom, powerShareDenom)
	minCollateral := big.Div(num, denom)
	return minCollateral, builtin.TotalFilecoin
}

func DealClientCollateralBounds(_ abi.PaddedPieceSize, _ abi.ChainEpoch) (min abi.TokenAmount, max abi.TokenAmount) {
	return abi.NewTokenAmount(0), builtin.TotalFilecoin
}

// Computes the weight for a deal proposal, which is a function of its size and duration.
func DealWeight(proposal *DealProposal) abi.DealWeight {
	dealDuration := big.NewInt(int64(proposal.Duration()))
	dealSize := big.NewIntUnsigned(uint64(proposal.PieceSize))
	dealSpaceTime := big.Mul(dealDuration, dealSize)
	return dealSpaceTime
}

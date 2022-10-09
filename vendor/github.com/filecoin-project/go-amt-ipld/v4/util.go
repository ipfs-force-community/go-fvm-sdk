package amt

import "math"

// Given height 'height', how many nodes in a maximally full tree can we
// build? (bitWidth^2)^height = width^height. If we pass in height+1 we can work
// out how many elements a maximally full tree can hold, width^(height+1).
func nodesForHeight(bitWidth uint, height int) uint64 {
	heightLogTwo := uint64(bitWidth) * uint64(height)
	if heightLogTwo >= 64 {
		// The max depth layer may not be full.
		return math.MaxUint64
	}
	return 1 << heightLogTwo
}

package simulated

import (
	"testing"
)

func TestBeaconRandomness(t *testing.T) {
	makeRandomness(343, 438, []byte{})
}

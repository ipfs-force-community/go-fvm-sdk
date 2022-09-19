//go:build simulate
// +build simulate

package contract

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSayHello(t *testing.T) {
	simulated.Begin()

	testState := State{}
	sdk.SaveState(&testState)

	newSt := new(State)
	sdk.LoadState(newSt)
	assert.Equal(t, *newSt, testState)
	simulated.End()

}

package contract

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSayHello(t *testing.T) {
	simulated.Begin()
	testState := State{}
	a := testState.SayHello()
	assert.Equal(t, string(a), "1", "The two words should be the same.")

	b := testState.SayHello()
	assert.Equal(t, string(b), "2", "The two words should be the same.")

	simulated.End()
}

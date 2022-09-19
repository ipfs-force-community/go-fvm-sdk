package contract

import (
	"fmt"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"testing"
)

func TestSayHello(t *testing.T) {
	simulated.Begin()
	testState := State{}
	result := testState.SayHello()
	fmt.Printf("%s\n", result)
	simulated.End()
}

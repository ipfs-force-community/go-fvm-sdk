//go:build simulate

package contract

import (
	"fmt"
	gomock "github.com/golang/mock/gomock"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
	"testing"
)

func TestSayHello(t *testing.T) {
	simulated.Begin()

	erc20 := makeErc20Token()
	sdk.SaveState(&erc20)

	newSt := new(Erc20Token)
	sdk.LoadState(newSt)
	assert.Equal(t, *newSt, erc20)
	simulated.End()

}

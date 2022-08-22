//go:build simulate

package contract

import (
	"fmt"
	"testing"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs/go-cid"
)

func init() {

}

func TestErc20Token_Approval(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%v\n", err)
		}
	}()

	sys.SimulatedInstance.EXPECT().Enabled().Return(true, nil)
	LoggerInit()
	_ = Erc20Token{}
	sys.Finish()

}

func TestErc20Token_SelfRoot(t *testing.T) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%v\n", err)
		}
	}()

	sys.SimulatedInstance.EXPECT().Enabled().Return(true, nil)

	cidout, _ := cid.Decode("bafy2bzacecdjkk2tzogitpcybu3eszr4uptrjogstqmyt6u4q2p3hh4chmf3i")
	sys.SimulatedInstance.EXPECT().SelfRoot().Return(cidout, nil)
	token := Erc20Token{}
	sys.Finish()

	got, err := token.SelfRoot()
	fmt.Printf("%v-%v\n", got, err)
	sys.Finish()
}

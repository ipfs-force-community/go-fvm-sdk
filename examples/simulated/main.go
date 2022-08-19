//go:build simulate

package simulated

import (
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs/go-cid"
)

func main() {

}

func Open(id cid.Cid) {
	out, err := sys.Open(id)
	fmt.Printf("%v-%v\n", out, err)
}

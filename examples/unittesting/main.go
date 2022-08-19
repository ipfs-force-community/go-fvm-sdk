package main

import (
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/fvm"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

type FakeReporter struct {
}

func (f *FakeReporter) Errorf(format string, args ...interface{}) {
	println("Errorf\n")
}
func (f *FakeReporter) Fatalf(format string, args ...interface{}) {
	fmt.Printf("Fatalf:\n%v\n", args...)
}

func main() {
	// report := FakeReporter{}
	// fvm.InitMockFvm(&report)


	result := types.IpldOpen{1, 2, 3}
	h, _ := mh.Sum([]byte("TEST"), mh.SHA3, 4)
	in := cid.NewCidV1(7, h)
	println(in.String())



	fvm.OpenExpect(in, &result, nil)


	h1, _ := mh.Sum([]byte("TEST"), mh.SHA3, 4)
	in1 := cid.NewCidV1(7, h1)
	println(h1.String())
	out, _ := fvm.MockFvmInstance.Open(in1)
	fmt.Printf("%v\n", out)



	fvm.EpochFinish()
}

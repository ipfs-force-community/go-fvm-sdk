//go:build simulate

package simulated

import (
	"testing"

	"github.com/ipfs/go-cid"
)


func TestOpen(t *testing.T) {
	defer func() {
		sys.EpochFinish()
	}()

	out := types.IpldOpen{1, 2, 3}
	args1, _ := mh.Sum([]byte("TEST"), mh.SHA3, 4)
	argsin := cid.NewCidV1(7, h)

	sys.SetOpenExpect(argsin, &out, nil)


	type args struct {
		in  interface{}
		out interface{}
	}
	tests := []struct {
		name string
		args args
	}{
	  {name:"test1",args:args{in:argsin,out:out}}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Open(args1)
		})
	}
}

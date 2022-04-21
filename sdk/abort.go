package sdk

import "github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"

func Abort(code uint32, msg string) {
	sys.Abort(code, msg)
}

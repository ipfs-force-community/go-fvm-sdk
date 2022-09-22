package main

import "github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

func main() {
	amount := &types.TokenAmount{
		Lo: 9,
	}
	amount.Big()
}

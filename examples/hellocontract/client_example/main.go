package main

import (
	"context"
	"encoding/hex"
	"flag"
	"hellocontract/client"
	"io/ioutil"
	"log"

	"github.com/filecoin-project/go-address"

	_ "github.com/filecoin-project/specs-actors/v8/actors/builtin/init"
	v0 "github.com/filecoin-project/venus/venus-shared/api/chain/v0"
)

func main() {
	var ip = flag.String("ip", "", "full node url")
	var token = flag.String("token", "", "full node token")
	var fromAddrStr = flag.String("from", "", "send message from")
	flag.Parse()

	ctx := context.Background()
	v0FullNode, closer, err := v0.DialFullNodeRPC(ctx, *ip, *token, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer closer()

	addr, err := address.NewFromString(*fromAddrStr)
	if err != nil {
		log.Fatalln(err)
		return
	}
	helloClient := client.NewStateClient(v0FullNode, client.SetFromAddressOpt(addr))

	code, err := ioutil.ReadFile("../hellocontract.wasm")
	if err != nil {
		log.Fatalln(err)
		return
	}

	installRet, err := helloClient.Install(ctx, code)
	if err != nil {
		log.Fatalln(err)
		return
	}
	createParams, err := hex.DecodeString("846556656e757361560647005af3107a3fff")
	if err != nil {
		log.Fatalln(err)
		return
	}

	execRet, err := helloClient.CreateActor(ctx, installRet.CodeCid, createParams)
	if err != nil {
		log.Fatalln(err)
		return
	}
	println("actor id", execRet.IDAddress.String())

	ret, err := helloClient.SayHello(ctx)
	if err != nil {
		log.Fatalln(err)
		return
	}
	println(string(ret))
}

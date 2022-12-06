package gen

import (
	"encoding/hex"
	"fmt"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func GenCase(name string, methodNum uint64, input, output cbg.CBORMarshaler, send_from uint32, send_value abi.TokenAmount, expectCode ferrors.ExitCode, expectMsg string) {
	msg := fmt.Sprintf(` {
           "name": "%s",
           "method_num": %d,
           "params": "%s",
           "return_data": "%s",
           "send_from": %d,
           "send_value": %s,
           "expect_code": %d,
           "expect_message":"%s"
         }`, name, methodNum,
		hex.EncodeToString(sdk.MustCborMarshal(input)),
		hex.EncodeToString(sdk.MustCborMarshal(output)),
		send_from,
		send_value.String(),
		expectCode,
		expectMsg,
	)
	println(msg)
}

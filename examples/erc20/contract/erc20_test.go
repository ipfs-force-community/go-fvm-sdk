package contract_test

import (
	"bytes"
	"encoding/hex"
	"erc20/contract"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/go-state-types/big"
)

func TestXX(t *testing.T) {
	v := big.NewInt(100)
	c := contract.ConstructorReq{
		Name:        "mock test",
		Symbol:      "testCoin",
		Decimals:    5,
		TotalSupply: &v,
	}
	r := bytes.NewBufferString("")
	err := c.MarshalCBOR(r)
	assert.Nil(t, err)
	bbb := r.Bytes()
	fmt.Println(hex.EncodeToString(bbb))

	var req contract.ConstructorReq
	err = req.UnmarshalCBOR(bytes.NewReader(bbb))
	assert.Nil(t, err)
}

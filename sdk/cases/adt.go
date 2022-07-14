package main

import (
	"context"
	"fmt"

	"github.com/filecoin-project/go-state-types/big"

	"reflect"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/stretchr/testify/assert"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/adt"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/testing"
)

func main() {} //nolint

//
// If the values are not of like type, the returned strings will be prefixed
// with the type name, and the value will be enclosed in parenthesis similar
// to a type conversion in the Go grammar.
func formatUnequalValues(expected, actual interface{}) (e string, a string) {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		return fmt.Sprintf("%T(%s)", expected, truncatingFormat(expected)),
			fmt.Sprintf("%T(%s)", actual, truncatingFormat(actual))
	}
	return truncatingFormat(expected), truncatingFormat(actual)
}

// truncatingFormat formats the data and truncates it if it's too long.
//
// This helps keep formatted error messages lines from exceeding the
// bufio.MaxScanTokenSize max line length that the go testing framework imposes.
func truncatingFormat(data interface{}) string {
	value := fmt.Sprintf("%#v", data)
	return value
}

//go:export invoke
func Invoke(_ uint32) uint32 { //nolint
	t := testing.NewTestingT()
	defer t.CheckResult()

	emptyArray, err := adt.MakeEmptyArray(adt.AdtStore(context.Background()), adt.BalanceTableBitwidth)
	assert.Nil(t, err)

	emptyArrRoot, err := emptyArray.Root()
	assert.Nil(t, err)
	assert.Equal(t, "bafy2bzaceaa2jny7gkgdwnid4kuldau6bnvgyss5bszo4uy6uikrncvdu5mc2", emptyArrRoot.String())

	val := big.NewInt(100)
	err = emptyArray.Set(10, &val)
	assert.Nil(t, err)

	arrRoot, err := emptyArray.Root()
	assert.Nil(t, err)
	assert.Equal(t, "bafy2bzacebomththj2xbwgezqwseyzb3mruxt6lr4ryeiqkhyke3a632tqjlw", arrRoot.String())

	emptyMap, err := adt.MakeEmptyMap(adt.AdtStore(context.Background()), adt.BalanceTableBitwidth)
	assert.Nil(t, err)
	emptyMapRoot, err := emptyMap.Root()
	assert.Nil(t, err)
	assert.Equal(t, "bafy2bzaceamp42wmmgr2g2ymg46euououzfyck7szknvfacqscohrvaikwfay", emptyMapRoot.String())

	val2 := big.NewInt(100000)
	err = emptyMap.Put(types.StringKey("key"), &val2)
	assert.Nil(t, err)

	mapRoot, err := emptyMap.Root()
	assert.Nil(t, err)
	assert.Equal(t, "bafy2bzacebgqqbr6nie2x3md44qk32tvw7n7emypldcedoiissdkb7itwulve", mapRoot.String())

	return 0
}

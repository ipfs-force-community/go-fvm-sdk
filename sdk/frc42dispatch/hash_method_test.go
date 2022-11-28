package frc42dispatch

import (
	"testing"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/stretchr/testify/assert"
)

func TestHashMethodName(t *testing.T) {
	t.Run("construct", func(t *testing.T) {
		methodNumber, err := GenMethodNumber("Constructor")
		assert.NoError(t, err)
		assert.Equal(t, CONSTRUCTORMETHODNUMBER, methodNumber)
	})

	t.Run("success", func(t *testing.T) {
		methodNm, err := GenMethodNumber("GetName")
		assert.NoError(t, err)
		assert.Equal(t, abi.MethodNum(0xa30674d4), methodNm)
	})

	t.Run("first char is _", func(t *testing.T) {
		methodNumber, err := GenMethodNumber("_Error")
		assert.Nil(t, err)
		assert.Equal(t, abi.MethodNum(1786550896), methodNumber)
	})

	//fail case
	t.Run("empty string", func(t *testing.T) {
		_, err := GenMethodNumber("")
		assert.Error(t, err)
	})

	t.Run("bad method namne", func(t *testing.T) {
		_, err := GenMethodNumber("Bad!Method!Name!")
		assert.Error(t, err)
	})

	t.Run("non ascii", func(t *testing.T) {
		_, err := GenMethodNumber("你好")
		assert.Error(t, err)
	})

	t.Run("first char is lower case", func(t *testing.T) {
		_, err := GenMethodNumber("error")
		assert.Error(t, err)
	})
}

func check(t *testing.T, name string, value abi.MethodNum) {
	t.Run("case:"+name, func(t *testing.T) {
		methodNumber, err := GenMethodNumber(name)
		assert.NoError(t, err, name)
		assert.Equal(t, value, methodNumber, name)
	})
}

func TestCompatableWithRust(t *testing.T) {
	check(t, "Method", abi.MethodNum(0xa20642fc))
	// this case from https://github.com/filecoin-project/filecoin-actor-utils/blob/main/frc42_dispatch/macros/tests/build-success.rs
	check(t, "Name", abi.MethodNum(0x02ea015c))
	check(t, "Symbol", abi.MethodNum(0x7adab63e))
	check(t, "TotalSupply", abi.MethodNum(0x06da7a35))
	check(t, "BalanceOf", abi.MethodNum(0x8710e1ac))
	check(t, "Allowance", abi.MethodNum(0xfaa45236))
	check(t, "IncreaseAllowance", abi.MethodNum(0x69ecb918))
	check(t, "DecreaseAllowance", abi.MethodNum(0x5b286f21))
	check(t, "RevokeAllowance", abi.MethodNum(0xa4d840b1))
	check(t, "Burn", abi.MethodNum(0x5584159a))
	check(t, "TransferFrom", abi.MethodNum(0xd7d4deed))
	check(t, "Transfer", abi.MethodNum(0x04cbf732))
	check(t, "Mint", abi.MethodNum(0x06f84ab2))
}

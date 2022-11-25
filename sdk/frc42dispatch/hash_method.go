// Package frc42dispatch implement frc42 reference https://github.com/filecoin-project/FIPs/blob/master/FRCs/frc-0042.md
package frc42dispatch

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/minio/blake2b-simd"

	"github.com/filecoin-project/go-state-types/abi"

	"golang.org/x/exp/utf8string"
)

// CONSTRUCTORMETHODNAME default constructor name
const CONSTRUCTORMETHODNAME = "Constructor"

// CONSTRUCTORMETHODNUMBER default constructor method number
const CONSTRUCTORMETHODNUMBER abi.MethodNum = 1

// DIGESTCHUNKLENGTH chunk to generate u32 method number
const DIGESTCHUNKLENGTH = 4

// FIRSTMETHODNUMBER Method numbers below FIRST_METHOD_NUMBER are reserved for other use
const FIRSTMETHODNUMBER = 1 << 24

// GenMethodNumber generate method number by method name
func GenMethodNumber(name string) (abi.MethodNum, error) {
	err := checkMethodName(name)
	if err != nil {
		return 0, err
	}
	if name == CONSTRUCTORMETHODNAME {
		return CONSTRUCTORMETHODNUMBER, nil
	}
	methodName := fmt.Sprintf("1|%s", name)
	hasher := blake2b.New512()
	_, err = hasher.Write([]byte(methodName))
	if err != nil {
		return 0, err
	}
	methodHashBytes := hasher.Sum(nil)
	var first, end = 0, 1
	for i := 1; i <= len(methodHashBytes); i++ {
		if i%DIGESTCHUNKLENGTH == 0 {
			methodId := binary.BigEndian.Uint32(methodHashBytes[first:end]) //nolint
			if methodId >= FIRSTMETHODNUMBER {
				return abi.MethodNum(methodId), nil
			}
		}
		end++
	}
	return 0, fmt.Errorf("unable to calculate method id, choose a another method name %s", name)
}

func checkMethodName(name string) error {
	if len(name) == 0 {
		return errors.New("method name is empty")
	}

	s := utf8string.NewString(name)
	if !s.IsASCII() {
		return errors.New("method name must be ascii")
	}

	for i := 0; i < s.RuneCount(); i++ {
		rune := s.At(i)
		if i == 0 {
			if ('A' <= rune && rune <= 'Z') || rune == '_' {
				continue
			}
			return fmt.Errorf("method %s first char must be upper case", name)
		}

		if rune == '_' ||
			(0 <= rune && rune <= 9) ||
			('a' <= rune && rune <= 'z') ||
			('A' <= rune && rune <= 'Z') {
			continue
		}
		return fmt.Errorf("method %s must be range _ | a-z | A-Z", name)
	}

	return nil
}

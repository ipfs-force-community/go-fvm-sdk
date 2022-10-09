package system

import (
	"github.com/filecoin-project/go-state-types/abi"
)

var Methods = []interface{}{
	1: *new(func(interface{}, *abi.EmptyValue) *abi.EmptyValue), // Constructor
}

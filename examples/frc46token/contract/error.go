package contract

import (
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
)

func InsufficientAllowanceError(operator, owner address.Address, delta, allowance abi.TokenAmount) error {
	return fmt.Errorf("%s attempted to utilise %s of allowance %s set by %s %w", operator, delta, allowance, owner, ferrors.USR_INSUFFICIENT_FUNDS)
}

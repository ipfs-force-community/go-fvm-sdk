package simulated

import "errors"

var (
	ErrorIDValid           = errors.New("id is valid")
	ErrorKeyExists         = errors.New("key already exists")
	ErrorKeyMatchSucess    = errors.New("key match is ok")
	ErrorKeyMatchFail      = errors.New("key match is fail")
	ErrorKeyTypeException  = errors.New("key type is except")
	ErrorIllegalArgument   = errors.New("illegal argument")
	ErrorIllegalOperation  = errors.New("illegal operation")
	ErrorLimitExceeded     = errors.New("limit exceeded")
	ErrorAssertionFailed   = errors.New("filecoin assertion failed")
	ErrorInsufficientFunds = errors.New("insufficient funds")
	ErrorNotFound          = errors.New("resource not found")
	ErrorInvalidHandle     = errors.New("invalid ipld block handle")
	ErrorIllegalCid        = errors.New("illegal cid specification")
	ErrorIllegalCodec      = errors.New("illegal ipld codec")
	ErrorSerialization     = errors.New("serialization error")
	ErrorForbidden         = errors.New("operation forbidden")
	ErrorBufferTooSmall    = errors.New("buffer too small")
)

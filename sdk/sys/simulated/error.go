package simulated

import "errors"

var (
	ErrorIdValid       = errors.New("id is valid")
	ErrorNotFound      = errors.New("key is not found ")
	ErrorKeyExists     = errors.New("key already exists")
	ErrorKeyMatchSucess = errors.New("key match is ok")
	ErrorKeyMatchFail = errors.New("key match is fail")
	ErrorKeyTypeException = errors.New("key type is except")
)

package simulated

import "errors"

var (
	ErrorIdValid   = errors.New("id is valid")
	ErrorNotFound  = errors.New("key is not found ")
	ErrorKeyExists = errors.New("key already exists")
)

package fvm

import (
	"github.com/golang/mock/gomock"
)



type ExpectOptions struct {
	Do          func(f interface{})
	DoAndReturn func(f interface{})
	MaxTimes    *int
	MinTimes    *int
	AnyTimes    *int
	Times       *int
}

var (
	MatchAny = gomock.Any
	MatchEq  = gomock.Eq
	MatchNil = gomock.Nil
	MatchLen = gomock.Len
	MatchNot = gomock.Not
)

func initcall(call *gomock.Call, op *ExpectOptions) {
	if op != nil {
		return
	}
	if op.Do != nil {
		call = call.Do(op.Do)
	}
	if op.DoAndReturn != nil {
		call = call.DoAndReturn(op.DoAndReturn)
	}

	if op.MaxTimes != nil {
		call = call.MaxTimes(*op.MaxTimes)

	}

	if op.MinTimes != nil {
		call = call.MinTimes(*op.MinTimes)

	}

	if op.AnyTimes != nil {
		call = call.AnyTimes()
	}

	if op.Times != nil {
		call = call.Times(*op.Times)
	}

}

func OpenExpect(in interface{}, out interface{}, op *ExpectOptions) {
	call := mockFvmInstance.EXPECT().Open(in).Return(out)
	initcall(call, op)
}

package sys

import (
	"unsafe"
)

type sliceHeader struct {
	data uintptr
	len  int
	cap  int
}

type stringHeader struct {
	data uintptr
	len  int
}

func GetSlicePointerAndLen(data interface{}) (uintptr, uint32) {
	s := (*sliceHeader)(unsafe.Pointer(&data))
	return s.data, uint32(s.len)
}

func GetStringPointerAndLen(str string) (uintptr, uint32) {
	s := (*sliceHeader)(unsafe.Pointer(&str))
	return s.data, uint32(s.len)
}

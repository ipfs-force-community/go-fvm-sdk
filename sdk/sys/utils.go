package sys

import (
	"reflect"
	"unsafe"
)

type emptyInterface struct {
	_    uintptr
	word unsafe.Pointer
}

func GetSlicePointerAndLen(data interface{}) (uintptr, uint32) {
	eface := (*emptyInterface)(unsafe.Pointer(&data))
	s := (*reflect.SliceHeader)(eface.word)
	return s.Data, uint32(s.Len)
}

func GetStringPointerAndLen(str string) (uintptr, uint32) {
	s := (*reflect.StringHeader)(unsafe.Pointer(&str))
	return s.Data, uint32(s.Len)
}

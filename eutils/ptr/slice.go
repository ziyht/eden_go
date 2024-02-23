package ptr

import (
	"reflect"
	"unsafe"
)

type slice struct {
}
var sSlice = &slice{}

func Slice() *slice { 
	return sSlice 
}

func (s *slice)SetByteAtOffset(slice unsafe.Pointer, offset int, val byte) {
	header := (*reflect.SliceHeader)(slice)
	SetAtOffset(unsafe.Pointer(header.Data), offset, val)
}

func (s *slice)GetByteAtOffset(slice unsafe.Pointer, offset int) byte {
	header := (*reflect.SliceHeader)(slice)
	return GetAtOffset[byte](unsafe.Pointer(header.Data), offset)
}

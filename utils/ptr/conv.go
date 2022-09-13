package ptr

import (
	"unsafe"
)

func StringToBytes(s string) (b []byte) {
	*(*string)(unsafe.Pointer(&b)) = s
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&b)) + 2*unsafe.Sizeof(&b))) = len(s)
	return
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ValToBytes[T any](val T)[]byte {
	len := unsafe.Sizeof(val)
	out := make([]byte, len)
	copy(unsafe.Slice((*byte)(unsafe.Pointer(&out[0])), len), unsafe.Slice((*byte)(unsafe.Pointer(&val)), len))
	return out
}

func BytesToVal[T any](b []byte)(out T){
	return *(*T)(unsafe.Pointer(&b[0]))
}


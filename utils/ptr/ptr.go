package ptr

import (
	"unsafe"
)

// Be careful by using this, we do not check the input ptr and memory's accessment
func SetAtOffset[T any](ptr unsafe.Pointer, offset int, val T){
	*(*T)(unsafe.Pointer(uintptr(ptr) + uintptr(offset))) = val
}

// Be careful by using this, we do not check the input ptr and memory's accessment 
func GetAtOffset[T any](ptr unsafe.Pointer, offset int) T {
	return *(*T)(unsafe.Pointer(uintptr(ptr) + uintptr(offset)))
}

// Be careful by using this, we do not check the input ptr and memory's accessment 
func SetAtIndex[T any](ptr unsafe.Pointer, index int, val T){
	*(*T)(unsafe.Pointer(uintptr(ptr) + uintptr(index) * unsafe.Sizeof(val) )) = val
}

// Be careful by using this, we do not check the input ptr and memory's accessment 
func GetAtIndex[T any](ptr unsafe.Pointer, index int) (out T) {
	return *(*T)(unsafe.Pointer(uintptr(ptr) + uintptr(index) * unsafe.Sizeof(out)))
}
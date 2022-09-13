package ptr

import (
	"testing"
	"unsafe"
)

func TestToBytes(t *testing.T) {
	a := 0x010203ff

	SetAtOffset[byte](unsafe.Pointer(&a), 0, 'a')
	SetAtOffset[byte](unsafe.Pointer(&a), 1, 'b')
	SetAtOffset[byte](unsafe.Pointer(&a), 2, 'c')
	SetAtOffset[byte](unsafe.Pointer(&a), 3, 'd')
	SetAtOffset[byte](unsafe.Pointer(&a), 4, 'a')
	SetAtOffset[byte](unsafe.Pointer(&a), 5, 'b')
	SetAtOffset[byte](unsafe.Pointer(&a), 6, 'c')
	SetAtOffset[byte](unsafe.Pointer(&a), 7, 'd')

	data := ValToBytes(a)
	data2 := []byte("abcdabcd")

	t.Logf("1: %s\n", data)
	t.Logf("2: %s\n", data2)
}
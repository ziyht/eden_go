package ptr

import (
	"testing"
	"unsafe"
)

func TestGetSet(t *testing.T) {
	a := 0x010203ff

	t.Logf("1: %d", GetAtOffset[byte](unsafe.Pointer(&a), 0))
	t.Logf("2: %d", GetAtOffset[byte](unsafe.Pointer(&a), 1))
	t.Logf("3: %d", GetAtOffset[byte](unsafe.Pointer(&a), 2))
	t.Logf("4: %d", GetAtOffset[byte](unsafe.Pointer(&a), 3))
	
	SetAtOffset[byte](unsafe.Pointer(&a), 0, 4)
	SetAtOffset[byte](unsafe.Pointer(&a), 1, 3)
	SetAtOffset[byte](unsafe.Pointer(&a), 2, 2)
	SetAtOffset[byte](unsafe.Pointer(&a), 3, 1)

	t.Logf("1: %d", GetAtOffset[byte](unsafe.Pointer(&a), 0))
	t.Logf("2: %d", GetAtOffset[byte](unsafe.Pointer(&a), 1))
	t.Logf("3: %d", GetAtOffset[byte](unsafe.Pointer(&a), 2))
	t.Logf("4: %d", GetAtOffset[byte](unsafe.Pointer(&a), 3))
}
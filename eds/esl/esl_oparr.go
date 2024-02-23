package esl

import (
	"sync/atomic"
	"unsafe"
)

const (
	op1 = 4
	op2 = maxLevel - op1
)

type opArray struct {
	base  [op1]unsafe.Pointer
	extra *([op2]unsafe.Pointer)
}

var __initVal = &opArray{}
func (a *opArray) init(i int) {
	a.base = __initVal.base
	if i > op1 {
		a.extra = new([op2]unsafe.Pointer)
	}
}

func (a *opArray) load(i int) unsafe.Pointer {
	if i < op1 {
		return a.base[i]
	}
	return a.extra[i-op1]
}

func (a *opArray) store(i int, p unsafe.Pointer) {
	if i < op1 {
		a.base[i] = p
		return
	}
	a.extra[i-op1] = p
}

func (a *opArray) atomicLoad(i int) unsafe.Pointer {
	if i < op1 {
		return atomic.LoadPointer(&a.base[i])
	}
	return atomic.LoadPointer(&a.extra[i-op1])
}

func (a *opArray) atomicStore(i int, p unsafe.Pointer) {
	if i < op1 {
		atomic.StorePointer(&a.base[i], p)
		return
	}
	atomic.StorePointer(&a.extra[i-op1], p)
}
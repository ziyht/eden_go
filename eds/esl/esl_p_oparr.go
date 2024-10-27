package esl

import (
	"sync/atomic"
	"unsafe"
)

const (
	op1 = 4
	op2 = _MAX_LEVEL - op1
)

type opArray struct {
	base  [op1]unsafe.Pointer
	extra *([op2]unsafe.Pointer)
}

func (a *opArray) load0() unsafe.Pointer {
	return a.base[0]
}

func (a *opArray) atomicLoad0() unsafe.Pointer {
	return atomic.LoadPointer(&a.base[0])
}

func (a *opArray) load(layer int) unsafe.Pointer {
	if layer < op1 {
		return a.base[layer]
	}
	return a.extra[layer-op1]
}

func (a *opArray) store(layer int, p unsafe.Pointer) {
	if layer < op1 {
		a.base[layer] = p
		return
	}
	a.extra[layer-op1] = p
}

func (a *opArray) atomicLoad(layer int) unsafe.Pointer {
	if layer < op1 {
		return atomic.LoadPointer(&a.base[layer])
	}
	return atomic.LoadPointer(&a.extra[layer-op1])
}

func (a *opArray) atomicStore(layer int, p unsafe.Pointer) {
	if layer < op1 {
		atomic.StorePointer(&a.base[layer], p)
		return
	}
	atomic.StorePointer(&a.extra[layer-op1], p)
}


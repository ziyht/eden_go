package esl

import (
	"sync/atomic"
	"unsafe"
)

type opArray struct {
	base  [_MAX_LEVEL]unsafe.Pointer				// _MAX_LEVEL is only for convenience
}

func (a *opArray) load0() unsafe.Pointer {
	return a.base[0]
}

func (a *opArray) atomicLoad0() unsafe.Pointer {
	return atomic.LoadPointer(&a.base[0])
}

func (a *opArray) load(layer int) unsafe.Pointer {
	return a.base[layer]
}

func (a *opArray) store(layer int, p unsafe.Pointer) {
	a.base[layer] = p
}

func (a *opArray) atomicLoad(layer int) unsafe.Pointer {
	return atomic.LoadPointer(&a.base[layer])
}

func (a *opArray) atomicStore(layer int, p unsafe.Pointer) {
	atomic.StorePointer(&a.base[layer], p)
}

package esl

import "sync/atomic"

const (
	fullyLinked = 1 << iota
	deleting
)

type bitFlag struct {
	data uint32
}

func (f *bitFlag) SetTrue(flags uint32) {
	for {
		old := atomic.LoadUint32(&f.data)
		if old&flags != flags {
			// Flag is 0, need set it to 1.
			n := old | flags
			if atomic.CompareAndSwapUint32(&f.data, old, n) {
				return
			}
			continue
		}
		return
	}
}

func (f *bitFlag) SetFalse(flags uint32) {
	for {
		old := atomic.LoadUint32(&f.data)
		check := old & flags
		if check != 0 {
			// Flag is 1, need set it to 0.
			n := old ^ check
			if atomic.CompareAndSwapUint32(&f.data, old, n) {
				return
			}
			continue
		}
		return
	}
}

func (f *bitFlag) Get(flag uint32) bool {
	return (atomic.LoadUint32(&f.data) & flag) != 0
}

func (f *bitFlag) Check(check, expect uint32) bool {
	return (atomic.LoadUint32(&f.data) & check) == expect
}

package esync

import (
	"testing"
	"sync"
)

type iRWMutex interface {
	Lock()
	TryLock() bool
	Unlock()
	RLock()
	TryRLock() bool
	RUnlock()
}

type nonMutex struct {
}

func (nmu *nonMutex)Lock(){}
func (nmu *nonMutex)TryLock()bool{return true}
func (nmu *nonMutex)Unlock(){}
func (nmu *nonMutex)RLock(){}
func (nmu *nonMutex)TryRLock()bool{return true}
func (nmu *nonMutex)RUnlock(){}

func newRWMutex(safe ...bool) iRWMutex {
	if len(safe) == 0 || len(safe) > 0 && safe[0] {
		return &sync.RWMutex{}
	}

	return &nonMutex{}
}

// Benchmark for performance
func BenchmarkRWMutex(b *testing.B) {
	b.Run("SafeLock", func(b *testing.B) {
		lock := newRWMutex(true)
		benchmarkLock(b, lock)
	})

	b.Run("UnsafeLock", func(b *testing.B) {
		lock := newRWMutex(false)
		benchmarkLock(b, lock)
	})
}

func benchmarkLock(b *testing.B, lock iRWMutex) {
	for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Unlock()
	}
}

func BenchmarkInjectedRWMutex(b *testing.B) {
	b.Run("SafeLock", func(b *testing.B) {
		lock := NewERWMutex(true)
		benchmarkLock2(b, lock)
	})

	b.Run("UnsafeLock", func(b *testing.B) {
		lock := NewERWMutex(false)
		benchmarkLock2(b, lock)
	})
}

func benchmarkLock2(b *testing.B, lock *ERWMutex) {
	for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Unlock()
	}
}
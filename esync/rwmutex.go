package esync

import (
	"sync"
)

type RWMutex interface {
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
func (nmu *nonMutex)TryLock()bool{return false}
func (nmu *nonMutex)Unlock(){}
func (nmu *nonMutex)RLock(){}
func (nmu *nonMutex)TryRLock()bool{return false}
func (nmu *nonMutex)RUnlock(){}

func NewRWMutex(safe ...bool) RWMutex {
	if len(safe) > 0 && safe[0] {
		return &sync.RWMutex{}
	}

	return &nonMutex{}
}

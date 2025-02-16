package esync

import (
	"sync"
)

type ERWMutex struct {
	sync.RWMutex
}

func NewERWMutex(safe ...bool) (mu *ERWMutex) {
	if len(safe) == 0 || safe[0] {
		return &ERWMutex{}
	}

	return nil
}

func (mu *ERWMutex) Lock()    { if mu != nil { mu.RWMutex.Lock()    } }
func (mu *ERWMutex) Unlock()  { if mu != nil { mu.RWMutex.Unlock()  } }
func (mu *ERWMutex) RLock()   { if mu != nil { mu.RWMutex.RLock()   } }
func (mu *ERWMutex) RUnlock() { if mu != nil { mu.RWMutex.RUnlock() } }

func (mu *ERWMutex) TryLock()  bool { if mu != nil { return mu.RWMutex.TryLock()  }; return true }
func (mu *ERWMutex) TryRLock() bool { if mu != nil { return mu.RWMutex.TryRLock() }; return true }

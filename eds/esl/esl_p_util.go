package esl

import (
	"github.com/ziyht/eden_go/erand"
)

func randomLevel() int {
	level := 1
	for erand.Uint32n(1/p) == 0 {
		level++
	}
	if level > maxLevel {
		return maxLevel
	}
	return level
}

// ordered is a constraint that permits any ordered type: any type
// that supports the operators < <= >= >.
type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | // sign
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | // unsign
		~float32 | ~float64 | // float
		~string
}
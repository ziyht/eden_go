package esl

import (
	"github.com/ziyht/eden_go/erand"
)

func randomLevel() int {
	level := 1
	for erand.Uint32n(1/P) == 0 {
		level++
	}
	if level > MAX_LEVEL {
		return MAX_LEVEL
	}
	return level
}
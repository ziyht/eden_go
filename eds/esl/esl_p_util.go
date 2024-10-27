package esl

import (
	"github.com/ziyht/eden_go/erand"
)

func randomLevel() int {
	level := 1
	for erand.Uint32n(1/_P) == 0 {
		level++
	}
	if level > _MAX_LEVEL {
		return _MAX_LEVEL
	}
	return level
}
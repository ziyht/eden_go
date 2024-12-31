package elem

import (
	cst "golang.org/x/exp/constraints"
)

type PairElem[K cst.Ordered, V any] struct {
  key K
	Val V
}

func (e *PairElem[K, V]) Key() K {
	return e.key
}

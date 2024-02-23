package esl

import (
	"unsafe"
)

// ESLNode represents one actual Node in the skiplist structure.
// It saves the actual element, pointers to the next nodes and a pointer to one previous node.
type ESLNode [K ordered, V any] struct {
	nexts opArray
	prev  *ESLNode[K, V]
	level int
	key   K
	value V
}

func genESLNode[K ordered, V any]() *ESLNode[K, V]{
	return &ESLNode[K, V]{}
}

func newESLNode[K ordered, V any](level int, key K, val V) *ESLNode[K, V]{
	node := &ESLNode[K, V]{
		level: level,
		key:   key,
		value: val,
	}
	node.nexts.init(level)

	return node
}

// GetValue extracts the ListElement value from a skiplist node.
func (e *ESLNode[K, V]) GetValue() V {
	return e.value
}

func (e *ESLNode[K, V]) Next() *ESLNode[K, V] {
	return e.atomicLoadNext(0)
}

func (e *ESLNode[K, V]) Prev() *ESLNode[K, V] {
	return e.prev
}

func (n *ESLNode[K, V]) loadNext(level int) *ESLNode[K, V] {
	return (*ESLNode[K, V])(n.nexts.load(level))
}

func (n *ESLNode[K, V]) storeNext(level int, next *ESLNode[K, V]) {
	n.nexts.store(level, unsafe.Pointer(next))
}

func (n *ESLNode[K, V]) atomicLoadNext(level int) *ESLNode[K, V] {
	return (*ESLNode[K, V])(n.nexts.atomicLoad(level))
}

func (n *ESLNode[K, V]) atomicStoreNext(level int, next *ESLNode[K, V]) {
	n.nexts.atomicStore(level, unsafe.Pointer(next))
}


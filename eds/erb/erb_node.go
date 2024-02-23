package erb

import (
	. "golang.org/x/exp/constraints"
)

type Node[K Ordered, V any] struct {
	color   int8
	parent 	*Node[K, V]
	left    *Node[K, V]
	right  	*Node[K, V]

	key     K
	Val     V
}

// Key returns the node's key.
func (n *Node[K, V]) Key() K {
	return n.key
}

// Next returns the node's successor.
func (n *Node[K, V]) Next() *Node[K, V] {
	return successor(n)
}

// Prev returns the node's predecessor.
func (n *Node[K, V]) Prev() *Node[K, V] {
	return predecessor(n)
}
package erb

import (
	cst "golang.org/x/exp/constraints"
)

// color of node
const (
	RED   = 0
	BLACK = 1
	VER   = "1.0.0"
)

// ERB is a struct of red-black tree.
type ERB[K cst.Ordered, V any] struct {
	root *Node[K, V]
	size int64
}

// New creates a new rbtree.
func New[K cst.Ordered, V any]() *ERB[K, V] {
	return &ERB[K, V]{}
}

// Add inserts the key-value pair into the rbtree. If the key already exists, will ignore and return false.
func (t *ERB[K, V])Add(key K, val V) bool {
	link := &t.root
	var parent *Node[K, V]

	for *link  != nil {
		parent = *link
		if key == parent.key {
			return false
		} else if key < parent.key { link = &parent.left
		} else                     { link = &parent.right
		}
	}

	node := &Node[K, V]{parent: parent, color: RED, key: key, Val: val}

	*link = node
	t.insert(node)

	return true
}

// Set sets the key-value pair into the rbtree
// If the key already exists, will replace the value and return the previous value.
func (t *ERB[K, V])Set(key K, val V)(prev V, replaced bool) {
	link := &t.root
	var parent *Node[K, V]

	for *link  != nil {
		parent = *link
		if key == parent.key {
			prev = parent.Val
			parent.Val = val
			return prev, true
		} else if key < parent.key { link = &parent.left
		} else                     { link = &parent.right
		}
	}

	node := &Node[K, V]{parent: parent, color: RED, key: key, Val: val}

	*link = node
	t.insert(node)

	return
}

// Get finds the node and return its value if found.
func (t *ERB[K, V]) Get(key K, df... V) (V, bool) {
	n := t.find(key)
	if n != nil {
		return n.Val, true
	}

	if len(df) > 0 {
		return df[0], false
	}

	var result V
	return result, false
}

// Val finds the node and return its value.
func (t *ERB[K, V]) Val(key K, df... V) V {
	n := t.find(key)
	if n != nil {
		return n.Val
	}

	if len(df) > 0 {
		return df[0]
	}

	var result V
	return result
}

// Find finds the node and return it
func (t *ERB[K, V]) Find(key K) (*Node[K, V]) {
	return t.find(key)
}

// First return the first node.
func (t *ERB[K, V]) First() *Node[K, V] {
	if t.root == nil {
		return nil
	}
	return t.root.minimum()
}

// First return the last node.
func (t *ERB[K, V]) Last() *Node[K, V] {
	if t.root == nil {
		return nil
	}
	return t.root.maximum()
}

// IsEmpty checks whether the rbtree is empty.
func (t *ERB[K, V]) IsEmpty() bool {
	return t.root == nil
}

// Size returns the size of the rbtree.
func (t *ERB[K, V]) Size() int64 {
	return t.size
}

// Clear unlink all node from rbtree.
func (t *ERB[K, V]) Clear() {
	t.root = nil
	t.size = 0
}

// Del deletes one node by key
func (t *ERB[K, V]) Del(key K) (bool) {
	n := t.find(key)
	t.erase(n)
	return n != nil
}

func (t *ERB[K, V]) PopFirst() *Node[K, V] {
	n := t.First(); t.erase(n); return n
}

func (t *ERB[K, V]) PopLast() *Node[K, V] {
	n := t.Last(); t.erase(n); return n
}

func Version() string {
	return VER
}

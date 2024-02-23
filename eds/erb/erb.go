package erb

import (
	. "golang.org/x/exp/constraints"
)

// color of node
const (
	RED   = 0
	BLACK = 1
)

// ERb is a struct of red-black tree.
type ERb[K Ordered, V any] struct {
	root *Node[K, V]
	size int64
}

// New creates a new rbtree.
func New[K Ordered, V any]() *ERb[K, V] {
	return &ERb[K, V]{}
}

// Add inserts the key-value pair into the rbtree. 
// multi defines if you can add the same key in tree or not, default is false
func (t *ERb[K, V]) Add(key K, value V, multi... bool) bool {
	var pos __pos_s[K, V]

	t.searchKeyPos(key, &pos, len(multi) > 0 && multi[0])

	if pos.canadd {
		node := &Node[K, V]{parent: pos.parent, color: RED, key: key, Val: value}
		if node.parent == nil {
			node.color = BLACK
		}

		*pos.instPos = node
		t.rbInsertFixup(node)

		t.size++
		return true
	}

	return false
}

// Get finds the node and return its value.
func (t *ERb[K, V]) Get(key K) (V, bool) {
	n := t.findnode(key)
	if n != nil {
		return n.Val, true
	}
	var result V
	return result, false
}

// Get finds the first node of key and return its value.
func (t *ERb[K, V]) GetFirst(key K) (V, bool) {
	n := t.findFirst(key)
	if n == nil {
		var result V
		return result, false
	}

	return n.Val, true
}

// Val finds the node and return its value.
func (t *ERb[K, V]) Val(key K) V {
	n := t.findnode(key)
	if n == nil {
		return n.Val
	}
	var result V
	return result
}

// Val finds the node and return it
func (t *ERb[K, V]) Node(key K) (*Node[K, V]) {
	return t.findnode(key)
}

// FindEX finds the node and return its value.
func (t *ERb[K, V]) FindEX(key K, cmp func (*Node[K, V]) bool) V {
	n := t.findnode(key)
	if n != nil && cmp(n){
		return n.Val
	}

	// find from prevs
	iter := n
	for {
		iter = predecessorSameKey(iter)
		if iter == nil {
			break
		}
		if cmp(iter){
			return iter.Val
		}
	}

	// find from nexts
	iter = n
	for {
		iter = successorSameKey(iter)
		if iter == nil {
			break
		}
		if cmp(iter){
			return iter.Val
		}
	}

	var result V
	return result
}

// FindNode finds the node and return it as an iterator.
func (t *ERb[K, V]) FindNode(key K) *Node[K, V] {
	return t.findnode(key)
}

// Empty checks whether the rbtree is empty.
func (t *ERb[K, V]) IsEmpty() bool {
	return t.root == nil
}

// First return the minmum node.
func (t *ERb[K, V]) First() *Node[K, V] {
	return minimum(t.root)
}

// First return the maxmum node.
func (t *ERb[K, V]) Last() *Node[K, V] {
	return maximum(t.root)
}

// Size returns the size of the rbtree.
func (t *ERb[K, V]) Size() int64 {
	return t.size
}

// Clear unlink all node from rbtree.
func (t *ERb[K, V]) Clear() {
	t.root = nil
	t.size = 0
}

func (t *ERb[K, V]) Set(key K, value V) {
	var pos __pos_s[K, V]

	t.searchKeyPos(key, &pos, false)
	if pos.find {
		pos.parent.Val = value
	} else {
		node := &Node[K, V]{parent: pos.parent, color: RED, key: key, Val: value}
		if node.parent == nil {
			node.color = BLACK
		}

		*pos.instPos = node
		t.rbInsertFixup(node)

		t.size++
	}
}

func (t *ERb[K, V])SetMulti(key K, value V) int64 {
	return 0
}

// Del deletes the node by key
func (t *ERb[K, V]) Del(key K) (bool) {
	z := t.findnode(key)
	t.deleteNode(z)
	return z != nil
}

func (t *ERb[K, V]) Pop(key K) *Node[K, V] {
	n := t.findnode(key); t.deleteNode(n); return n
}

func (t *ERb[K, V]) PopKeyFirst(key K) *Node[K, V] {
	n := t.findFirst(key); 
	t.deleteNode(n); 
	return n
}

func (t *ERb[K, V]) PopFirst() *Node[K, V] {
	n := t.First(); t.deleteNode(n); return n
}

func (t *ERb[K, V]) PopLast() *Node[K, V] {
	n := t.Last(); t.deleteNode(n); return n
}



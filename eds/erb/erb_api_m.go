package erb

import (
	cst "golang.org/x/exp/constraints"
)

type ERBM[K cst.Ordered, V any] ERB[K, V]

func (tm *ERB[K, V]) MulAdd(key K, val V) {
	t := (*ERB[K, V])(tm)
	link := &t.root
	var parent *Node[K, V]

	for *link != nil {
		parent = *link
		if     key < parent.key { link = &parent.left
		} else                  { link = &parent.right
		}
	}

	node := &Node[K, V]{parent: parent, color: RED, key: key, Val: val}
	if parent == nil {
		node.color = BLACK
	}
	*link = node
	t.insert(node)
}

// Get finds the first node of key and return its value.
func (t *ERB[K, V]) GetFirst(key K) (V, bool) {
	n := t.findFirst(key)
	if n == nil {
		var result V
		return result, false
	}

	return n.Val, true
}

func (t *ERB[K, V]) DelFirst(key K) *Node[K, V] {
	n := t.findFirst(key); 
	(*ERB[K, V])(t).erase(n); 
	return n
}

func (t *ERB[K, V]) DelLast(key K) *Node[K, V] {
	n := t.findLast(key); 
	t.erase(n); 
	return n
}

// Del deletes the nodes by key
func (t *ERB[K, V]) DelAll(key K) (int64) {
	var cnt int64

	for {
		n := t.find(key)
		if n == nil {
			break
		}

		t.erase(n)
		cnt++
	}

	return cnt
}
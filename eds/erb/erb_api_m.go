package erb

import (
	cst "golang.org/x/exp/constraints"
)

type ERBM[K cst.Ordered, V any] ERB[K, V]

func (t *ERB[K, V]) MulAdd(key K, val V) {
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

func (t *ERB[K, V]) SetAll(key K, val V) (cnt int64) {
	t.root.traverseNodeFromToInOrder(key, key, func (n *Node[K, V]) bool	{
		n.Val = val
		cnt++
		return true
	})

	return
}

func (t *ERB[K, V]) ValFirst(key K, df ...V) (ret V) {
	n := t.findFirst(key)
	if n == nil {
		if len(df) > 0 { return df[0] }
		return 
	}

	return n.Val
}

func (t *ERB[K, V]) ValLast(key K, df ...V) (ret V) {
	n := t.findLast(key)
	if n == nil {
		if len(df) > 0 { return df[0] }
		return 
	}

	return n.Val
}

// GetFirst finds the first node of specific key and return it.
func (t *ERB[K, V]) GetFirst(key K) (ret V, found bool) {
	n := t.findFirst(key)
	if n == nil {
		return 
	}

	return n.Val, true
}

func (t *ERB[K, V]) GetLast(key K) (ret V, found bool) {
	n := t.findLast(key)
	if n == nil {
		return 
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

	n := t.find(key)
	if n == nil {
		return cnt
	}

	for {
		b := n.prevBrotherAny()
		if b == nil { break }
		t.erase(b)
		cnt++
	}

	for {
		b := n.nextBrotherAny()
		if b == nil { break }
		t.erase(b)
		cnt++
	}

	t.erase(n)
	cnt++
	
	return cnt
}
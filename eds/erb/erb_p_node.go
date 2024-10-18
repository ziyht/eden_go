package erb

func (n *Node[K, V])setParentColor(parent *Node[K, V], color int8) {
	n.parent = parent
	n.color  = color
}

// predecessor returns the predecessor of the node
func (n *Node[K, V])predecessor() *Node[K, V] {
	/*
		* If we have a left-hand child, go down and then right as far
		* as we can.
		*/
	if n.left != nil {
		return n.left.maximum()
	}

	/*
		* No left-hand children. Go up till we find an ancestor which
		* is a right-hand child of its parent.
		*/
	parent := n.parent
	for parent != nil && n == parent.left {
		n = parent
		parent = n.parent
	}
	return parent
}

// successor returns the successor of the node
func (n *Node[K, V])successor() *Node[K, V] {
	/*
		* If we have a right-hand child, go down and then left as far
		* as we can.
		*/
	if n.right != nil {
		return n.right.minimum()
	}

	/*
		* No right-hand children. Everything down and left is smaller than us,
		* so any 'next' node must be in the general direction of our parent.
		* Go up the tree; any time the ancestor is a right-hand child of its
		* parent, keep going up. First time it's a left-hand child of its
		* parent, said parent is our 'next' node.
		*/
	parent := n.parent
	for parent != nil && n == parent.right {
		n = parent
		parent = n.parent
	}
	return parent
}

func (n *Node[K, V])prevBrother() *Node[K, V] {
	iter := n.predecessor()
	if iter == nil || iter.key != n.key{
		return nil
	}
	
	return iter
}

func (n *Node[K, V])nextBrother() *Node[K, V] {
	iter := n.successor()
	if iter == nil || iter.key != n.key{
		return nil
	}
	
	return iter
}

func (n *Node[K, V])prevBrotherAny() *Node[K, V] {
  if n.left != nil && n.left.key == n.key {
		return n.left
  }

	return n.prevBrother()
}

func (n *Node[K, V])nextBrotherAny() *Node[K, V] {
	if n.right != nil && n.right.key == n.key {
		return n.right
	}

	return n.nextBrother()
}

// minimum finds the minimum node of subtree n.
func (n *Node[K, V])minimum() *Node[K, V] {
	for n.left != nil {
		n = n.left
	}
	return n
}

// minimum finds the minimum node of subtree n.
func (n *Node[K, V])maximum() *Node[K, V] {
	for n.right != nil {
		n = n.right
	}
	return n
}

// firstBrother returns the first node with the same key of current node in pre-order
// if not have any more, return nil
func (n *Node[K, V])firstBrother() (fd *Node[K, V]) {
	key := n.key
	n = n.left

	for n != nil {
		if n.key == key {
			fd = n
			n = n.left
		} else {
			n = n.right
		}
	}

	return
}

// lastBrother returns the last node with the same key of current node in post-order
// if not have any more, return nil
func (n *Node[K, V])lastBrother() (fd *Node[K, V]) {
	key := n.key
	n = n.right

	for n != nil {
		if n.key == key {
			fd = n
			n = n.right
		} else {
			n = n.left
		}
	}

	return
}

func (n *Node[K, V]) traverseInOrder(callback func(k K, v V) bool) bool {
	if n != nil {
		if !n.left.traverseInOrder(callback) {
			return false
		}

		if !callback(n.key, n.Val) {
			return false
		}

		if !n.right.traverseInOrder(callback) {
			return false
		}
	}

	return true
}

func (n *Node[K, V]) traverseReverseInOrder(callback func(k K, v V) bool) bool {
	if n != nil {
		if !n.right.traverseReverseInOrder(callback) {
			return false
		}

		if !callback(n.key, n.Val) {
			return false
		}

		if !n.left.traverseReverseInOrder(callback) {
			return false
		}
	}

	return true
}

func (n *Node[K, V]) traverseNodeFromInOrder(from K, to K, callback func(n *Node[K, V]) bool) bool {
	if n != nil {
		if n.key >= from {
			if !n.left.traverseNodeFromInOrder(from, to, callback) {
				return false
			}
		}

		if n.key >= from && n.key <= to {
			if !callback(n) {
				return false
			}
		}

		if n.key <= to {
			if !n.right.traverseNodeFromInOrder(from, to, callback) {
				return false
			}
		}
	}

	return true
}

func (n *Node[K, V]) traverseNodeFromReverseInOrder(from K, to K, callback func(n *Node[K, V]) bool) bool {
	if n != nil {
		if n.key <= from {
			if !n.right.traverseNodeFromReverseInOrder(from, to, callback) {
				return false
			}
		}

		if n.key >= to && n.key <= from {
			if !callback(n) {
				return false
			}
		}

		if n.key >= to {
			if !n.left.traverseNodeFromReverseInOrder(from, to, callback) {
				return false
			}
		}
	}

	return true
}

func (n *Node[K, V]) traverseFromInOrder(from K, to K, callback func(k K, v V) bool) bool {
	return n.traverseNodeFromInOrder(from, to, func(n *Node[K, V]) bool {
		return callback(n.key, n.Val)
	})
}

func (n *Node[K, V]) traverseFromReverseInOrder(from K, to K, callback func(k K, v V) bool) bool {
	return n.traverseNodeFromReverseInOrder(from, to, func(n *Node[K, V]) bool {
		return callback(n.key, n.Val)
	})
}
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

func (n *Node[K, V])predecessorSameKey() *Node[K, V] {
	iter := n.predecessor()
	if iter == nil || iter.key != n.key{
		return nil
	}
	
	return iter
}

func (n *Node[K, V])successorSameKey() *Node[K, V] {
	iter := n.successor()
	if iter == nil || iter.key != n.key{
		return nil
	}
	
	return iter
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

func (n *Node[K, V])firstBrother2() (fd *Node[K, V]) {
	key := n.key
	tmp := n
	for tmp != nil {
		if       key < tmp.key { tmp = tmp.left
		}else if key > tmp.key { tmp = tmp.right
		}else                  { 
			fd = tmp;
			for tmp != nil {
				if     tmp.left != nil && key == tmp.left.key { fd = tmp.left; tmp = tmp.left
				} else                                        { fd = tmp     ; tmp = tmp.predecessorSameKey()
				}
			}
		}
	}
	return
}

func (n *Node[K, V])firstBrother() (fd *Node[K, V]) {
	key := n.key
	
	for n != nil {
		fd = n
		if n.left != nil && key == n.left.key {  
			n = n.left
		} else { 
			n = n.predecessor()
			if n == nil || key != n.key {
				return
			}
		}
	}

	return
}

func (n *Node[K, V])lastBrother() (fd *Node[K, V]) {
	tmp := n
	key := n.key

	for tmp != nil {
		if     tmp.right != nil && key == tmp.right.key { fd = tmp.right; tmp = tmp.right
		} else                                          { fd = tmp      ; tmp = tmp.successorSameKey()
		}
	}

	return
}

package erb

import (
	. "golang.org/x/exp/constraints"
)

type __pos_s [K Ordered, V any] struct {
	canadd   bool        
	find     bool
	parent 	 *Node[K, V]
	instPos  **Node[K, V]
};

// func (t *ERb[K, V]) __Add(key K, value V, multi... bool) bool {
// 	tmp := t.root
// 	var find *Node[K, V]

// 	for tmp != nil {
// 		find = tmp
// 		if key < tmp.key {
// 			tmp = tmp.left
// 		} else {
// 			tmp = tmp.right
// 		}
// 	}

// 	new := &Node[K, V]{parent: find, color: RED, key: key, Val: value}
	
// 	if find == nil {
// 		new.color = BLACK
// 		t.root = new
// 		t.size++
// 		return true
// 	} else if new.key < find.key {
// 		find.left = new
// 		t.size++
// 	} else {
// 		if (len(multi) == 0 || !multi[0]) && key == find.key{
// 			return false
// 		}
// 		find.right = new
// 		t.size++
// 	}
// 	t.rbInsertFixup(new)
	
// 	return true
// }

func (t *ERb[K, V]) searchKeyPos(k K, pos *__pos_s[K, V], multi bool){
  pos.instPos = &t.root

	if multi {
		for *pos.instPos != nil {
			pos.parent = *pos.instPos
			if k < pos.parent.key { pos.instPos = &pos.parent.left
			} else                { pos.instPos = &pos.parent.right
				pos.find = k == pos.parent.key
			}
		}
	} else {
		for *pos.instPos != nil {
			pos.parent = *pos.instPos
			if        k < pos.parent.key { pos.instPos = &pos.parent.left
			} else if k > pos.parent.key { pos.instPos = &pos.parent.right
			} else {
				pos.find = true
				return
			} 
		}
	}

	pos.canadd = true
}

func (t *ERb[K, V]) rbInsertFixup(z *Node[K, V]) {
	var y *Node[K, V]
	for z.parent != nil && z.parent.color == RED {
		if z.parent == z.parent.parent.left {
			y = z.parent.parent.right
			if y != nil && y.color == RED {
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					z = z.parent
					t.leftRotate(z)
				}
				z.parent.color = BLACK
				z.parent.parent.color = RED
				t.rightRotate(z.parent.parent)
			}
		} else {
			y = z.parent.parent.left
			if y != nil && y.color == RED {
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					t.rightRotate(z)
				}
				z.parent.color = BLACK
				z.parent.parent.color = RED
				t.leftRotate(z.parent.parent)
			}
		}
	}
	t.root.color = BLACK
}

func (t *ERb[K, V]) rbDeleteFixup(x, parent *Node[K, V]) {
	var w *Node[K, V]

	for x != t.root && getColor(x) == BLACK {
		if x != nil {
			parent = x.parent
		}
		if x == parent.left {
			w = parent.right
			if w.color == RED {
				w.color = BLACK
				parent.color = RED
				t.leftRotate(parent)
				w = parent.right
			}
			if getColor(w.left) == BLACK && getColor(w.right) == BLACK {
				w.color = RED
				x = parent
			} else {
				if getColor(w.right) == BLACK {
					if w.left != nil {
						w.left.color = BLACK
					}
					w.color = RED
					t.rightRotate(w)
					w = parent.right
				}
				w.color = parent.color
				parent.color = BLACK
				if w.right != nil {
					w.right.color = BLACK
				}
				t.leftRotate(parent)
				x = t.root
			}
		} else {
			w = parent.left
			if w.color == RED {
				w.color = BLACK
				parent.color = RED
				t.rightRotate(parent)
				w = parent.left
			}
			if getColor(w.left) == BLACK && getColor(w.right) == BLACK {
				w.color = RED
				x = parent
			} else {
				if getColor(w.left) == BLACK {
					if w.right != nil {
						w.right.color = BLACK
					}
					w.color = RED
					t.leftRotate(w)
					w = parent.left
				}
				w.color = parent.color
				parent.color = BLACK
				if w.left != nil {
					w.left.color = BLACK
				}
				t.rightRotate(parent)
				x = t.root
			}
		}
	}
	if x != nil {
		x.color = BLACK
	}
}

func (t *ERb[K, V]) leftRotate(x *Node[K, V]) {
	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func (t *ERb[K, V]) rightRotate(x *Node[K, V]) {
	y := x.left
	x.left = y.right
	if y.right != nil {
		y.right.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		t.root = y
	} else if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}
	y.right = x
	x.parent = y
}

// findnode finds the node by key and return it, if not exists return nil.
func (t *ERb[K, V]) findnode(key K) *Node[K, V] {
	tmp := t.root
	for tmp != nil {
		if        key  < tmp.key { tmp = tmp.left
		} else if key == tmp.key { return tmp
		} else                   { tmp = tmp.right
		}
	}
	return nil
}

// findnode finds the node by key and return it, if not exists return nil.
func (t *ERb[K, V]) findFirst(key K) (fd *Node[K, V]) {
	tmp := t.root
	for tmp != nil {
		if        key  < tmp.key { tmp = tmp.left
		} else if key == tmp.key { 
			for tmp != nil {
				fd = tmp
				tmp = predecessorSameKey(tmp)
			}
			return
		} else                   { tmp = tmp.right
		}
	}
	return
}

// findnode finds the node by key and return it, if not exists return nil.
func (t *ERb[K, V]) findFirst2(key K) (fd *Node[K, V]) {
	tmp := t.root
	for tmp != nil {
		if        key  < tmp.key { tmp = tmp.left
		} else if key == tmp.key { fd = tmp; tmp = tmp.left
		} else                   { tmp = tmp.right
		}
	}
	return
}


func (t *ERb[K, V])__rb_change_child(old, _new, parent *Node[K, V]){
    if (parent != nil) {
        if (parent.left == old) {
					parent.left = _new
				} else {
					parent.right = _new
				}
    } else {
			t.root = _new
		}
}

func (t *ERb[K, V]) deleteNode2(node *Node[K, V]) *Node[K, V] {
    child := node.right
		tmp   := node.left
		var parent, rebalance *Node[K, V]

    //var pc uint64

    if (tmp == nil) {
        /*
         * Case 1: node to erase has no more than 1 child (easy!)
         *
         * Note that if there is one child it must be red due to 5)
         * and node must be black due to 4). We adjust colors locally
         * so as to bypass __rb_erase_color() later on.
         */
        //pc = node->__rb_parent_color;
				parent = node.parent
				t.__rb_change_child(node, child, parent)
        if (child != nil) {
						child.parent = node.parent
						child.color  = node.color
            //rebalance = NULL;
        } else {
					if node.color == BLACK {
						rebalance = parent
					}
				}
        tmp = parent;
    } else if (child == nil) {
        /* Still case 1, but this time the child is node->rb_left */
				tmp.parent = node.parent; tmp.color = node.color
				t.__rb_change_child(node,tmp, node.parent)
        //rebalance = NULL;
        tmp = parent;
    } else {
				successor := child; var child2 *Node[K, V]

        tmp = child.left;
        if (tmp == nil) {
            /*
             * Case 2: node's successor is its right child
             *
             *    (n)          (s)
             *    / \          / \
             *  (x) (s)  ->  (x) (c)
             *        \
             *        (c)
             */
            parent = successor;
            child2 = successor.right;

            // augment->copy(node, successor);
        } else {
            /*
             * Case 3: node's successor is leftmost under
             * node's right child subtree
             *
             *    (n)          (s)
             *    / \          / \
             *  (x) (y)  ->  (x) (y)
             *      /            /
             *    (p)          (p)
             *    /            /
             *  (s)          (c)
             *    \
             *    (c)
             */
						for {
							parent = successor
							successor = tmp
							tmp = tmp.left
							if tmp == nil {
								break
							}
						}
            child2 = successor.right
						parent.left =child2
						successor.right = child
						child.parent = successor

            // augment->copy(node, successor);
            // augment->propagate(parent, successor);
        }

        tmp = node.left
				successor.left = tmp
				tmp.parent = successor

				t.__rb_change_child(node, successor, node.parent)

        if (child2 != nil) {
						successor.parent = node.parent; successor.color = node.color
						child2.parent = parent
						child2.color  = BLACK
            rebalance = nil;
        } else {
						c := successor.color
						successor.parent = node.parent; successor.color = node.color
						if c == BLACK {
							rebalance = parent
						}
        }
        tmp = successor
    }

    // augment->propagate(tmp, NULL);
    return rebalance
}

// deleteNode finds the node by key and return it, if not exists return nil.
func (t *ERb[K, V]) deleteNode(del *Node[K, V]) {
	if del == nil {
		return
	}

	var x, y *Node[K, V]
	if del.left != nil && del.right != nil {
		y = successor(del)
	} else {
		y = del
	}

	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}

	xparent := y.parent
	if x != nil {
		x.parent = xparent
	}
	if y.parent == nil {
		t.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	if y != del {
		del.key = y.key
		del.Val = y.Val
	}

	if y.color == BLACK {
		t.rbDeleteFixup(x, xparent)
	}
	t.size--
}

func predecessorSameKey[K Ordered, V any](n *Node[K, V]) *Node[K, V] {
	iter := n

	iter = predecessor(iter)
	if iter == nil || iter.key != n.key{
		return nil
	}
	
	return iter
}

func successorSameKey[K Ordered, V any](n *Node[K, V]) *Node[K, V] {
	iter := n

	iter = successor(iter)
	if iter == nil || iter.key != n.key{
		return nil
	}
	
	return iter
}

// predecessor returns the predecessor of the node
func predecessor[K Ordered, V any](n *Node[K, V]) *Node[K, V] {
	/*
		* If we have a left-hand child, go down and then right as far
		* as we can.
		*/
	if n.left != nil {
		return maximum(n.left)
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
func successor[K Ordered, V any](n *Node[K, V]) *Node[K, V] {
	/*
		* If we have a right-hand child, go down and then left as far
		* as we can.
		*/
	if n.right != nil {
		return minimum(n.right)
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

// getColor gets color of the node.
func getColor[K Ordered, V any](n *Node[K, V]) int8 {
	if n == nil {
		return BLACK
	}
	return n.color
}

// minimum finds the minimum node of subtree n.
func minimum[K Ordered, V any](n *Node[K, V]) *Node[K, V] {
	for n.left != nil {
		n = n.left
	}
	return n
}

// minimum finds the minimum node of subtree n.
func maximum[K Ordered, V any](n *Node[K, V]) *Node[K, V] {
	for n.right != nil {
		n = n.right
	}
	return n
}
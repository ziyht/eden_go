package erb

import (
	"fmt"
)

func (t *ERB[K, V]) rotateLeft(n *Node[K, V]) {
	right := n.right
	n.right = right.left
	if right.left != nil {
		right.left.parent = n
	}
	right.parent = n.parent
	if n.parent == nil {
		t.root = right
	} else if n == n.parent.left {
		n.parent.left = right
	} else {
		n.parent.right = right
	}
	right.left = n
	n.parent = right
}

func (t *ERB[K, V]) rotateRight(n *Node[K, V]) {
	left := n.left
	n.left = left.right
	if left.right != nil {
		left.right.parent = n
	}
	left.parent = n.parent
	if n.parent == nil {
		t.root = left
	} else if n == n.parent.right {
		n.parent.right = left
	} else {
		n.parent.left = left
	}
	left.right = n
	n.parent = left
}

func (t *ERB[K, V])__rb_change_child(old, _new, parent *Node[K, V]){
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

func (r *ERB[K, V]) insert(n *Node[K, V]) {

	var parent, gParent *Node[K, V]

	for parent = n.parent; parent != nil && parent.color == RED; parent = n.parent {

		gParent = parent.parent
		if parent == gParent.left {

			uncle := gParent.right
			if uncle != nil && uncle.color == RED {
				uncle.color = BLACK
				parent.color = BLACK
				gParent.color = RED
				n = gParent
				continue
			}

			if parent.right == n {
				r.rotateLeft(parent)
				parent, n = n, parent
			}

			parent.color = BLACK
			gParent.color = RED
			r.rotateRight(gParent)
		} else {
			uncle := gParent.left
			if uncle != nil && uncle.color == RED {
				uncle.color = BLACK
				parent.color = BLACK
				gParent.color = RED
				n = gParent
				continue
			}

			if parent.left == n {
				r.rotateRight(parent)
				parent, n = n, parent
			}
			parent.color = BLACK
			gParent.color = RED
			r.rotateLeft(gParent)
		}
	}
	r.root.color = BLACK //黑根
	r.size++
}

// find finds the node by key and return it, if not exists return nil.
func (t *ERB[K, V]) find(key K) *Node[K, V] {
	tmp := t.root
	for tmp != nil {
		if        key  < tmp.key { tmp = tmp.left
		} else if key == tmp.key { return tmp
		} else                   { tmp = tmp.right
		}
	}
	return nil
}

func (t *ERB[K, V]) __rotate_set_parents(old, new *Node[K, V], color int8) {
  parent := old.parent
	new.setParentColor(old.parent, old.color)
	old.setParentColor(new, color)
	t.__rb_change_child(old, new, parent)
}

/*
 * Inline version for rb_erase() use - we want to be able to inline
 * and eliminate the dummy_rotate callback there
 */
func (t *ERB[K, V]) __erase_color(parent *Node[K, V]){
    var node, sibling, tmp1, tmp2 *Node[K, V];

    for ;; {
        /*
         * Loop invariants:
         * - node is black (or NULL on first iteration)
         * - node is not the root (parent is not NULL)
         * - All leaf paths going through parent and node have a
         *   black node count that is 1 lower than other leaf paths.
         */
        sibling = parent.right;
        if (node != sibling) {	/* node == parent->rb_left */
            if (sibling.color == RED) {
                /*
                 * Case 1 - left rotate at parent
                 *
                 *     P               S
                 *    / \             / \
                 *   N   s    -->    p   Sr
                 *      / \         / \
                 *     Sl  Sr      N   Sl
                 */
                tmp1 = sibling.left
								parent.right = tmp1
                sibling.left = parent
                
								tmp1.setParentColor(parent, BLACK)
								t.__rotate_set_parents(parent, sibling, RED)
                
                sibling = tmp1;
            }
            tmp1 = sibling.right;
            if (tmp1 == nil || tmp1.color == BLACK) {
                tmp2 = sibling.left;
                if (tmp2 == nil || tmp2.color == BLACK) {
                    /*
                     * Case 2 - sibling color flip
                     * (p could be either color here)
                     *
                     *    (p)           (p)
                     *    / \           / \
                     *   N   S    -->  N   s
                     *      / \           / \
                     *     Sl  Sr        Sl  Sr
                     *
                     * This leaves us violating 5) which
                     * can be fixed by flipping p to black
                     * if it was red, or by recursing at p.
                     * p is red when coming from Case 1.
                     */
										sibling.setParentColor(parent, RED)
                    if (parent.color == RED){
											parent.color = BLACK
										} else {
                        node = parent;
                        parent = node.parent
                        if (parent != nil) {
													continue
												}
                    }
                    break;
                }
                /*
                 * Case 3 - right rotate at sibling
                 * (p could be either color here)
                 *
                 *   (p)           (p)
                 *   / \           / \
                 *  N   S    -->  N   Sl
                 *     / \             \
                 *    sl  Sr            s
                 *                       \
                 *                        Sr
                 */
                tmp1         = tmp2.right
								sibling.left = tmp1
								tmp2.right   = sibling
								parent.right = tmp2
                if (tmp1 != nil) {
									tmp1.setParentColor(sibling, BLACK)
								}
                tmp1    = sibling;
                sibling = tmp2;
            }
            /*
             * Case 4 - left rotate at parent + color flips
             * (p and sl could be either color here.
             *  After rotation, p becomes black, s acquires
             *  p's color, and sl keeps its color)
             *
             *      (p)             (s)
             *      / \             / \
             *     N   S     -->   P   Sr
             *        / \         / \
             *      (sl) sr      N  (sl)
             */
            tmp2 = sibling.left;
						parent.right = tmp2
            sibling.left = parent
            tmp1.setParentColor(sibling, BLACK)
            if (tmp2 != nil){
							tmp2.parent = parent
						}
						t.__rotate_set_parents(parent, sibling, BLACK)
            break;
        } else {
            sibling = parent.left;
            if (sibling.color == RED) {
                /* Case 1 - right rotate at parent */
                tmp1 = sibling.right
								parent.left   = tmp1
                sibling.right = parent
                tmp1.setParentColor(parent, BLACK)
                t.__rotate_set_parents(parent, sibling, RED)
                sibling = tmp1
            }
            tmp1 = sibling.left
            if (tmp1 == nil || tmp1.color == BLACK) {
                tmp2 = sibling.right;
                if (tmp2 == nil || tmp2.color == BLACK) {
                    /* Case 2 - sibling color flip */
										sibling.setParentColor(parent, RED)
                    if (parent.color == RED){
											parent.color = BLACK
										} else {
                        node = parent
                        parent = node.parent
                        if (parent != nil) {
													continue
												}
                    }
                    break;
                }
                /* Case 3 - right rotate at sibling */
                tmp1 = tmp2.left
								sibling.right = tmp1
                tmp2.left = sibling
								parent.left = tmp2
                if (tmp1 != nil) {
									tmp1.setParentColor(sibling, BLACK)
								}
                tmp1    = sibling
                sibling = tmp2
            }
            /* Case 4 - left rotate at parent + color flips */
            tmp2 = sibling.right
						parent.left = tmp2
						sibling.right = parent
						tmp1.setParentColor(sibling, BLACK)
            if (tmp2 != nil) {
							tmp2.parent = parent
						}
						t.__rotate_set_parents(parent, sibling, BLACK)
            break;
        }
    }
}

func (t *ERB[K, V]) erase(node *Node[K, V]) *Node[K, V] {
	if node == nil {
		return nil
	}

	child := node.right
	tmp   := node.left
	var parent, rebalance *Node[K, V]

	if (tmp == nil) {
		/*
			* Case 1: node to erase has no more than 1 child (easy!)
			*
			* Note that if there is one child it must be red due to 5)
			* and node must be black due to 4). We adjust colors locally
			* so as to bypass __rb_erase_color() later on.
			*/
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
		//tmp = parent;
	} else if (child == nil) {
			/* Still case 1, but this time the child is node->rb_left */
			tmp.parent = node.parent; tmp.color = node.color
			t.__rb_change_child(node,tmp, node.parent)
			//rebalance = NULL;
			//tmp = parent;
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
			//tmp = successor
	}

  //augment->propagate(tmp, NULL);

	if (rebalance != nil) {
		t.__erase_color(rebalance)
	}

	t.size--

	return rebalance
}

func (t *ERB[K, V])isRBTree() (bool) {
	ok, _ := t.isRBTree2()
	return ok
}

func (t *ERB[K, V])isRBTree2() (bool, error) {
	// RBTree Properties:
	// 1. Each node is either red or black.
	// 2. The root is BLACK.
	// 3. All leaves (NIL) are black.
	// 4. If a node is red, both its children are black.
	// 5. Every path from a given node to any of its descendant NIL nodes contains the same number of black nodes.
	_, property, ok := t.test(t.root)
	if !ok {
		return false, fmt.Errorf("violate property %v", property)
	}
	return true, nil
}

func (t *ERB[K, V])test(n *Node[K, V]) (int, int, bool) {

	if n == nil { // property 3:
		return 1, 0, true
	}

	if n == t.root && n.color == RED { // property 2:
		return 1, 2, false
	}
	leftBlackCount, property, ok := t.test(n.left)
	if !ok {
		return leftBlackCount, property, ok
	}
	rightBlackCount, property, ok := t.test(n.right)
	if !ok {
		return rightBlackCount, property, ok
	}

	if rightBlackCount != leftBlackCount { // property 5:
		return leftBlackCount, 5, false
	}
	blackCount := leftBlackCount

	if n.color == RED {
		if !(n.left == nil || n.left.color == BLACK) || !(n.right == nil || n.right.color == BLACK) { // property 4:
			return 0, 4, false
		}
	} else {
		blackCount++
	}

	// if n == t.root {
	// 	fmt.Printf("blackCount:%v \n", blackCount)
	// }
	return blackCount, 0, true
}
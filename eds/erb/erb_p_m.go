package erb

// findFirst finds the node by key and return it, if not exists return nil.
func (t *ERB[K, V]) findFirst(key K) (fd *Node[K, V]) {
	tmp := t.root
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

// findFirst finds the node by key and return it, if not exists return nil.
func (t *ERB[K, V]) findLast(key K) (fd *Node[K, V]) {
	tmp := t.root
	for tmp != nil {
		if       key < tmp.key { tmp = tmp.left
		}else if key > tmp.key { tmp = tmp.right
		}else                  { 
			fd = tmp;
			for tmp != nil {
				if     tmp.right != nil && key == tmp.right.key { fd = tmp.right; tmp = tmp.right
				} else                                          { fd = tmp      ; tmp = tmp.successorSameKey()
				}
			}
		}
	}
	return
}

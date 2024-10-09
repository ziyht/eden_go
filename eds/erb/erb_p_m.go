package erb

// findFirst finds the node by key and return it, if not exists return nil.
func (t *ERB[K, V]) findFirst(key K) (fd *Node[K, V]) {
	tmp := t.root
	for tmp != nil {
		if       key < tmp.key { tmp = tmp.left
		}else if key > tmp.key { tmp = tmp.right
		}else                  {
			fd = tmp.firstBrother()
			if fd == nil { return tmp }
			return
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
			fd = tmp.lastBrother()
			if fd == nil { return tmp }
			return
		}
	}
	return
}

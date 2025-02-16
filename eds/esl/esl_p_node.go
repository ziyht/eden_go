package esl

import (
	"sync"
	"sync/atomic"
	"unsafe"

	cst "golang.org/x/exp/constraints"
)

// ESL represents a set based on skip list.

type meta struct {
	mu    sync.Mutex
	flags bitFlag
	level uint32
}

type elem[K cst.Ordered, V any] struct {
	key K
	Val V
}

type Node[K cst.Ordered, V any] struct {
	meta
	elem[K, V]
	next [MAX_LEVEL]unsafe.Pointer
}

func newEslNode[K cst.Ordered, V any](key K, val V, level int) *Node[K, V] {
	switch level{
		case  1: n := struct{ meta; elem[K, V]; next [ 1]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case  2: n := struct{ meta; elem[K, V]; next [ 2]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case  3: n := struct{ meta; elem[K, V]; next [ 3]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case  4: n := struct{ meta; elem[K, V]; next [ 4]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case  5: n := struct{ meta; elem[K, V]; next [ 5]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case  6: n := struct{ meta; elem[K, V]; next [ 6]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case  7: n := struct{ meta; elem[K, V]; next [ 7]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case  8: n := struct{ meta; elem[K, V]; next [ 8]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case  9: n := struct{ meta; elem[K, V]; next [ 9]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 10: n := struct{ meta; elem[K, V]; next [10]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 11: n := struct{ meta; elem[K, V]; next [11]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 12: n := struct{ meta; elem[K, V]; next [12]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 13: n := struct{ meta; elem[K, V]; next [13]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 14: n := struct{ meta; elem[K, V]; next [14]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 15: n := struct{ meta; elem[K, V]; next [15]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 16: n := struct{ meta; elem[K, V]; next [16]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 17: n := struct{ meta; elem[K, V]; next [17]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 18: n := struct{ meta; elem[K, V]; next [18]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 19: n := struct{ meta; elem[K, V]; next [19]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 20: n := struct{ meta; elem[K, V]; next [20]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 21: n := struct{ meta; elem[K, V]; next [21]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 22: n := struct{ meta; elem[K, V]; next [22]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 23: n := struct{ meta; elem[K, V]; next [23]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 24: n := struct{ meta; elem[K, V]; next [24]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 25: n := struct{ meta; elem[K, V]; next [25]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 26: n := struct{ meta; elem[K, V]; next [26]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 27: n := struct{ meta; elem[K, V]; next [27]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 28: n := struct{ meta; elem[K, V]; next [28]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 29: n := struct{ meta; elem[K, V]; next [29]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 30: n := struct{ meta; elem[K, V]; next [30]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 31: n := struct{ meta; elem[K, V]; next [31]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		case 32: n := struct{ meta; elem[K, V]; next [32]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
		default: n := struct{ meta; elem[K, V]; next [MAX_LEVEL]unsafe.Pointer }{ meta: meta{ level: uint32(level) }, elem: elem[K, V]{ key: key, Val: val } }; return (*Node[K, V])(unsafe.Pointer(&n))
	}
}

func (n *Node[K, V]) Next() *Node[K, V] {
	return (*Node[K, V])(n.next[0])
}

func (n *Node[K, V]) atomicNext() *Node[K, V] {
	return (*Node[K, V])(atomic.LoadPointer(&n.next[0]))
}

func (n *Node[K, V]) loadNext(layer int) *Node[K, V] {
	return (*Node[K, V])(n.next[layer])
}

func (n *Node[K, V]) atomicLoadNext(layer int) *Node[K, V] {
	return (*Node[K, V])(atomic.LoadPointer(&n.next[layer]))
}

func (n *Node[K, V]) storeNext(layer int, next *Node[K, V]) {
	n.next[layer] = unsafe.Pointer(next)
}

func (n *Node[K, V]) atomicStoreNext(layer int, next *Node[K, V]) {
	atomic.StorePointer(&n.next[layer], unsafe.Pointer(next))
}

// findNodeRemove takes a value and two maximal-height arrays then searches exactly as in a sequential skip-list.
// The returned preds and succs always satisfy preds[i] > value >= succs[i].
func (s *ESL[K, V]) findNodeRemove(key K, preds *[MAX_LEVEL]*Node[K, V], succs *[MAX_LEVEL]*Node[K, V]) int {
	// lFound represents the index of the first layer at which it found a node.
	lFound, x := -1, s.tail
	if x != nil && x.key < key {			// not exist
		return lFound
	}

	x = &s.header
	for i := int(atomic.LoadUint64(&s.level)) - 1; i >= 0; i-- {
		succ := x.atomicLoadNext(i)
		for succ != nil && (succ.key < key) {
			x = succ
			succ = x.atomicLoadNext(i)
		}
		preds[i] = x
		succs[i] = succ

		// Check if the key already in the skip list.
		if lFound == -1 && succ != nil && succ.key == key {
			lFound = i
		}
	}
	return lFound
}

// findNodeAdd takes a key and two maximal-height arrays then searches exactly as in a sequential skip-set.
// The returned preds and succs always satisfy preds[i] > value >= succs[i].
func (s *ESL[K, V]) findNodeAdd(key K, preds *[MAX_LEVEL]*Node[K, V], succs *[MAX_LEVEL]*Node[K, V]) int {
	x := &s.header
	for i := int(atomic.LoadUint64(&s.level)) - 1; i >= 0; i-- {
		succ := x.atomicLoadNext(i)
		for succ != nil && (succ.key < key) {
			x = succ
			succ = x.atomicLoadNext(i)
		}
		preds[i] = x
		succs[i] = succ

		// Check if the key already in the skip list.
		if succ != nil && succ.key == key {
			return i
		}
	}
	return -1
}

func unlockOrdered[K cst.Ordered, V any](preds [MAX_LEVEL]*Node[K, V], highestLevel int) {
	var prevPred *Node[K, V]
	for i := highestLevel; i >= 0; i-- {
		if preds[i] != prevPred { // the node could be unlocked by previous loop
			preds[i].mu.Unlock()
			prevPred = preds[i]
		}
	}
}

// Add adds the value into skip set, returns true if this process insert the value into skip set,
// returns false if this process can't insert this value, because another process has insert the same value.
//
// If the value is in the skip set but not fully linked, this process will wait until it is.
func (s *ESL[K, V]) add(key K, val V, set bool) (prev V, replaced bool) {
	level := s.randomLevel()
	var preds, succs [MAX_LEVEL]*Node[K, V]
	for {
		lFound := s.findNodeAdd(key, &preds, &succs)
		if lFound != -1 { // indicating the value is already in the skip-list
			nodeFound := succs[lFound]
			if !nodeFound.flags.Get(deleting) {
				for !nodeFound.flags.Get(fullyLinked) {
					// The node is not yet fully linked, just waits until it is.
				}

				if set {
					prev = nodeFound.elem.Val
					nodeFound.elem.Val = val
					replaced = true
					return 
				}

				// replaced now means insert failed
				return 
			}
			// If the node is marked, represents some other thread is in the process of deleting this node,
			// we need to add this node in next loop.
			continue
		}
		// Add this node into skip list.
		var (
			highestLocked        = -1 // the highest level being locked by this process
			valid                = true
			pred, succ, prevPred *Node[K, V]
		)
		for layer := 0; valid && layer < level; layer++ {
			pred = preds[layer]   // target node's previous node
			succ = succs[layer]   // target node's next node
			if pred != prevPred { // the node in this layer could be locked by previous loop
				pred.mu.Lock()
				highestLocked = layer
				prevPred = pred
			}
			// valid check if there is another node has inserted into the skip list in this layer during this process.
			// It is valid if:
			// 1. The previous node and next node both are not marked.
			// 2. The previous node's next node is succ in this layer.
			valid = !pred.flags.Get(deleting) && (succ == nil || !succ.flags.Get(deleting)) && pred.loadNext(layer) == succ
		}
		if !valid {
			unlockOrdered(preds, highestLocked)
			continue
		}

		nn := newEslNode(key, val, level)
		for layer := 0; layer < level; layer++ {
			nn.storeNext(layer, succs[layer])
			preds[layer].atomicStoreNext(layer, nn)
		}
		if s.tail == nil {
			s.tail = nn
		} else if nn.key >=s.tail.key {
			s.tail = nn
		}
		nn.flags.SetTrue(fullyLinked)
		unlockOrdered(preds, highestLocked)
		atomic.AddInt64(&s.length, 1)

		// replaced now means insert ok
		replaced = !set
		return 
	}
}

func (s *ESL[K, V]) randomLevel() int {
	// Generate random level.
	level := randomLevel()
	// Update highest level if possible.
	for {
		hl := atomic.LoadUint64(&s.level)
		if level <= int(hl) {
			break
		}

		// Can only increase 2 level at a time
		if level > int(hl + 2) {
			level = int(hl + 2)
		}

		if atomic.CompareAndSwapUint64(&s.level, hl, uint64(level)) {
			break
		}
	}
	return level
}

func (s *ESL[K, V]) find(key K) (*Node[K, V]) {
	x := s.tail
	if x != nil{
		if x.key == key { return x}
		if x.key <  key { return nil   }
	}

	x = &s.header
	for i := int(atomic.LoadUint64(&s.level)) - 1; i >= 0; i-- {
		nex := x.atomicLoadNext(i)
		for nex != nil && (nex.key < key) {
			x = nex
			nex = x.atomicLoadNext(i)
		}

		// Check if the value already in the skip list.
		if nex != nil && nex.key == key {
			if nex.flags.Check(fullyLinked|deleting, fullyLinked) {
				return nex
			}
			// deleting means node expect not in the skip list
			return nil
		}
	}

	return nil
}

// Remove removes a node from the skip set.
func (s *ESL[K, V]) Remove(key K) bool {
	var (
		nodeToRemove *Node[K, V]
		isDeleting   bool // represents if this operation mark the node
		topLayer     = -1
		preds, succs [MAX_LEVEL]*Node[K, V]
	)
	for {
		lFound := s.findNodeRemove(key, &preds, &succs)
		if isDeleting || // this process deleting this node or we can find this node in the skip list
				lFound != -1 && succs[lFound].flags.Check(fullyLinked|deleting, fullyLinked) && (int(succs[lFound].level)-1) == lFound {
			if !isDeleting { // we don't deleting this node for now
				nodeToRemove = succs[lFound]
				topLayer = lFound
				nodeToRemove.mu.Lock()
				if nodeToRemove.flags.Get(deleting) {
					// The node is deleting by another process,
					// the physical deletion will be accomplished by another process.
					nodeToRemove.mu.Unlock()
					return false
				}
				nodeToRemove.flags.SetTrue(deleting)
				isDeleting = true
			}
			// Accomplish the physical deletion.
			var (
				highestLocked        = -1 // the highest level being locked by this process
				valid                = true
				pred, succ, prevPred *Node[K, V]
			)
			for layer := 0; valid && (layer <= topLayer); layer++ {
				pred, succ = preds[layer], succs[layer]
				if pred != prevPred { // the node in this layer could be locked by previous loop
					pred.mu.Lock()
					highestLocked = layer
					prevPred = pred
				}
				// valid check if there is another node has inserted into the skip list in this layer
				// during this process, or the previous is removed by another process.
				// It is valid if:
				// 1. the previous node exists.
				// 2. no another node has inserted into the skip list in this layer.
				valid = !pred.flags.Get(deleting) && pred.loadNext(layer) == succ
			}
			if !valid {
				unlockOrdered(preds, highestLocked)
				continue
			}
			for i := topLayer; i >= 0; i-- {
				// Now we own the nodeToRemove, no other goroutine will modify it.
				// So we don't need nodeToRemove.loadNext
				preds[i].atomicStoreNext(i, nodeToRemove.loadNext(i))
			}
			if nodeToRemove.Next() == nil {
				s.tail = preds[0]
				if s.tail == &s.header {
					s.tail = nil
				}
			}

			// Down level if possible.
			for s.level > INIT_LEVEL && s.header.next[s.level-1] == nil {
				s.level--
			}
			nodeToRemove.mu.Unlock()
			unlockOrdered(preds, highestLocked)
			atomic.AddInt64(&s.length, -1)
			return true
		}
		return false
	}
}

// traverseNode calls f sequentially for each value present in the skip list.
// If f returns false, range stops the iteration.
func (s *ESL[K, V]) traverseNode(f func( *Node[K, V]) bool) {
	x := s.header.atomicLoadNext(0)
	for x != nil {
		if !x.flags.Check(fullyLinked|deleting, fullyLinked) {
			x = x.atomicLoadNext(0)
			continue
		}
		if !f(x) {
			break
		}
		x = x.atomicLoadNext(0)
	}
}

// traverse calls f sequentially for each value present in the skip list.
// If f returns false, range stops the iteration.
func (s *ESL[K, V]) traverse(f func(key K, val V) bool) {
	x := s.header.atomicLoadNext(0)
	for x != nil {
		if !x.flags.Check(fullyLinked|deleting, fullyLinked) {
			x = x.atomicLoadNext(0)
			continue
		}
		if !f(x.key, x.Val) {
			break
		}
		x = x.atomicLoadNext(0)
	}
}

// traverseFrom calls f sequentially for all elems with `key >= start` in the skip list.
// If f returns false, range stops the iteration.
func (s *ESL[K, V]) traverseFrom(start K, f func(key K, val V) bool) {
	var (
		x   = &s.header
		nex *Node[K, V]
	)
	for i := int(atomic.LoadUint64(&s.level)) - 1; i >= 0; i-- {
		nex = x.atomicLoadNext(i)
		for nex != nil && (nex.key < start) {
			x = nex
			nex = x.atomicLoadNext(i)
		}
		// Check if the value already in the skip list.
		if nex != nil && nex.key == start {
			break
		}
	}

	for nex != nil {
		if !nex.flags.Check(fullyLinked|deleting, fullyLinked) {
			nex = nex.atomicNext()
			continue
		}
		if !f(nex.key, nex.Val) {
			break
		}
		nex = nex.atomicNext()
	}
}

// traverseFrom calls f sequentially for all elems with `start <= key <= end` in the skip list.
// If f returns false, range stops the iteration.
func (s *ESL[K, V]) traverseFromTo(start K, end K, f func(key K, val V) bool) {
	var (
		x   = &s.header
		nex *Node[K, V]
	)
	for i := int(atomic.LoadUint64(&s.level)) - 1; i >= 0; i-- {
		nex = x.atomicLoadNext(i)
		for nex != nil && (nex.key < start) {
			x = nex
			nex = x.atomicLoadNext(i)
		}
		// Check if the value already in the skip list.
		if nex != nil && nex.key == start {
			break
		}
	}

	for nex != nil {
		if !nex.flags.Check(fullyLinked|deleting, fullyLinked) {
			nex = nex.atomicNext()
			continue
		}
		if nex.key > end {
			return
		}
		if !f(nex.key, nex.Val) {
			break
		}
		nex = nex.atomicNext()
	}
}

func (s *ESL[K, V]) traverseIn(start K, end K, f func(key K, val V) bool) {
	var (
		x   = &s.header
		nex *Node[K, V]
	)
	for i := int(atomic.LoadUint64(&s.level)) - 1; i >= 0; i-- {
		nex = x.atomicLoadNext(i)
		for nex != nil && (nex.key < start) {
			x = nex
			nex = x.atomicLoadNext(i)
		}
		// Check if the value already in the skip list.
		if nex != nil && nex.key == start {
			break
		}
	}

	for nex != nil {
		if !nex.flags.Check(fullyLinked|deleting, fullyLinked) {
			nex = nex.atomicNext()
			continue
		}
		if nex.key >= end {
			return
		}
		if !f(nex.key, nex.Val) {
			break
		}
		nex = nex.atomicNext()
	}
}
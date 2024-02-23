package esl

import (
	"math/bits"
	"math/rand"
)

const (
	// maxLevel denotes the maximum height of the skiplist. This height will keep the skiplist
	// efficient for up to 34m entries. If there is a need for much more, please adjust this constant accordingly.
	maxLevel = 25
	eps      = 0.00001
)

// // ListElement is the interface to implement for elements that are inserted into the skiplist.
// type ListElement interface {
// 	// ExtractKey() returns a float64 representation of the key that is used for insertion/deletion/find. It needs to establish an order over all elements
// 	ExtractKey() float64
// 	// A string representation of the element. Can be used for pretty-printing the list. Otherwise just return an empty string.
// 	String() string
// }

type ordered interface {
	~int  | ~int8  | ~int16  | ~int32  | ~int64  | // sign
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | // unsign
	~float32 | ~float64 | // float
	~string
}



// SkipList is the actual skiplist representation.
// It saves all nodes accessible from the start and end and keeps track of element count, eps and levels.
type SkipList [K ordered, V any] struct {
	head         [maxLevel]*ESLNode[K, V]
	tail         [maxLevel]*ESLNode[K, V]
	maxNewLevel  int
	maxLevel     int
	count        int
}

// New returns a new empty, initialized Skiplist.
func New[K ordered, V any]() *SkipList[K, V] {
	return &SkipList[K, V]{
		head:         [maxLevel]*ESLNode[K, V]{},
		tail:         [maxLevel]*ESLNode[K, V]{},
		maxNewLevel:  maxLevel,
		maxLevel:     0,
		count:        0,
	}
}

// IsEmpty checks, if the skiplist is empty.
func (t *SkipList[K, V]) IsEmpty() bool {
	return t.head[0] == nil
}

func (t *SkipList[K, V]) generateLevel(maxLevel int) int {
	level := maxLevel - 1
	// First we apply some mask which makes sure that we don't get a level
	// above our desired level. Then we find the first set bit.
	var x uint64 = rand.Uint64() & ((1 << uint(maxLevel-1)) - 1)
	zeroes := bits.TrailingZeros64(x)
	if zeroes <= maxLevel {
		level = zeroes
	}

	return level
}

func (t *SkipList[K, V]) findEntryIndex(key K, level int) int {
	// Find good entry point so we don't accidentally skip half the list.
	for i := t.maxLevel; i >= 0; i-- {
		if t.head[i] != nil && t.head[i].key <= key || i <= level {
			return i
		}
	}
	return 0
}

func (t *SkipList[K, V]) findExtended(key K, findGreaterOrEqual bool) (foundElem *ESLNode[K, V], ok bool) {

	foundElem = nil
	ok = false

	if t.IsEmpty() {
		return
	}

	index := t.findEntryIndex(key, 0)
	var currentNode *ESLNode[K, V]

	currentNode = t.head[index]
	nextNode := currentNode

	// In case, that our first element is already greater-or-equal!
	if findGreaterOrEqual && currentNode.key > key {
		foundElem = currentNode
		ok = true
		return
	}

	for {
		if currentNode.key == key {
			foundElem = currentNode
			ok = true
			return
		}

		nextNode = (*ESLNode[K, V])(currentNode.nexts.load(index))

		// Which direction are we continuing next time?
		if nextNode != nil && nextNode.key <= key {
			// Go right
			currentNode = nextNode
		} else {
			if index > 0 {

				next := currentNode.nexts.load(0)

				// Early exit
				if next != nil && currentNode.key == key {
					foundElem = (*ESLNode[K, V])(next)
					ok = true
					return
				}
				// Go down
				index--
			} else {
				// Element is not found and we reached the bottom.
				if findGreaterOrEqual {
					foundElem = nextNode
					ok = nextNode != nil
				}

				return
			}
		}
	}
}

// Find tries to find an element in the skiplist based on the key from the given ListElement.
// elem can be used, if ok is true.
// Find runs in approx. O(log(n))
func (t *SkipList[K, V]) Find(key K) (elem *ESLNode[K, V], ok bool) {

	if t == nil {
		return
	}

	elem, ok = t.findExtended(key, false)
	return
}

// FindGreaterOrEqual finds the first element, that is greater or equal to the given ListElement e.
// The comparison is done on the keys (So on ExtractKey()).
// FindGreaterOrEqual runs in approx. O(log(n))
func (t *SkipList[K, V]) FindGreaterOrEqual(key K) (elem *ESLNode[K, V], ok bool) {

	if t == nil {
		return
	}

	elem, ok = t.findExtended(key, true)
	return
}

// Delete removes an element equal to e from the skiplist, if there is one.
// If there are multiple entries with the same value, Delete will remove one of them
// (Which one will change based on the actual skiplist layout)
// Delete runs in approx. O(log(n))
func (t *SkipList[K, V]) Delete(key K) {

	if t == nil || t.IsEmpty() {
		return
	}

	index := t.findEntryIndex(key, 0)

	var currentNode *ESLNode[K, V]
	nextNode := currentNode

	for {

		if currentNode == nil {
			nextNode = t.head[index]
		} else {
			nextNode = (*ESLNode[K, V])(currentNode.nexts.load(index))
		}

		// Found and remove!
		if nextNode != nil && nextNode.key == key {

			if currentNode != nil {
				currentNode.atomicStoreNext(index, nextNode.loadNext(index))
			}

			next := nextNode.loadNext(index)
			if index == 0 {
				if next != nil {
					next.prev = currentNode
				}
				t.count--
			}

			// Link from start needs readjustments.
			if t.head[index] == nextNode {
				t.head[index] = nextNode.loadNext(index)
				// This was our currently highest node!
				if t.head[index] == nil {
					t.maxLevel = index - 1
				}
			}

			// Link from end needs readjustments.
			if next == nil {
				t.tail[index] = currentNode
			}
			nextNode.atomicStoreNext(index, nil)
		}

		if nextNode != nil && nextNode.key < key {
			// Go right
			currentNode = nextNode
		} else {
			// Go down
			index--
			if index < 0 {
				break
			}
		}
	}

}

// Insert inserts the given ListElement into the skiplist.
// Insert runs in approx. O(log(n))
func (t *SkipList[K, V]) Insert(key K, val V) {

	if t == nil {
		return
	}

	level := t.generateLevel(t.maxNewLevel)

	// Only grow the height of the skiplist by one at a time!
	if level > t.maxLevel {
		level = t.maxLevel + 1
		t.maxLevel = level
	}

	elem := newESLNode(level, key, val)

	t.count++

	newFirst := true
	newLast := true
	if !t.IsEmpty() {
		newFirst = elem.key < t.head[0].key
		newLast = elem.key > t.tail[0].key
	}

	normallyInserted := false
	if !newFirst && !newLast {

		normallyInserted = true

		index := t.findEntryIndex(elem.key, level)

		var currentNode, nextNode *ESLNode[K, V]
		for {

			if currentNode == nil {
				nextNode = t.head[index]
			} else {
				nextNode = currentNode.loadNext(index)
			}

			// Connect node to next
			if index <= level && (nextNode == nil || nextNode.key > elem.key) {
				elem.atomicStoreNext(index, nextNode)
				if currentNode != nil {
					currentNode.atomicStoreNext(index, elem)
				}
				if index == 0 {
					elem.prev = currentNode
					if nextNode != nil {
						nextNode.prev = elem
					}
				}
			}

			if nextNode != nil && nextNode.key <= elem.key {
				// Go right
				currentNode = nextNode
			} else {
				// Go down
				index--
				if index < 0 {
					break
				}
			}
		}
	}

	// Where we have a left-most position that needs to be referenced!
	for i := level; i >= 0; i-- {

		didSomething := false

		if newFirst || normallyInserted {

			if t.head[i] == nil || t.head[i].key > elem.key {
				if i == 0 && t.head[i] != nil {
					t.head[i].prev = elem
				}
				elem.atomicStoreNext(i, t.head[i])
				t.head[i] = elem
			}

			// link the endLevels to this element!
			if elem.atomicLoadNext(i) == nil {
				t.tail[i] = elem
			}

			didSomething = true
		}

		if newLast {
			// Places the element after the very last element on this level!
			// This is very important, so we are not linking the very first element (newFirst AND newLast) to itself!
			if !newFirst {
				if t.tail[i] != nil {
					t.tail[i].atomicStoreNext(i, elem)
				}
				if i == 0 {
					elem.prev = t.tail[i]
				}
				t.tail[i] = elem
			}

			// Link the startLevels to this element!
			if t.head[i] == nil || t.head[i].key > elem.key {
				t.head[i] = elem
			}

			didSomething = true
		}

		if !didSomething {
			break
		}
	}
}


// GetSmallestNode returns the very first/smallest node in the skiplist.
// GetSmallestNode runs in O(1)
func (t *SkipList[K, V]) Head() *ESLNode[K, V] {
	return t.head[0]
}

// GetLargestNode returns the very last/largest node in the skiplist.
// GetLargestNode runs in O(1)
func (t *SkipList[K, V]) Tail() *ESLNode[K, V] {
	return t.tail[0]
}

// GetNodeCount returns the number of nodes currently in the skiplist.
func (t *SkipList[K, V]) Count() int {
	return t.count
}

// // ChangeValue can be used to change the actual value of a node in the skiplist
// // without the need of Deleting and reinserting the node again.
// // Be advised, that ChangeValue only works, if the actual key from ExtractKey() will stay the same!
// // ok is an indicator, wether the value is actually changed.
// func (t *SkipList[K, V]) ChangeValue(e *SkipListElement, newValue ListElement) (ok bool) {
// 	// The key needs to stay correct, so this is very important!
// 	if math.Abs(newValue.ExtractKey() - e.key) <= t.eps {
// 		e.value = newValue
// 		ok = true
// 	} else {
// 		ok = false
// 	}
// 	return
// }

// // String returns a string format of the skiplist. Useful to get a graphical overview and/or debugging.
// func (t *SkipList[K, V]) String() string {
// 	s := ""

// 	s += " --> "
// 	for i, l := range t.startLevels {
// 		if l == nil {
// 			break
// 		}
// 		if i > 0 {
// 			s += " -> "
// 		}
// 		next := "---"
// 		if l != nil {
// 			next = l.value.String()
// 		}
// 		s += fmt.Sprintf("[%v]", next)

// 		if i == 0 {
// 			s += "    "
// 		}
// 	}
// 	s += "\n"

// 	node := t.startLevels[0]
// 	for node != nil {
// 		s += fmt.Sprintf("%v: ", node.value)
// 		for i := 0; i <= node.level; i++ {

// 			l := node.next[i]

// 			next := "---"
// 			if l != nil {
// 				next = l.value.String()
// 			}

// 			if i == 0 {
// 				prev := "---"
// 				if node.prev != nil {
// 					prev = node.prev.value.String()
// 				}
// 				s += fmt.Sprintf("[%v|%v]", prev, next)
// 			} else {
// 				s += fmt.Sprintf("[%v]", next)
// 			}
// 			if i < node.level {
// 				s += " -> "
// 			}

// 		}
// 		s += "\n"
// 		node = node.next[0]
// 	}

// 	s += " --> "
// 	for i, l := range t.endLevels {
// 		if l == nil {
// 			break
// 		}
// 		if i > 0 {
// 			s += " -> "
// 		}
// 		next := "---"
// 		if l != nil {
// 			next = l.value.String()
// 		}
// 		s += fmt.Sprintf("[%v]", next)
// 		if i == 0 {
// 			s += "    "
// 		}
// 	}
// 	s += "\n"
// 	return s
// }
package esl

import (
	"sync/atomic"
	"unsafe"

	cst "golang.org/x/exp/constraints"
)

//
// a non-blocking and thread-safe skip list implementation
// rebuild from https://github.com/zhangyunhao116/skipset
//

const (
	MAX_LEVEL  = 16               // max level is based on P, MAX_LEVEL = 64 * P
	P          = 0.25             // factor to raise node level, must be less than 1, recommended less than 0.5, and 0.25 is a good choice, redis use it
	INIT_LEVEL = 1                // 1 is ok
	_VER       = "1.1.0"
)

type ESL[K cst.Ordered, V any] struct {
	length   int64
	level    uint64           // cur highest level, it will changed in operations.
	header   Node[K, V]
	tail     *Node[K, V]
}

func New[K cst.Ordered, V any]() *ESL[K, V] {
	out := &ESL[K, V]{
		level:  INIT_LEVEL,
	}
	out.header.flags.SetTrue(fullyLinked)
	return out
}

func (s *ESL[K, V]) Add(key K, val V) bool {
	_, ok := s.add(key, val, false)
	return ok
}

func (s *ESL[K, V]) Set(key K, val V) (prev V, replaced bool) {
	return s.add(key, val, true)
}

// Get returns the value of the node with the specified key in the list.
// if key not found, return default value(if given, else return zero value)
func (s *ESL[K, V]) Get(key K, df... V) (V, bool) {
	n := s.find(key)
	if n != nil {
		return n.Val, true
	}

	if len(df) > 0 {
		return df[0], false
	}

	var result V
	return result, false
}

// Val returns the value of the node with the specified key in the list.
// if key not found, return default value(if given, else return zero value)
func (s *ESL[K, V]) Val(key K, df... V) (V) {
	n := s.find(key)
	if n != nil {
		return n.Val
	}

	if len(df) > 0 {
		return df[0]
	}

	var result V
	return result
}

// Keys returns all keys in the list.
func (s *ESL[K, V]) Keys() (keys []K) {
	s.traverseNode(func(n *Node[K, V]) bool {
		keys = append(keys, n.key)
		return true
	})

	return
}

// Vals returns all values in the list.
func (s *ESL[K, V]) Vals() (vals []V) {
	s.traverseNode(func(n *Node[K, V]) bool {
		vals = append(vals, n.Val)
		return true
	})

	return
}

// Find returns the node with the specified key in the list.
// if key not found, return nil
func (s *ESL[K, V]) Find(key K) (*Node[K, V]) {
	return s.find(key)
}

// Contains checks if the value is in the skip set.
func (s *ESL[K, V]) Contains(key K) bool {
	return s.find(key) != nil
}

// First returns the first node in the list.
func (s *ESL[K, V]) First() (*Node[K, V]) {
	return s.header.Next()
}

// Last returns the last node in the list. 
func (s *ESL[K, V]) Last() (*Node[K, V]) {
	return s.tail
}

// IsEmpty returns true if the list is empty.
func (s *ESL[K, V]) IsEmpty() bool {
	return s.header.Next() == nil
}

// Len returns the number of elements in the list.
func (s *ESL[K, V]) Size() int64 {
	return atomic.LoadInt64(&s.length)
}

// Len returns the number of elements in the list.
func (s *ESL[K, V]) Len() int64 {
	return atomic.LoadInt64(&s.length)
}

// Clear clears the list.
func (s *ESL[K, V]) Clear() {
	s.header.mu.Lock()
	defer s.header.mu.Unlock()

	s.header.next = [MAX_LEVEL]unsafe.Pointer{}
	s.tail  = nil
	s.level = INIT_LEVEL
	atomic.StoreInt64(&s.length, 0)
}

// Del deletes the elem with the specified key in the list.
func (s *ESL[K, V]) Del(key K) (bool) {
	return s.Remove(key)
}

// Dels delete the elems with the specified keys in the list.
func (s *ESL[K, V]) Dels(keys ...K) (cnt int64) {
	for _, key := range keys {
		if s.Del(key) {
			cnt++
		}
	}
	return
}

// PopFirst delete the first elem in the list.
func (s *ESL[K, V]) PopFirst() (*Node[K, V]) {
  nn := s.header.Next()
	if nn != nil {
		if s.Remove(nn.key) {
			return nn
		}
	}
	return nil
}

// Range calls the function `cb` for all items in the list, can limit by param limit.
func (s *ESL[K, V]) Range(cb func(K, V) bool, limit... int) {
	if len(limit) == 0 {
		s.traverse(cb)
		return
	}

	limit_ := limit[0]
	s.traverse(func(k K, v V ) bool {
		if limit_ <= 0 {
			return false
		}

		if !cb(k, v) {
			return false
		}

		limit_--
		return true
	})
}

// RangeFrom traverse the items in range [start, ...] by call function cb
func (s *ESL[K, V]) RangeFrom(start K, cb func(K, V) bool, limit... int) {

	if len(limit) == 0 {
		s.traverseFrom(start, cb)
		return
	}

	limit_ := limit[0]
	s.traverseFrom(start, func(k K, v V ) bool {
		if limit_ <= 0 {
			return false
		}

		if !cb(k, v) {
			return false
		}

		limit_--
		return true
	})
}

// RangeFromTo range the items in [start, end] by call function cb
func (s *ESL[K, V]) RangeFromTo(start K, end K, cb func(K, V) bool, limit... int) {

	if len(limit) == 0 {
		s.traverseFromTo(start, end, cb)
		return
	}

	limit_ := limit[0]

	s.traverseFromTo(start, end, func(k K, v V ) bool {
		if limit_ <= 0 {
			return false
		}

		if !cb(k, v) {
			return false
		}

		limit_--
		return true
	})
}

// RangeIn range the items in [start, end) by call function cb
func (s *ESL[K, V]) RangeIn(start K, end K, cb func(K, V) bool, limit... int) {

	if len(limit) == 0 {
		s.traverseIn(start, end, cb)
		return
	}

	limit_ := limit[0]

	s.traverseIn(start, end, func(k K, v V ) bool {
		if limit_ <= 0 {
			return false
		}

		if !cb(k, v) {
			return false
		}

		limit_--
		return true
	})
}

func Version() string {
	return _VER
}
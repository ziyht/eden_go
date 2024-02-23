package earr

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/ziyht/eden_go/esync"
	"golang.org/x/exp/slices"
)

const (
	growthFactor = float32(2.0)  // growth by 100%
	shrinkFactor = float32(0.25) // shrink when size is 25% of capacity (0 means never shrink)
)

// EArr holds the elements in a slice
type EArr [T comparable] struct {
	elements []T
	size     int
	opcnt    int32
	mu       esync.RWMutex
}

// New instantiates a new earr and adds the passed values, if any, to the earr
func New [T comparable] (safe ...bool) *EArr[T] {
	return &EArr[T]{mu: esync.NewRWMutex(safe...)}
}

// Add appends a value at the end of the earr
func (earr *EArr[T]) Add(values ...T) {
	earr.mu.Lock()
	defer earr.mu.Unlock()

	earr.growBy(len(values))
	for _, value := range values {
		earr.elements[earr.size] = value
		earr.size++
	}
}

// Get returns the element at index.
// Second return parameter is true if index is within bounds of the array and array is not empty, otherwise false.
func (earr *EArr[T]) Get(index int) (e T, ok bool) {
	earr.mu.Lock()
	defer earr.mu.Unlock()

	if !earr.withinRange(index) {
		return
	}

	return earr.elements[index], true
}

// Del removes the element at the given index from the earr.
func (earr *EArr[T]) Del(index int) {
	earr.mu.Lock()
	defer earr.mu.Unlock()

	if !earr.withinRange(index) {
		return
	}

	// cleanup reference, this is needed if T is a pointer type
	{
		var e T
		earr.elements[index] = e  
	}
	                                  
	copy(earr.elements[index:], earr.elements[index+1:earr.size]) // shift to the left by one (slow operation, need ways to optimize this)
	earr.size--

	earr.shrink()
}

// Contains checks if elements (one or more) are present in the set.
// All elements have to be present in the set for the method to return true.
// Performance time complexity of n^2.
// Returns true if no arguments are passed at all, i.e. set is always super-set of empty set.
func (earr *EArr[T]) Contains(values ...T) bool {
	earr.mu.Lock()
	defer earr.mu.Unlock()

	for _, searchValue := range values {
		found := false
		for index := 0; index < earr.size; index++ {
			if earr.elements[index] == searchValue {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Values returns all elements in the earr.
func (earr *EArr[T]) Values() []T {
	newElements := make([]T, earr.size)
	copy(newElements, earr.elements[:earr.size])
	return newElements
}

//IndexOf returns index of provided element
func (earr *EArr[T]) IndexOf(value interface{}) int {
	if earr.size == 0 {
		return -1
	}
	for index, element := range earr.elements {
		if element == value {
			return index
		}
	}
	return -1
}

// Empty returns true if earr does not contain any elements.
func (earr *EArr[T]) Empty() bool {
	return earr.size == 0
}

// Size returns number of elements within the earr.
func (earr *EArr[T]) Size() int {
	return earr.size
}

// Clear removes all elements from the earr.
func (earr *EArr[T]) Clear() {
	earr.size = 0
	earr.elements = []T{}
}

// Sort sorts values (in-place) using.
func (earr *EArr[T]) Sort(less func(a, b T) bool) {
	if len(earr.elements) < 2 {
		return
	}

	slices.SortFunc(earr.elements[:earr.size], less)
}

// Swap swaps the two values at the specified positions.
func (earr *EArr[T]) Swap(i, j int) {
	if earr.withinRange(i) && earr.withinRange(j) {
		earr.elements[i], earr.elements[j] = earr.elements[j], earr.elements[i]
	}
}

// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Does not do anything if position is negative or bigger than earr's size
// Note: position equal to earr's size is valid, i.e. append.
func (earr *EArr[T]) Insert(index int, values ...T) {

	if !earr.withinRange(index) {
		// Append
		if index == earr.size {
			earr.Add(values...)
		}
		return
	}

	l := len(values)
	earr.growBy(l)
	earr.size += l
	copy(earr.elements[index+l:], earr.elements[index:earr.size-l])
	copy(earr.elements[index:], values)
}

// Set the value at specified index
// Does not do anything if position is negative or bigger than earr's size
// Note: position equal to earr's size is valid, i.e. append.
func (earr *EArr[T]) Set(index int, value T) {

	if !earr.withinRange(index) {
		// Append
		if index == earr.size {
			earr.Add(value)
		}
		return
	}

	earr.elements[index] = value
}

// String returns a string representation of container
func (earr *EArr[T]) String() string {
	str := "Arrayearr\n"
	values := []string{}
	for _, value := range earr.elements[:earr.size] {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

// Check that the index is within bounds of the earr
func (earr *EArr[T]) withinRange(index int) bool {
	return index >= 0 && index < earr.size
}

func (earr *EArr[T]) resize(cap int) {
	newElements := make([]T, cap)
	copy(newElements, earr.elements)
	earr.elements = newElements
}

// Expand the array if necessary, i.e. capacity will be reached if we add n elements
func (earr *EArr[T]) growBy(n int) {
	// When capacity is reached, grow by a factor of growthFactor and add number of elements

	currentCapacity := cap(earr.elements)

	needed := earr.size + n - currentCapacity
	if needed > 0 {
		newCapacity := int(growthFactor * float32(currentCapacity+n))
		if newCapacity <= 10240 {
			earr.resize(newCapacity)
		} else {
			earr.resize(earr.size + n + 10240)
		}
	}
}

// Shrink the array if necessary, i.e. when size is shrinkFactor percent of current capacity
func (earr *EArr[T]) shrink() {
	if shrinkFactor == 0.0 {
		return
	}

	// check 100 times to ignore frequency operations may happen in some cases
	if atomic.AddInt32(&earr.opcnt, 1) % 100 != 0 {
		return
	}

	// Shrink when size is at shrinkFactor * capacity
	currentCapacity := cap(earr.elements)
	if earr.size <= int(float32(currentCapacity)*shrinkFactor) {
		earr.resize(earr.size)
	}
}
package etimer

import (
	heap "github.com/ziyht/eden_go/etimer/heap"
	"math"
	"sync"
	"sync/atomic"
)

// priorityQueueItem stores the queue item which has a `priority` attribute to sort itself in heap.
type priorityQueueItem struct {
	value    interface{}
	priority int64
}

// priorityQueueHeap is a heap manager, of which the underlying `array` is a array implementing a heap structure.
type priorityQueueHeap struct {
	array []priorityQueueItem
}

// priorityQueue is an abstract data type similar to a regular queue or stack data structure in which
// each element additionally has a "priority" associated with it. In a priority queue, an element with
// high priority is served before an element with low priority.
// priorityQueue is based on heap structure.
type priorityQueue struct {
	mu           sync.Mutex
	heap         *priorityQueueHeap // the underlying queue items manager using heap.
	nextPriority int64              // nextPriority stores the next priority value of the heap, which is used to check if necessary to call the Pop of heap by Timer.
}


// newPriorityQueue creates and returns a priority queue.
func newPriorityQueue() *priorityQueue {
	queue := &priorityQueue{
		heap:         &priorityQueueHeap{array: make([]priorityQueueItem, 0)},
		nextPriority: math.MaxInt64,
	}
	heap.Init(queue.heap)
	return queue
}

// NextPriority retrieves and returns the minimum and the most priority value of the queue.
func (q *priorityQueue) NextPriority() int64 {
	return  atomic.LoadInt64(&q.nextPriority)
}

// Push pushes a value to the queue.
// The `priority` specifies the priority of the value.
// The lesser the `priority` value the higher priority of the `value`.
func (q *priorityQueue) Push(value interface{}, priority int64) {
	q.mu.Lock()
	defer q.mu.Unlock()
	heap.Push(q.heap, priorityQueueItem{
		value:    value,
		priority: priority,
	})
	// Update the minimum priority using atomic operation.
	nextPriority := atomic.LoadInt64(&q.nextPriority)
	if priority >= nextPriority {
		return
	}
	atomic.StoreInt64(&q.nextPriority, priority)
}

// Pop retrieves, removes and returns the most high priority value from the queue.
func (q *priorityQueue) Pop() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if v := heap.Pop(q.heap); v != nil {
		var nextPriority int64 = math.MaxInt64
		if len(q.heap.array) > 0 {
			nextPriority = q.heap.array[0].priority
		}
		atomic.StoreInt64(&q.nextPriority, nextPriority)
		return v.(priorityQueueItem).value
	}
	return nil
}

// Pop retrieves, removes and returns the most high priority value from the queue.
func (q *priorityQueue) Fetch() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if v := heap.Fetch(q.heap); v != nil {
		var nextPriority int64 = math.MaxInt64
		if len(q.heap.array) > 0 {
			nextPriority = q.heap.array[0].priority
		}
		atomic.StoreInt64(&q.nextPriority, nextPriority)
		return v.(priorityQueueItem).value
	}
	return nil
}

// Len is used to implement the interface of sort.Interface.
func (h *priorityQueueHeap) Len() int {
	return len(h.array)
}

// Less is used to implement the interface of sort.Interface.
// The least one is placed to the top of the heap.
func (h *priorityQueueHeap) Less(i, j int) bool {
	return h.array[i].priority < h.array[j].priority
}

// Swap is used to implement the interface of sort.Interface.
func (h *priorityQueueHeap) Swap(i, j int) {
	if len(h.array) == 0 {
		return
	}
	h.array[i], h.array[j] = h.array[j], h.array[i]
}

// Push pushes an item to the heap.
func (h *priorityQueueHeap) Push(x interface{}) {
	h.array = append(h.array, x.(priorityQueueItem))
}

// Pop retrieves, removes and returns the most high priority item from the heap.
func (h *priorityQueueHeap) Pop() interface{} {
	length := len(h.array)
	if length == 0 {
		return nil
	}
	item := h.array[length-1]
	h.array = h.array[0 : length-1]
	return item
}

// Pop retrieves, removes and returns the most high priority item from the heap.
func (h *priorityQueueHeap) Fetch() interface{} {
	length := len(h.array)
	if length == 0 {
		return nil
	}
	item := h.array[length-1]
	return item
}

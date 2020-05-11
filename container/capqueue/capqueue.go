/*
Package capqueue implements a key-value priority queue with limited number of entries.
This differs from a standard heap in that it maintains a doubly-linked list running through all of its entries.
When a new entry is added to a full queue, the oldest element (not the element with lowest priority) gets deleted.

The underlying heap implementation uses container/heap which is based on a binary heap, providing O(log n)
complexity for q.Add() and q.Remove() and O(1) for q.Max().
*/
package capqueue

import (
	"container/heap"
	"container/list"
)

// CapQueue represents a priority queue with limited number of entries.
type CapQueue struct {
	heap binHeap
	cap  int

	index map[string]*item
	order *list.List
}

// item represents one entry of CapQueue.
type item struct {
	*list.Element // position of the item in the list

	key   string
	value int
	index int // index of the item in the heap<
}

// binary heap of the items
type binHeap []*item

// New crates a new CapQueue instance.
func New(cap int) *CapQueue {
	h := &CapQueue{
		heap:  make(binHeap, 0, cap),
		cap:   cap,
		index: make(map[string]*item, cap),
		order: list.New(),
	}
	heap.Init(&h.heap)
	return h
}

// Add adds a new key-value pair to the queue.
// If the queue is already full, the oldest element gets removed.
func (h *CapQueue) Add(key string, value int) {
	var it *item
	// assure that there is always space in the heap
	if h.Len() == h.cap {
		it = h.first()
		h.order.Remove(it.Element)
		delete(h.index, it.key)
		// replace with new key/value
		it.key = key
		it.value = value
		heap.Fix(&h.heap, it.index)
	} else {
		// create a new item
		it = &item{key: key, value: value}
		heap.Push(&h.heap, it)
	}
	// add the item to the map and list
	h.index[key] = it
	it.Element = h.order.PushBack(it)
}

// Delete removes the element with the given key.
// It returns true, if an element was removed or false when no element with the given key exists.
func (h *CapQueue) Delete(key string) bool {
	it, ok := h.index[key]
	if !ok {
		return false
	}

	delete(h.index, it.key)
	h.order.Remove(it.Element)
	heap.Remove(&h.heap, it.index)
	return true
}

// Value returns the value of the given key or 0 if no such key exists.
func (h *CapQueue) Value(key string) int {
	it, ok := h.index[key]
	if !ok {
		return 0
	}
	return it.value
}

// Len returns the number of elements contained in the queue.
// The number of elements will never be larger than the initial capacity of the queue.
func (h *CapQueue) Len() int {
	return h.heap.Len()
}

// Cap returns the maximum capacity of the queue.
func (h *CapQueue) Cap() int {
	return h.cap
}

// Max returns the key-value pair with the highest value.
// This will panic if the queue is empty.
func (h *CapQueue) Max() (string, int) {
	if h.Len() == 0 {
		panic("empty queue")
	}
	it := h.heap[0]
	return it.key, it.value
}

// First returns the oldest key-value pair.
// This returns the element that was added to the queue first, not the one with the lowest value.
// If more than capacity elements are added to the queue, the oldest element gets removed.
func (h *CapQueue) First() (string, int) {
	if h.Len() == 0 {
		panic("empty queue")
	}
	it := h.first()
	return it.key, it.value
}

// first returns the oldest element in the queue.
func (h *CapQueue) first() *item {
	return h.order.Front().Value.(*item)
}

func (h binHeap) Len() int {
	return len(h)
}

func (h binHeap) Less(i, j int) bool {
	return h[i].value > h[j].value
}

func (h binHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *binHeap) Push(x interface{}) {
	n := len(*h)
	if n == cap(*h) {
		panic("insufficient capacity")
	}
	item := x.(*item)
	item.index = n
	*h = append(*h, item)
}

func (h *binHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*h = old[0 : n-1]
	return item
}

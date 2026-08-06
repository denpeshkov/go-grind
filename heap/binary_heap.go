// Package heap implements a binary heap.
package heap

import "iter"

// BinaryHeap is a binary max-heap.
type BinaryHeap[T any] struct {
	s    []*item[T]
	less func(a, b T) bool
}

// item is a heap element together with its index.
type item[T any] struct {
	value T
	index int // -1 for deleted item.
}

// A Handle refers to an element stored in a [BinaryHeap].
// It stays valid while the element remains in the heap.
// It is invalidated when its element is removed from the heap.
// Passing an invalidated handle to [BinaryHeap.Update] or [BinaryHeap.Remove] panics.
type Handle[T any] struct{ item *item[T] }

// New builds a heap from the elements of s, ordered by less.
// The elements are copied into the heap; s itself is not retained.
func New[T any](s []T, less func(a, b T) bool) *BinaryHeap[T] {
	h := &BinaryHeap[T]{
		s:    make([]*item[T], len(s)),
		less: less,
	}
	for i, v := range s {
		h.s[i] = &item[T]{value: v, index: i}
	}
	h.buildHeap()
	return h
}

// Len returns the number of elements in the heap.
func (h *BinaryHeap[T]) Len() int {
	return len(h.s)
}

// Push adds v to the heap and returns it's [Handle].
func (h *BinaryHeap[T]) Push(v T) Handle[T] {
	it := &item[T]{value: v, index: len(h.s)}
	h.s = append(h.s, it)
	h.siftUp(it.index)
	return Handle[T]{item: it}
}

// Peek returns the max element without removing it.
func (h *BinaryHeap[T]) Peek() T {
	if len(h.s) == 0 {
		panic("heap: heap is empty")
	}
	return h.s[0].value
}

// Pop removes and returns the max element.
// The [Handle] of the returned element is invalidated.
func (h *BinaryHeap[T]) Pop() T {
	if len(h.s) == 0 {
		panic("heap: heap is empty")
	}
	return h.removeAt(0)
}

// Update replaces the value of the element identified by handle.
// It panics if handle is invalid.
func (h *BinaryHeap[T]) Update(handle Handle[T], v T) {
	it := h.itemOf(handle)
	old := it.value
	it.value = v
	if h.less(old, v) {
		h.siftUp(it.index)
	} else {
		h.siftDown(it.index)
	}
}

// Remove removes the element identified by handle and returns its value.
// Handle is invalidated by this call. It panics if handle is invalid.
func (h *BinaryHeap[T]) Remove(handle Handle[T]) T {
	it := h.itemOf(handle)
	return h.removeAt(it.index)
}

// Merge moves every element of g into h and leaves g empty.
// Handles from both heaps remain valid and afterwards refer to h.
// g must order elements the same way as h, and g must not be h.
func (h *BinaryHeap[T]) Merge(g *BinaryHeap[T]) {
	h.s = append(h.s, g.s...)
	g.s = nil
	for i, it := range h.s {
		it.index = i
	}
	h.buildHeap()
}

// All returns an iterator over the heap's elements.
func (h *BinaryHeap[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		vals := make([]T, len(h.s))
		for i, it := range h.s {
			vals[i] = it.value
		}
		tmp := New(vals, h.less)
		for tmp.Len() > 0 {
			if !yield(tmp.Pop()) {
				return
			}
		}
	}
}

// itemOf returns the item referenced by handle.
func (h *BinaryHeap[T]) itemOf(hd Handle[T]) *item[T] {
	it := hd.item
	if it == nil || it.index < 0 || it.index >= len(h.s) || h.s[it.index] != it {
		panic("heap: invalid handle")
	}
	return it
}

// removeAt removes the element at index i, invalidating its item.
func (h *BinaryHeap[T]) removeAt(i int) T {
	it := h.s[i]
	v := it.value

	n := len(h.s) - 1
	if i != n {
		h.swap(i, n)
	}
	h.s[n] = nil // Zero out.
	h.s = h.s[:n]

	it.value, it.index = *new(T), -1 // Invalidate the handle.

	if i < n {
		h.siftDown(i)
		h.siftUp(i)
	}
	return v
}

func (h *BinaryHeap[T]) buildHeap() {
	for i := len(h.s)/2 - 1; i >= 0; i-- {
		h.siftDown(i)
	}
}

func (h *BinaryHeap[T]) swap(i, j int) {
	h.s[i], h.s[j] = h.s[j], h.s[i]
	h.s[i].index = i
	h.s[j].index = j
}

func (h *BinaryHeap[T]) siftUp(i int) {
	for i > 0 {
		p := (i - 1) / 2
		if !h.less(h.s[p].value, h.s[i].value) { // parent not less than child
			break
		}
		h.swap(i, p)
		i = p
	}
}

func (h *BinaryHeap[T]) siftDown(i int) {
	for l := 2*i + 1; l < len(h.s); l = 2*i + 1 {
		largest := l
		if r := l + 1; r < len(h.s) && h.less(h.s[l].value, h.s[r].value) {
			largest = r
		}
		if !h.less(h.s[i].value, h.s[largest].value) { // i not less than largest child
			break
		}
		h.swap(i, largest)
		i = largest
	}
}

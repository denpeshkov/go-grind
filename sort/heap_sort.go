package sort

import (
	"cmp"

	"github.com/denpeshkov/go-grind/heap"
)

// HeapSort sorts s in ascending order using a heapsort.
func HeapSort[S ~[]E, E cmp.Ordered](s S) {
	h := heap.New(s, cmp.Less)
	n := len(s) - 1
	// We can reuse original slice as [heap.BinaryHeap.Pop] decreases s length.
	// We insert elements from the end, as [heap.BinaryHeap] is a max-heap.
	for v := range h.All() {
		s[n] = v
		n--
	}
}

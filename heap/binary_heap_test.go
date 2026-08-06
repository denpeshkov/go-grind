package heap

import (
	"cmp"
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/nalgeon/be"
)

// checkHeap verifies the max-heap invariant under h.less and that every item's
// recorded index matches its position, so Handles stay accurate.
func checkHeap[T any](t *testing.T, h *BinaryHeap[T]) {
	t.Helper()
	for i, it := range h.s {
		be.Equal(t, it.index, i)
		if l := 2*i + 1; l < len(h.s) {
			be.True(t, !h.less(it.value, h.s[l].value))
		}
		if r := 2*i + 2; r < len(h.s) {
			be.True(t, !h.less(it.value, h.s[r].value))
		}
	}
}

func TestBinaryHeap(t *testing.T) {
	h := New([]int(nil), cmp.Less)

	assertPopAndSize := func(want, n int) {
		t.Helper()
		got := h.Pop()
		be.Equal(t, got, want)
		be.Equal(t, h.Len(), n)
		checkHeap(t, h)
	}

	h.Push(3)
	h.Push(4)
	h.Push(7)
	h.Push(2)

	// Pop all elements in max order.
	assertPopAndSize(7, 3)
	assertPopAndSize(4, 2)
	assertPopAndSize(3, 1)
	assertPopAndSize(2, 0)

	// Push+pop, push+pop...
	h.Push(3)
	assertPopAndSize(3, 0)
	h.Push(9)
	assertPopAndSize(9, 0)
	h.Push(1)
	assertPopAndSize(1, 0)

	// Push in random orders; they should still pop largest-first.
	words := []int{1, 2, 3, 4, 5}
	for range 100 {
		w := slices.Clone(words)
		rand.Shuffle(len(w), func(i, j int) { w[i], w[j] = w[j], w[i] })
		for _, x := range w {
			h.Push(x)
		}
		assertPopAndSize(5, 4)
		assertPopAndSize(4, 3)
		assertPopAndSize(3, 2)
		assertPopAndSize(2, 1)
		assertPopAndSize(1, 0)
	}
}

func TestBinaryHeap_Empty(t *testing.T) {
	h := New([]int(nil), cmp.Less)
	be.Equal(t, h.Len(), 0)
}

func TestBinaryHeap_Peek(t *testing.T) {
	h := New([]int{3, 1, 2}, cmp.Less)
	v := h.Peek()
	be.Equal(t, v, 3)
	be.Equal(t, h.Len(), 3) // Peek does not remove.
	checkHeap(t, h)
}

func TestBinaryHeap_Remove(t *testing.T) {
	h := New([]int(nil), cmp.Less)
	handles := map[int]Handle[int]{}
	for _, v := range []int{5, 3, 8, 1, 9, 2, 7} {
		handles[v] = h.Push(v)
	}

	got := h.Remove(handles[8]) // remove a non-root element by handle
	be.Equal(t, got, 8)
	be.Equal(t, h.Len(), 6)
	checkHeap(t, h)

	be.Equal(t, slices.Collect(h.All()), []int{9, 7, 5, 3, 2, 1})
}

func TestBinaryHeap_Update(t *testing.T) {
	h := New([]int(nil), cmp.Less)
	hs := map[int]Handle[int]{}
	for _, v := range []int{5, 3, 8, 1} {
		hs[v] = h.Push(v)
	}

	h.Update(hs[3], 20) // increase: becomes the new max
	checkHeap(t, h)
	v := h.Peek()
	be.Equal(t, v, 20)

	h.Update(hs[8], 0) // decrease: sinks to the bottom
	checkHeap(t, h)

	be.Equal(t, slices.Collect(h.All()), []int{20, 5, 1, 0})
}

// TestBinaryHeap_HandleStability is the core regression test for the handle
// API: a Handle keeps referring to its element across churn, even though the
// element's position in the backing array moves.
func TestBinaryHeap_HandleStability(t *testing.T) {
	h := New([]int(nil), cmp.Less)
	target := h.Push(-1) // smaller than everything else, so it survives the pops
	for _, v := range []int{5, 3, 8, 1, 9, 2, 7, 50, 60} {
		h.Push(v)
	}
	for range 3 {
		h.Pop() // removes the three largest and shuffles positions around
	}

	h.Update(target, 0) // still the right element
	checkHeap(t, h)

	be.Equal(t, h.Remove(target), 0)
	checkHeap(t, h)
}

func TestBinaryHeap_InvalidHandle(t *testing.T) {
	assertPanics := func(f func()) {
		t.Helper()
		defer func() { be.True(t, recover() != nil) }()
		f()
	}

	h := New([]int(nil), cmp.Less)
	hd := h.Push(1)
	h.Push(2) // keep h non-empty after removing hd
	h.Remove(hd)

	assertPanics(func() { h.Remove(hd) })            // already removed
	assertPanics(func() { h.Update(hd, 5) })         // already removed
	assertPanics(func() { h.Remove(Handle[int]{}) }) // zero handle

	other := New([]int(nil), cmp.Less)
	oh := other.Push(9)
	assertPanics(func() { h.Remove(oh) }) // handle from another heap
}

func TestBinaryHeap_Merge(t *testing.T) {
	a := New([]int(nil), cmp.Less)
	three := a.Push(3)
	a.Push(7)
	a.Push(1)
	b := New([]int{8, 2, 5}, cmp.Less)

	a.Merge(b)
	be.Equal(t, a.Len(), 6)
	be.Equal(t, b.Len(), 0)
	checkHeap(t, a)

	// A handle from a is still valid after the merge.
	a.Update(three, 100)
	checkHeap(t, a)

	be.Equal(t, slices.Collect(a.All()), []int{100, 8, 7, 5, 2, 1})
}

func TestBinaryHeap_All(t *testing.T) {
	h := New([]int{3, 1, 2, 5, 4}, cmp.Less)

	first := slices.Collect(h.All())
	second := slices.Collect(h.All())

	be.Equal(t, first, []int{5, 4, 3, 2, 1})
	be.Equal(t, second, []int{5, 4, 3, 2, 1}) // restartable
	be.Equal(t, h.Len(), 5)                   // non-consuming
	checkHeap(t, h)
}

func TestBinaryHeap_MinHeap(t *testing.T) {
	// A reversed less turns it into a min-heap.
	h := New([]int{5, 3, 8, 1}, func(a, b int) bool { return cmp.Less(b, a) })
	be.Equal(t, slices.Collect(h.All()), []int{1, 3, 5, 8})
}

func TestBinaryHeap_NaN(t *testing.T) {
	// cmp.Less gives a total order (NaN sorts below everything), so the heap
	// invariant holds even with NaN present.
	nan := math.NaN()
	h := New([]float64{nan, 1, 5, nan, 3}, cmp.Less)
	checkHeap(t, h)
	v := h.Peek()
	be.Equal(t, v, 5.0)
}

func TestBinaryHeap_NilCmpPanics(t *testing.T) {
	defer func() { be.True(t, recover() != nil) }()
	New([]int{1, 2}, nil)
}

func TestBinaryHeap_Random(t *testing.T) {
	for range 200 {
		n := rand.IntN(50)
		var want []int
		h := New([]int(nil), cmp.Less)
		for range n {
			v := rand.IntN(100)
			want = append(want, v)
			h.Push(v)
			checkHeap(t, h)
		}
		slices.SortFunc(want, func(a, b int) int { return cmp.Compare(b, a) }) // descending
		be.Equal(t, slices.Collect(h.All()), want)
	}
}

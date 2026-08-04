package sort_test

import (
	"cmp"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/denpeshkov/go-grind/sort"
	"github.com/nalgeon/be"
)

var data = []int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}

func testEmptyNilSlice[S ~[]E, E cmp.Ordered](t *testing.T, f func(S)) {
	emptySlice := []E{}
	var nilSlice []E

	panics := func(f func()) (b bool) {
		defer func() {
			if x := recover(); x != nil {
				b = true
			}
		}()
		f()
		return false
	}

	be.True(t, !panics(func() { f(emptySlice) }))
	be.True(t, !panics(func() { f(nilSlice) }))
}

func testData[S ~[]E, E cmp.Ordered](t *testing.T, f func(S), s S) {
	s = slices.Clone(s)
	f(s)
	be.True(t, slices.IsSorted(s))
}

func testRandomInts[E sort.Integer](t *testing.T, f func([]E)) {
	n := 10_000
	s := make([]E, n)
	for i := range n {
		s[i] = rand.N(E(127))
	}
	be.True(t, !slices.IsSorted(s))

	f(s)
	be.True(t, slices.IsSorted(s))
}

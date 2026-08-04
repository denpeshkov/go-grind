package partition_test

import (
	"cmp"
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/denpeshkov/go-grind/partition"
	"github.com/nalgeon/be"
)

// checkPartition asserts the invariant and that s is a permutation of orig.
func checkPartition[T cmp.Ordered](t *testing.T, orig, s []T, pivot T, i, j int) {
	t.Helper()

	be.True(t, 0 <= i && i <= j && j <= len(s))

	for _, v := range s[:i] {
		be.True(t, cmp.Compare(v, pivot) < 0)
	}
	for _, v := range s[i:j] {
		be.True(t, cmp.Compare(v, pivot) == 0)
	}
	for _, v := range s[j:] {
		be.True(t, cmp.Compare(v, pivot) > 0)
	}

	got, want := slices.Clone(s), slices.Clone(orig)
	slices.Sort(got)
	slices.Sort(want)
	// We can't use be.Equal as slices can contain NaN.
	be.True(t, slices.EqualFunc(got, want, func(a, b T) bool { return cmp.Compare(a, b) == 0 }))
}

func TestPartition3Way(t *testing.T) {
	tests := []struct {
		name  string
		s     []int
		pivot int
		i, j  int
	}{
		{"nil", nil, 5, 0, 0},
		{"empty", []int{}, 5, 0, 0},
		{"single less", []int{1}, 5, 1, 1},
		{"single equal", []int{5}, 5, 0, 1},
		{"single greater", []int{9}, 5, 0, 0},
		{"all less", []int{1, 2, 3, 4}, 5, 4, 4},
		{"all greater", []int{6, 7, 8}, 5, 0, 0},
		{"all equal", []int{5, 5, 5}, 5, 0, 3},
		{"pivot absent", []int{7, 1, 9, 2, 8}, 5, 2, 2},
		{"pivot with dups", []int{5, 3, 5, 8, 1, 5, 9}, 5, 2, 5},
		{"already grouped", []int{1, 2, 5, 5, 8, 9}, 5, 2, 4},
		{"reverse", []int{9, 8, 5, 3, 1}, 5, 2, 3},
		{"negatives", []int{-3, 5, -10, 5, 0, -3}, 0, 3, 4},
		{"pivot below all", []int{3, 4, 5}, 1, 0, 0},
		{"pivot above all", []int{3, 4, 5}, 9, 3, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := slices.Clone(tt.s)
			i, j := partition.ThreeWay(s, tt.pivot)

			be.Equal(t, i, tt.i)
			be.Equal(t, j, tt.j)
			checkPartition(t, tt.s, s, tt.pivot, i, j)
		})
	}
}

func TestPartition3Way_Types(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		s := []string{"banana", "apple", "cherry", "banana", "date"}
		sc := slices.Clone(s)
		i, j := partition.ThreeWay(sc, "banana")
		checkPartition(t, s, sc, "banana", i, j)
	})
	t.Run("float64", func(t *testing.T) {
		s := []float64{3.5, 1.2, 3.5, -0.5, 9.9, 3.5}
		sc := slices.Clone(s)
		i, j := partition.ThreeWay(sc, 3.5)
		checkPartition(t, s, sc, 3.5, i, j)
	})
}

func TestPartition3Way_NaN(t *testing.T) {
	t.Run("nan elements", func(t *testing.T) {
		s := []float64{math.NaN(), 1.0, math.NaN(), 5.0, 3.0}
		sc := slices.Clone(s)
		i, j := partition.ThreeWay(sc, 3.0)
		// The two NaNs sort with the values below the pivot.
		be.Equal(t, i, 3)
		be.Equal(t, j, 4)
		checkPartition(t, s, sc, 3.0, i, j)
	})
	t.Run("nan pivot", func(t *testing.T) {
		s := []float64{1.0, math.NaN(), 5.0, math.NaN()}
		sc := slices.Clone(s)
		i, j := partition.ThreeWay(sc, math.NaN())
		// Every non-NaN is greater than a NaN pivot; the NaNs are "equal".
		be.Equal(t, i, 0)
		be.Equal(t, j, 2)
		checkPartition(t, s, sc, math.NaN(), i, j)
	})
}

func TestPartition3Way_Random(t *testing.T) {
	for range 2000 {
		n := rand.IntN(64)
		s := make([]int, n)
		for i := range s {
			// Small value range → many duplicates and frequent pivot hits.
			s[i] = rand.IntN(10) - 5
		}
		orig := slices.Clone(s)

		var pivot int
		if n > 0 && rand.IntN(2) == 0 {
			pivot = s[rand.IntN(n)] // present in the slice
		} else {
			pivot = rand.IntN(12) - 6 // possibly absent
		}
		i, j := partition.ThreeWay(s, pivot)
		checkPartition(t, orig, s, pivot, i, j)
	}
}

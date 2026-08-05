package sort

import (
	"cmp"
	"math/rand/v2"

	"github.com/denpeshkov/go-grind/partition"
)

// QuickSort sorts s in ascending order using a quicksort.
func QuickSort[S ~[]E, E cmp.Ordered](s S) {
	if len(s) <= 1 {
		return
	}
	pivot := s[rand.IntN(len(s))]
	i, j := partition.ThreeWay(s, pivot)
	QuickSort(s[:i])
	QuickSort(s[j:])
}

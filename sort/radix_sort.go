package sort

import (
	"math"
	"slices"
	"unsafe"
)

// RadixSort sorts s in ascending order using an LSD radix sort.
func RadixSort[S ~[]E, E Unsigned](s S) {
	aux := make(S, len(s))
	// countingSort implements a stable counting sort on byte i of each element.
	countingSort := func(s S, i int) {
		cnts := make([]uint, math.MaxUint8+1)
		for _, v := range s {
			b := uint(v>>(i*8)) & 0xFF
			cnts[b]++
		}
		for i := 1; i < math.MaxUint8+1; i++ {
			cnts[i] += cnts[i-1]
		}
		for _, v := range slices.Backward(s) {
			b := uint(v>>(i*8)) & 0xFF
			aux[cnts[b]-1] = v
			cnts[b]--
		}
		copy(s, aux)
	}
	// One pass per byte of E, so narrow types don't do wasted no-op passes.
	var zero E
	for i := range int(unsafe.Sizeof(zero)) {
		countingSort(s, i)
	}
}

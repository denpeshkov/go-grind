package partition

import (
	"cmp"
)

// ThreeWay partitions s in place into three contiguous groups:
//   - s[:i] < pivot
//   - s[i:j] == pivot
//   - s[j:] > pivot
func ThreeWay[S ~[]E, E cmp.Ordered](s S, pivot E) (i int, j int) {
	// Invariant:
	// - s[:i] < pivot
	// - s[i:j] == pivot
	// - s[j:k+1] unexamined
	// - s[k+1:] > pivot
	i, j, k := 0, 0, len(s)-1

	for j <= k {
		switch c := cmp.Compare(s[j], pivot); {
		case c < 0:
			s[i], s[j] = s[j], s[i]
			i++
			j++
		case c == 0:
			j++
		default: // c > 0
			s[j], s[k] = s[k], s[j]
			k--
		}
	}
	return i, j
}

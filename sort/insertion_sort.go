package sort

import "cmp"

// InsertionSort sorts s in ascending order using an insertion sort.
func InsertionSort[S ~[]E, E cmp.Ordered](s S) {
	for i := 1; i < len(s); i++ {
		j, v := i, s[i]
		for j > 0 && s[j-1] > v {
			s[j] = s[j-1]
			j--
		}
		s[j] = v
	}
}

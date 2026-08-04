package sort

import "cmp"

// SelectionSort sorts s in ascending order using a selection sort.
func SelectionSort[S ~[]E, E cmp.Ordered](s S) {
	for i := range s {
		mini := i
		for j := i + 1; j < len(s); j++ {
			if s[j] < s[mini] {
				mini = j
			}
		}
		s[i], s[mini] = s[mini], s[i]
	}
}

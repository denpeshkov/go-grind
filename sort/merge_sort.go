package sort

import "cmp"

// MergeSort sorts s in ascending order using a merge sort.
func MergeSort[S ~[]E, E cmp.Ordered](s S) {
	merge := func(l, r, aux S) {
		i, j, k := 0, 0, 0
		for i < len(l) && j < len(r) {
			if l[i] <= r[j] {
				aux[k] = l[i]
				i++
			} else {
				aux[k] = r[j]
				j++
			}
			k++
		}
		k += copy(aux[k:], l[i:])
		copy(aux[k:], r[j:])
	}

	var sort func(s, aux S)
	sort = func(s, aux S) {
		if len(s) <= 1 {
			return
		}
		m := len(s) / 2
		sort(s[:m], aux[:m])
		sort(s[m:], aux[m:])
		merge(s[:m], s[m:], aux)
		copy(s, aux)
	}

	sort(s, make(S, len(s)))
}

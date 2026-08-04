package sort

import (
	"cmp"
	"slices"
)

// BucketSort sorts s in ascending order defined by key using a bucket sort.
// key must map each element to the half-open interval [0,1).
func BucketSort[S ~[]E, E any](s S, key func(E) float64) {
	if len(s) <= 1 {
		return
	}

	aux := make([][]E, len(s))
	for _, v := range s {
		i := int(key(v) * float64(len(s)))
		aux[i] = append(aux[i], v)
	}

	k := 0
	for _, bkt := range aux {
		slices.SortFunc(bkt, func(a, b E) int {
			return cmp.Compare(key(a), key(b))
		})
		for _, v := range bkt {
			s[k] = v
			k++
		}
	}
}

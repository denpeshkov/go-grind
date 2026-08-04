package sort

// CountingSort sorts s in ascending order using a counting sort.
//
// The memory footprint is proportional to k rather than to len(s), so
// CountingSort is only appropriate when the value range is small.
func CountingSort[S ~[]E, E Integer](s S) {
	if len(s) <= 1 {
		return
	}

	minv, maxv := s[0], s[0]
	for _, v := range s[1:] {
		minv, maxv = min(minv, v), max(maxv, v)
	}

	cnts := make([]int, uint64(maxv)-uint64(minv)+1)
	for _, v := range s {
		cnts[uint64(v)-uint64(minv)]++
	}

	k := 0
	for i, cnt := range cnts {
		v := minv + E(i)
		for range cnt {
			s[k] = v
			k++
		}
	}
}

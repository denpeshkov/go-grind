// Based on: https://github.com/jba/omap/blob/v0.7.0/ordered/map_test.go
package treap

import (
	"iter"
	"math/rand"
	"slices"
	"testing"

	"github.com/nalgeon/be"
)

// permute populates treap with the n odd keys 1, 3, …, 2n-1 inserted in random
// order, then overwrites the values of half of them (exercising Set's in-place update).
// It returns a slice indexed by key where slice[k] is the
// expected value for key k, or 0 if k is absent (all even keys, and 0/2n).
func permute(tr *Treap[int, int], n int) (slice []int) {
	perm := rand.Perm(n)
	slice = make([]int, 2*n+1)
	for i, x := range perm {
		tr.Set(2*x+1, i+1)
		slice[2*x+1] = i + 1
	}
	for i, x := range perm[:len(perm)/2] {
		tr.Set(2*x+1, i+100) // overwrite in place
		slice[2*x+1] = i + 100
	}
	return slice
}

func count(tr *Treap[int, int]) int {
	return len(slices.Collect(tr.All()))
}

func collectKeys(seq iter.Seq[*Node[int, int]]) (keys []int) {
	for n := range seq {
		keys = append(keys, n.key)
	}
	return keys
}

func TestGet(t *testing.T) {
	for n := range 11 {
		tr := &Treap[int, int]{}
		slice := permute(tr, n)
		for k, want := range slice {
			node := tr.Get(k)
			if want > 0 {
				be.True(t, node != nil)
				be.Equal(t, node.key, k)
				be.Equal(t, node.val, want)
			} else {
				be.Equal(t, node, nil)
			}
		}
	}
}

func TestSet(t *testing.T) {
	tr := &Treap[int, int]{}

	tr.Set(1, 10)
	be.Equal(t, tr.Get(1).val, 10)
	tr.Set(2, 20)
	be.Equal(t, tr.Get(2).val, 20)

	// Setting an existing key updates the value in place, without adding.
	tr.Set(1, 5)
	be.Equal(t, tr.Get(1).val, 5)
	tr.Set(1, 8)
	be.Equal(t, tr.Get(1).val, 8)
	be.Equal(t, count(tr), 2)
}

// TestSetPreservesNode checks that overwriting a key reuses the existing node.
func TestSet_PreservesNode(t *testing.T) {
	tr := &Treap[int, int]{}
	tr.Set(1, 10)
	n1 := tr.Get(1)
	tr.Set(1, 99)
	n2 := tr.Get(1)

	be.True(t, n1 == n2)
	be.Equal(t, n2.val, 99)
}

func TestDelete(t *testing.T) {
	for n := range 11 {
		tr := &Treap[int, int]{}
		slice := permute(tr, n)

		be.Equal(t, tr.Delete(0), nil) // absent (even) key
		be.Equal(t, count(tr), n)

		for _, i := range rand.Perm(n) {
			k := 2*i + 1
			node := tr.Delete(k)
			be.True(t, node != nil)
			be.Equal(t, node.key, k)
			be.Equal(t, node.val, slice[k])
			be.Equal(t, tr.Get(k), nil)
		}
		be.Equal(t, count(tr), 0)
	}
}

func TestMin(t *testing.T) {
	for n := range 11 {
		tr := &Treap[int, int]{}
		slice := permute(tr, n)
		if n == 0 {
			be.Equal(t, tr.Min(), nil)
		} else {
			be.Equal(t, tr.Min().key, 1)
			be.Equal(t, tr.Min().val, slice[1])
		}
	}
}

func TestMax(t *testing.T) {
	for n := range 11 {
		tr := &Treap[int, int]{}
		slice := permute(tr, n)
		if n == 0 {
			be.Equal(t, tr.Max(), nil)
		} else {
			be.Equal(t, tr.Max().key, 2*n-1)
			be.Equal(t, tr.Max().val, slice[2*n-1])
		}
	}
}

func TestFloorCeil(t *testing.T) {
	tr := &Treap[int, int]{}
	for _, k := range []int{10, 20, 30, 40, 50} {
		tr.Set(k, k)
	}

	be.Equal(t, tr.Floor(5), nil)     // below all
	be.Equal(t, tr.Floor(10).key, 10) // exact
	be.Equal(t, tr.Floor(25).key, 20) // between
	be.Equal(t, tr.Floor(55).key, 50) // above all

	be.Equal(t, tr.Ceil(5).key, 10)  // below all
	be.Equal(t, tr.Ceil(30).key, 30) // exact
	be.Equal(t, tr.Ceil(25).key, 30) // between
	be.Equal(t, tr.Ceil(55), nil)    // above all
}

func TestAll(t *testing.T) {
	for n := range 11 {
		tr := &Treap[int, int]{}
		slice := permute(tr, n)

		var got []int
		for node := range tr.All() {
			be.Equal(t, node.val, slice[node.key])
			got = append(got, node.key)
			if len(got) > n+5 { // too many; looping?
				break
			}
		}
		var nonzeroKeys []int
		for k, v := range slice {
			if v != 0 {
				nonzeroKeys = append(nonzeroKeys, k)
			}
		}
		be.Equal(t, got, nonzeroKeys)
	}
}

func TestAll_BreakEarly(t *testing.T) {
	tr := &Treap[int, int]{}
	permute(tr, 20) // keys 1, 3, …, 39

	var got []int
	for node := range tr.All() {
		got = append(got, node.key)
		if len(got) == 5 {
			break
		}
	}
	be.Equal(t, got, []int{1, 3, 5, 7, 9})
}

func TestRange(t *testing.T) {
	tr := &Treap[int, int]{}
	for v := range 100 { // keys 0..99
		tr.Set(v, v)
	}
	rng := func(lo, hi int) []int { return collectKeys(tr.Range(lo, hi)) }
	seq := func(lo, hi int) (s []int) {
		for k := lo; k <= hi; k++ {
			s = append(s, k)
		}
		return s
	}

	be.Equal(t, rng(40, 60), seq(40, 60))  // interior, inclusive both ends
	be.Equal(t, rng(0, 99), seq(0, 99))    // whole set
	be.Equal(t, rng(-10, 200), seq(0, 99)) // clamped to existing keys
	be.Equal(t, rng(50, 50), seq(50, 50))  // single element
	be.Equal(t, rng(200, 300), nil)        // entirely above → empty
	be.Equal(t, rng(60, 40), nil)          // lo > hi → empty
}

func TestRange_BreakEarly(t *testing.T) {
	tr := &Treap[int, int]{}
	for v := range 100 {
		tr.Set(v, v)
	}
	var got []int
	for node := range tr.Range(40, 60) {
		got = append(got, node.key)
		if len(got) == 3 {
			break
		}
	}
	be.Equal(t, got, []int{40, 41, 42})
}

func TestEmpty(t *testing.T) {
	tr := &Treap[int, int]{}
	be.Equal(t, tr.Min(), nil)
	be.Equal(t, tr.Max(), nil)
	be.Equal(t, tr.Get(5), nil)
	be.Equal(t, tr.Delete(5), nil)
	be.Equal(t, tr.Floor(5), nil)
	be.Equal(t, tr.Ceil(5), nil)
	be.Equal(t, count(tr), 0)
	be.Equal(t, len(collectKeys(tr.Range(0, 100))), 0)
}

func TestDuplicate(t *testing.T) {
	tr := &Treap[int, int]{}
	for _, v := range []int{5, 3, 8, 3, 5, 1} {
		tr.Set(v, v)
	}
	be.Equal(t, count(tr), 4)
	be.Equal(t, collectKeys(tr.All()), []int{1, 3, 5, 8})

	tr.Set(5, 500) // overwrite existing
	be.Equal(t, count(tr), 4)
	be.Equal(t, tr.Get(5).val, 500)
	be.Equal(t, collectKeys(tr.All()), []int{1, 3, 5, 8})
}

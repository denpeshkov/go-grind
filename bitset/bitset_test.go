// Based on: https://github.com/eliben/code-for-blog/blob/ec2c2dd5d32a84979f2ee30ab96c5347d54c3a54/2025/bloom/bloom_test.go
package bitset

import (
	"fmt"
	"slices"
	"testing"

	"github.com/nalgeon/be"
)

func TestNew(t *testing.T) {
	var tests = []struct {
		size    uint64
		wantLen int
	}{
		{1, 1},
		{20, 1},
		{64, 1},
		{65, 2},
		{66, 2},
		{640, 10},
		{641, 11},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.size), func(t *testing.T) {
			bs := New(tt.size)
			be.Equal(t, len(bs.s), tt.wantLen)
		})
	}
}

func TestSetGetClear(t *testing.T) {
	bs := New(68)
	oneIndices := []uint64{1, 20, 33, 61, 67}

	for _, idx := range oneIndices {
		bs.Set(idx)
	}
	for idx := range uint64(68) {
		want := slices.Contains(oneIndices, idx)
		got := bs.Get(idx)
		be.Equal(t, got, want)
	}
	for idx := range uint64(68) {
		bs.Clear(idx)
		be.True(t, !bs.Get(idx))
	}
}

func TestSetClearReturnValues(t *testing.T) {
	bs := New(64)
	be.True(t, !bs.Set(10))   // was clear
	be.True(t, bs.Set(10))    // now already set
	be.True(t, bs.Clear(10))  // was set
	be.True(t, !bs.Clear(10)) // now already clear
}

func TestLen(t *testing.T) {
	for _, size := range []uint64{1, 64, 65, 130} {
		bs := New(size)
		be.Equal(t, bs.Len(), size)
	}
}

func TestCount(t *testing.T) {
	bs := New(130)
	be.Equal(t, bs.Count(), 0)

	indices := []uint64{0, 1, 64, 129}
	for _, idx := range indices {
		bs.Set(idx)
	}
	be.Equal(t, bs.Count(), len(indices))

	// Setting an already-set bit does not change the count.
	bs.Set(0)
	be.Equal(t, bs.Count(), len(indices))

	// Clearing an already-clear bit does not change the count.
	bs.Clear(2)
	be.Equal(t, bs.Count(), len(indices))

	bs.Clear(0)
	be.Equal(t, bs.Count(), len(indices)-1)
}

func TestReset(t *testing.T) {
	bs := New(130)
	for _, idx := range []uint64{0, 63, 64, 129} {
		bs.Set(idx)
	}

	bs.Reset()

	be.Equal(t, bs.Count(), 0)
	for idx := range uint64(130) {
		be.True(t, !bs.Get(idx))
	}
}

func TestPanics(t *testing.T) {
	assertPanics := func(f func()) {
		t.Helper()
		defer func() {
			be.True(t, recover() != nil)
		}()
		f()
	}

	assertPanics(func() { New(0) })

	// Capacity is 128 bits, but the logical size is 68. Indices in [size, cap)
	// fit in the backing array yet must still panic, as must indices beyond it.
	bs := New(68)
	for _, idx := range []uint64{68, 100, 127, 128, 1000} {
		assertPanics(func() { bs.Set(idx) })
		assertPanics(func() { bs.Get(idx) })
		assertPanics(func() { bs.Clear(idx) })
	}
}

// Based on: https://github.com/eliben/code-for-blog/blob/main/2025/bloom/loom_test.go
package bloom

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/nalgeon/be"
)

func TestFilter(t *testing.T) {
	// We set parameters to get an extremely low error rate,
	// and check that we don't get any false answers.
	const (
		capacity = uint64(2000)
		fpRate   = 1e-15
	)
	m, k := Params(fpRate, capacity)
	bf := New(m, k)

	// Insert m random items, also holding them in fullSet.
	buf := make([]byte, 256)
	fullSet := make(map[string]bool)
	for range capacity {
		if _, err := rand.Read(buf); err != nil {
			t.Fatal(err)
		}
		bf.Insert(buf)
		fullSet[string(buf)] = true
	}

	// Check that true is returned for all inserted items (this can never be
	// false due to the guarantees of the Bloom filter)
	for k := range fullSet {
		be.True(t, bf.Test([]byte(k)))
	}

	// Now generate another 2000 random items; the chance of false positives is
	// so low that we expect Test to return false for all of these.
	for range capacity {
		if _, err := rand.Read(buf); err != nil {
			t.Fatal(err)
		}
		be.True(t, !bf.Test(buf))
	}
}

func TestErrorRate(t *testing.T) {
	// Test the error rate we get from a filter matches theoretical estimates.
	const (
		capacity = uint64(1000)
		fpRate   = 0.1
	)
	n, k := Params(fpRate, capacity)

	bf := New(n, k)
	buf := make([]byte, 256)
	fullSet := make(map[string]bool)
	for range capacity {
		if _, err := rand.Read(buf); err != nil {
			t.Fatal(err)
		}
		bf.Insert(buf)
		fullSet[string(buf)] = true
	}

	// Check that true is returned for all inserted items (this can never be
	// false due to the guarantees of the Bloom filter)
	for k := range fullSet {
		be.True(t, bf.Test([]byte(k)))
	}

	// Now calculate the empirical error rate, by testing a large number of
	// random items (that weren't previously inserted).
	N := 100000
	npos := 0
	for range N {
		if _, err := rand.Read(buf); err != nil {
			t.Fatal(err)
		}
		if bf.Test(buf) {
			npos++
		}
	}

	// Expect the count to be within 25% of our requested eps
	expectedFPs := float64(N) * fpRate
	nposfp64 := float64(npos)
	be.True(t, nposfp64 >= expectedFPs*0.75 && nposfp64 <= expectedFPs*1.25)
}

func TestParams(t *testing.T) {
	tests := []struct {
		fpRate   float64
		capacity uint64
		wantM    uint64
		wantK    uint64
	}{
		{0.1, 1000, 4793, 4},
		{0.01, 1000, 9586, 7},
		{0.001, 10000, 143776, 10},
		{1e-15, 2000, 143776, 50},
		{0.1, 1, 5, 4},
		{0.5, 100, 145, 1},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("fp=%v/cap=%d", tt.fpRate, tt.capacity), func(t *testing.T) {
			m, k := Params(tt.fpRate, tt.capacity)
			be.Equal(t, m, tt.wantM)
			be.Equal(t, k, tt.wantK)
		})
	}
}

func TestParamsPanics(t *testing.T) {
	// fpRate must be in (0,1).
	for _, fpRate := range []float64{-0.1, 0, 1, 2} {
		assertPanics(t, func() { Params(fpRate, 1000) })
	}
	// capacity must be > 0.
	assertPanics(t, func() { Params(0.1, 0) })
}

func TestNewPanics(t *testing.T) {
	assertPanics(t, func() { New(64, 0) }) // k must be >= 1
	assertPanics(t, func() { New(0, 4) })  // m must be >= 1
	assertPanics(t, func() { New(0, 0) })
}

func assertPanics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		be.True(t, recover() != nil)
	}()
	f()
}

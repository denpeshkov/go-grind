// Based on: https://github.com/eliben/code-for-blog/blob/main/2025/bloom/bloom.go

// Package bloom implements a Bloom filter.
package bloom

import (
	"hash/maphash"
	"math"

	"github.com/denpeshkov/go-grind/bitset"
)

// Params returns the number of bits (m) and hash functions (k) for a Bloom
// filter with the given false-positive rate and expected capacity.
func Params(fpRate float64, capacity uint64) (m uint64, k uint64) {
	if fpRate <= 0 || fpRate >= 1 {
		panic("bloom: fpRate must be in (0,1)")
	}
	if capacity == 0 {
		panic("bloom: capacity must be > 0")
	}

	log2 := math.Log(2)
	mdivn := -math.Log(fpRate) / (log2 * log2)
	m = uint64(math.Ceil(float64(capacity) * mdivn))
	k = uint64(math.Ceil(mdivn * log2))
	return m, k
}

// New creates a new [BloomFilter] having m bits, using k hash functions.
func New(m uint64, k uint64) *BloomFilter {
	if k < 1 || m < 1 {
		panic("bloom: k and m must be >=1")
	}
	return &BloomFilter{
		k:      k,
		bitset: bitset.New(m),
		seed1:  maphash.MakeSeed(),
		seed2:  maphash.MakeSeed(),
	}
}

// BloomFilter implements a Bloom filter.
//
// It is not safe for concurrent use by multiple goroutines.
type BloomFilter struct {
	k            uint64
	bitset       *bitset.BitSet
	seed1, seed2 maphash.Seed
}

// Insert adds a data item to the Bloom filter.
func (bf *BloomFilter) Insert(data []byte) {
	h1, h2 := maphash.Bytes(bf.seed1, data), maphash.Bytes(bf.seed2, data)
	for k := range bf.k {
		i := (h1 + k*h2) % bf.bitset.Len()
		bf.bitset.Set(i)
	}
}

// Test reports whether the given data item is in the Bloom filter.
// False positives are possible, but false negatives are not.
func (bf *BloomFilter) Test(data []byte) bool {
	h1, h2 := maphash.Bytes(bf.seed1, data), maphash.Bytes(bf.seed2, data)
	for k := range bf.k {
		i := (h1 + k*h2) % bf.bitset.Len()
		if !bf.bitset.Get(i) {
			return false
		}
	}
	return true
}

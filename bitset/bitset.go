// Package bitset implements a bitset.
package bitset

import "math/bits"

// BitSet implements a bitset.
//
// A BitSet is not safe for concurrent use by multiple goroutines.
type BitSet struct {
	s    []uint64
	size uint64
}

// New returns a [BitSet] that can hold size bits, indexed from 0.
// It panics if size is 0.
func New(size uint64) *BitSet {
	if size == 0 {
		panic("bitset: size must be > 0")
	}
	return &BitSet{
		s:    make([]uint64, (size-1)/64+1),
		size: size,
	}
}

// Set sets the bit at index i and reports whether it was already set.
// It panics if i is out of bounds.
func (b *BitSet) Set(i uint64) bool {
	if i >= b.size {
		panic("bitset: index out of bounds")
	}
	word, offset := i/64, i%64
	old := b.s[word] & (1 << offset)
	b.s[word] |= 1 << offset
	return old != 0
}

// Get reports whether the bit at index i is set.
// It panics if i is out of bounds.
func (b *BitSet) Get(i uint64) bool {
	if i >= b.size {
		panic("bitset: index out of bounds")
	}

	word, offset := i/64, i%64
	v := b.s[word] & (1 << offset)
	return v != 0
}

// Clear clears the bit at index i and reports whether it was previously set.
// It panics if i is out of bounds.
func (b *BitSet) Clear(i uint64) bool {
	if i >= b.size {
		panic("bitset: index out of bounds")
	}

	word, offset := i/64, i%64
	old := b.s[word] & (1 << offset)
	b.s[word] &^= 1 << offset
	return old != 0
}

// Reset turns off all bits in the set.
func (b *BitSet) Reset() {
	clear(b.s)
}

// Len returns the number of bits the set can hold.
func (b *BitSet) Len() uint64 {
	return b.size
}

// Count returns the number of set bits (the population count).
func (b *BitSet) Count() int {
	var n int
	for _, w := range b.s {
		n += bits.OnesCount64(w)
	}
	return n
}

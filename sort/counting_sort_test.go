package sort_test

import (
	"testing"

	"github.com/denpeshkov/go-grind/sort"
)

func TestCountingSort_EmptyNil(t *testing.T) {
	testEmptyNilSlice(t, sort.CountingSort[[]int])
}

func TestCountingSort_Data(t *testing.T) {
	testData(t, sort.CountingSort, data)
}

func TestCountingSort_NegativeRange(t *testing.T) {
	// Small ranges including negatives exercise the signed offset/reconstruction
	// arithmetic. int8 spans its full range, which a naive max-min+1 sizing in
	// the element type would overflow.
	t.Run("int", func(t *testing.T) {
		testData(t, sort.CountingSort[[]int], []int{-5, -100, 3, -100, 0, 50, -1})
	})
	t.Run("int8", func(t *testing.T) {
		testData(t, sort.CountingSort[[]int8], []int8{-5, -100, 3, -100, 0, 50, -1, 127, -128})
	})
}

func TestCountingSort_RandomInts(t *testing.T) {
	t.Run("int", func(t *testing.T) { testRandomInts(t, sort.CountingSort[[]int]) })
	t.Run("int8", func(t *testing.T) { testRandomInts(t, sort.CountingSort[[]int8]) })
	t.Run("int16", func(t *testing.T) { testRandomInts(t, sort.CountingSort[[]int16]) })
	t.Run("int32", func(t *testing.T) { testRandomInts(t, sort.CountingSort[[]int32]) })
	t.Run("int64", func(t *testing.T) { testRandomInts(t, sort.CountingSort[[]int64]) })
	t.Run("uint", func(t *testing.T) { testRandomInts(t, sort.CountingSort[[]uint]) })
	t.Run("uint8", func(t *testing.T) { testRandomInts(t, sort.CountingSort[[]uint8]) })
	t.Run("uint16", func(t *testing.T) { testRandomInts(t, sort.CountingSort[[]uint16]) })
	t.Run("uint32", func(t *testing.T) { testRandomInts(t, sort.CountingSort[[]uint32]) })
	t.Run("uint64", func(t *testing.T) { testRandomInts(t, sort.CountingSort[[]uint64]) })
	t.Run("uintptr", func(t *testing.T) { testRandomInts(t, sort.CountingSort[[]uintptr]) })
}

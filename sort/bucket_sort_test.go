package sort_test

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/denpeshkov/go-grind/sort"
)

// bucketSortFn adapts [sort.BucketSort] to the shared func(S) harness by building
// a key that normalizes each value into [0,1).
func bucketSortFn[S ~[]E, E sort.Integer]() func(S) {
	return func(s S) {
		if len(s) == 0 {
			sort.BucketSort(s, func(E) float64 { return 0 })
			return
		}
		minv, maxv := slices.Min(s), slices.Max(s)
		span := float64(uint64(maxv) - uint64(minv))
		sort.BucketSort(s, func(v E) float64 {
			if span == 0 {
				return 0
			}
			// Normalize into [0,1). The maximum maps to exactly 1.0, which is
			// outside the half-open interval BucketSort requires, so cap it just
			// below 1 (large ranges may also round up to 1.0 in float64).
			return min(float64(uint64(v)-uint64(minv))/span, math.Nextafter(1, 0))
		})
	}
}

func TestBucketSort_EmptyNil(t *testing.T) {
	testEmptyNilSlice(t, bucketSortFn[[]int]())
}

func TestBucketSort_Data(t *testing.T) {
	testData(t, bucketSortFn[[]int](), data)
}

func TestBucketSort_WideRange(t *testing.T) {
	// A wide value range must still sort correctly once normalized into [0,1).
	s := make([]uint64, 1000)
	for i := range s {
		s[i] = rand.Uint64N(1 << 55)
	}
	testData(t, bucketSortFn[[]uint64](), s)
}

func TestBucketSort_RandomInts(t *testing.T) {
	t.Run("int", func(t *testing.T) { testRandomInts(t, bucketSortFn[[]int]()) })
	t.Run("int8", func(t *testing.T) { testRandomInts(t, bucketSortFn[[]int8]()) })
	t.Run("int16", func(t *testing.T) { testRandomInts(t, bucketSortFn[[]int16]()) })
	t.Run("int32", func(t *testing.T) { testRandomInts(t, bucketSortFn[[]int32]()) })
	t.Run("int64", func(t *testing.T) { testRandomInts(t, bucketSortFn[[]int64]()) })
	t.Run("uint", func(t *testing.T) { testRandomInts(t, bucketSortFn[[]uint]()) })
	t.Run("uint8", func(t *testing.T) { testRandomInts(t, bucketSortFn[[]uint8]()) })
	t.Run("uint16", func(t *testing.T) { testRandomInts(t, bucketSortFn[[]uint16]()) })
	t.Run("uint32", func(t *testing.T) { testRandomInts(t, bucketSortFn[[]uint32]()) })
	t.Run("uint64", func(t *testing.T) { testRandomInts(t, bucketSortFn[[]uint64]()) })
	t.Run("uintptr", func(t *testing.T) { testRandomInts(t, bucketSortFn[[]uintptr]()) })
}

func TestBucketSort_Extremums(t *testing.T) {
	t.Run("int", func(t *testing.T) { testData(t, bucketSortFn[[]int](), []int{math.MaxInt, math.MinInt}) })
	t.Run("int8", func(t *testing.T) { testData(t, bucketSortFn[[]int8](), []int8{math.MaxInt8, math.MinInt8}) })
	t.Run("int16", func(t *testing.T) { testData(t, bucketSortFn[[]int16](), []int16{math.MaxInt16, math.MinInt16}) })
	t.Run("int32", func(t *testing.T) { testData(t, bucketSortFn[[]int32](), []int32{math.MaxInt32, math.MinInt32}) })
	t.Run("int64", func(t *testing.T) { testData(t, bucketSortFn[[]int64](), []int64{math.MaxInt64, math.MinInt64}) })
	t.Run("uint", func(t *testing.T) { testData(t, bucketSortFn[[]uint](), []uint{math.MaxUint, 0}) })
	t.Run("uint8", func(t *testing.T) { testData(t, bucketSortFn[[]uint8](), []uint8{math.MaxUint8, 0}) })
	t.Run("uint16", func(t *testing.T) { testData(t, bucketSortFn[[]uint16](), []uint16{math.MaxUint16, 0}) })
	t.Run("uint32", func(t *testing.T) { testData(t, bucketSortFn[[]uint32](), []uint32{math.MaxUint32, 0}) })
	t.Run("uint64", func(t *testing.T) { testData(t, bucketSortFn[[]uint64](), []uint64{math.MaxUint64, 0}) })
}

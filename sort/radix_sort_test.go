package sort_test

import (
	"math"
	"testing"

	"github.com/denpeshkov/go-grind/sort"
)

func TestRadixSort_EmptyNil(t *testing.T) {
	testEmptyNilSlice(t, sort.RadixSort[[]uint])
}

func TestRadixSort_RandomInts(t *testing.T) {
	t.Run("uint", func(t *testing.T) { testRandomInts(t, sort.RadixSort[[]uint]) })
	t.Run("uint8", func(t *testing.T) { testRandomInts(t, sort.RadixSort[[]uint8]) })
	t.Run("uint16", func(t *testing.T) { testRandomInts(t, sort.RadixSort[[]uint16]) })
	t.Run("uint32", func(t *testing.T) { testRandomInts(t, sort.RadixSort[[]uint32]) })
	t.Run("uint64", func(t *testing.T) { testRandomInts(t, sort.RadixSort[[]uint64]) })
	t.Run("uintptr", func(t *testing.T) { testRandomInts(t, sort.RadixSort[[]uintptr]) })
}

func TestRadixSort_Extremums(t *testing.T) {
	t.Run("uint", func(t *testing.T) { testData(t, sort.RadixSort[[]uint], []uint{math.MaxUint, 0}) })
	t.Run("uint8", func(t *testing.T) { testData(t, sort.RadixSort[[]uint8], []uint8{math.MaxUint8, 0}) })
	t.Run("uint16", func(t *testing.T) { testData(t, sort.RadixSort[[]uint16], []uint16{math.MaxUint16, 0}) })
	t.Run("uint32", func(t *testing.T) { testData(t, sort.RadixSort[[]uint32], []uint32{math.MaxUint32, 0}) })
	t.Run("uint64", func(t *testing.T) { testData(t, sort.RadixSort[[]uint64], []uint64{math.MaxUint64, 0}) })
}

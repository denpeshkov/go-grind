package sort_test

import (
	"testing"

	"github.com/denpeshkov/go-grind/sort"
)

func TestQuickSort_EmptyNil(t *testing.T) {
	testEmptyNilSlice(t, sort.QuickSort[[]int])
}

func TestQuickSort_Data(t *testing.T) {
	testData(t, sort.QuickSort, data)
}

func TestQuickSort_RandomInts(t *testing.T) {
	testRandomInts(t, sort.QuickSort[[]int])
}

package sort_test

import (
	"testing"

	"github.com/denpeshkov/go-grind/sort"
)

func TestHeapSort_EmptyNil(t *testing.T) {
	testEmptyNilSlice(t, sort.HeapSort[[]int])
}

func TestHeapSort_Data(t *testing.T) {
	testData(t, sort.HeapSort, data)
}

func TestHeapSort_RandomInts(t *testing.T) {
	testRandomInts(t, sort.HeapSort[[]int])
}

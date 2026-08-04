package sort_test

import (
	"testing"

	"github.com/denpeshkov/go-grind/sort"
)

func TestMergeSort_EmptyNil(t *testing.T) {
	testEmptyNilSlice(t, sort.MergeSort[[]int])
}

func TestMergeSort_Data(t *testing.T) {
	testData(t, sort.MergeSort, data)
}

func TestMergeSort_RandomInts(t *testing.T) {
	testRandomInts(t, sort.MergeSort[[]int])
}

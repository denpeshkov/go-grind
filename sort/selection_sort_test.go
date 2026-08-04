package sort_test

import (
	"testing"

	"github.com/denpeshkov/go-grind/sort"
)

func TestSelectionSort_EmptyNil(t *testing.T) {
	testEmptyNilSlice(t, sort.SelectionSort[[]int])
}

func TestSelectionSort_Data(t *testing.T) {
	testData(t, sort.SelectionSort, data)
}

func TestSelectionSort_RandomInts(t *testing.T) {
	testRandomInts(t, sort.SelectionSort[[]int])
}

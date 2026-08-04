package sort_test

import (
	"testing"

	"github.com/denpeshkov/go-grind/sort"
)

func TestInsertionSort_EmptyNil(t *testing.T) {
	testEmptyNilSlice(t, sort.InsertionSort[[]int])
}

func TestInsertionSort_Data(t *testing.T) {
	testData(t, sort.InsertionSort, data)
}

func TestInsertionSort_RandomInts(t *testing.T) {
	testRandomInts(t, sort.InsertionSort[[]int])
}

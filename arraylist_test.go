package collections

import (
	"cmp"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArrayList_BasicCRUD(t *testing.T) {
	l := NewArrayList[int]()
	assert.True(t, l.IsEmpty())
	assert.Equal(t, 0, l.Size())
	l.Add(1)
	l.AddAll(2, 3)
	l.AddSeq(seqOf([]int{4, 5}))
	assert.Equal(t, 5, l.Size())
	v, ok := l.Get(0)
	require.True(t, ok)
	assert.Equal(t, 1, v)
	old, ok := l.Set(1, 20)
	require.True(t, ok)
	assert.Equal(t, 2, old)
	v, ok = l.First()
	require.True(t, ok)
	assert.Equal(t, 1, v)
	v, ok = l.Last()
	require.True(t, ok)
	assert.Equal(t, 5, v)
	require.True(t, l.Insert(0, 0))
	require.True(t, l.InsertAll(2, 7, 8))
	removed, ok := l.RemoveAt(2) // remove 7
	require.True(t, ok)
	assert.Equal(t, 7, removed)
	assert.True(t, l.Remove(20, eqV[int]))
	rf := l.RemoveFunc(func(x int) bool { return x%2 == 0 })
	assert.Greater(t, rf, 0, "RemoveFunc expected to remove evens")
	l = NewArrayListFrom(1, 2, 3, 4, 5)
	rr := l.RetainFunc(func(x int) bool { return x > 3 })
	assert.Equal(t, 3, rr)
	assert.Equal(t, 2, l.Size())
}

func TestArrayList_SearchViewSortAggregate(t *testing.T) {
	l := NewArrayListFrom(1, 2, 3, 2, 1)
	assert.Equal(t, 1, l.IndexOf(2, eqV[int]))
	assert.Equal(t, 3, l.LastIndexOf(2, eqV[int]))
	assert.True(t, l.Contains(3, eqV[int]))
	v, ok := l.Find(func(x int) bool { return x%2 == 0 })
	require.True(t, ok)
	assert.Equal(t, 0, v%2)
	assert.Equal(t, 2, l.FindIndex(func(x int) bool { return x == 3 }))
	sub := l.SubList(1, 4).ToSlice()
	assert.True(t, slices.Equal(sub, []int{2, 3, 2}), "SubList=%v", sub)
	// Sort ascending for test
	l.Sort(func(a, b int) int { return cmp.Compare(a, b) })
	got := l.ToSlice()
	assert.True(t, slices.IsSortedFunc(got, func(a, b int) int { return cmp.Compare(a, b) }), "Sort ascending failed: %v", got)
	assert.True(t, l.Any(func(x int) bool { return x == 1 }))
	assert.True(t, l.Every(func(x int) bool { return x >= 1 }))
	// Seq and Reversed consistencies
	seqVals := make([]int, 0, l.Size())
	for v := range l.Seq() {
		seqVals = append(seqVals, v)
	}
	revVals := make([]int, 0, l.Size())
	for v := range l.Reversed() {
		revVals = append(revVals, v)
	}
	revCopy := slices.Clone(revVals)
	slices.Reverse(revCopy)
	assert.Equal(t, len(seqVals), len(revVals))
	assert.True(t, slices.Equal(seqVals, revCopy), "Seq vs Reversed mismatch: %v vs %v", seqVals, revVals)
}

func TestArrayList_InsertAtEnd(t *testing.T) {
	l := NewArrayListFrom(1, 2, 3)
	require.True(t, l.Insert(l.Size(), 4))
	got := l.ToSlice()
	want := []int{1, 2, 3, 4}
	assert.Equal(t, len(want), len(got))
	for i, v := range want {
		assert.Equal(t, v, got[i], "Values should match, got=%v want=%v", got, want)
	}
}

func TestArrayList_InsertOutOfBounds(t *testing.T) {
	l := NewArrayListFrom(1, 2, 3)
	assert.False(t, l.Insert(l.Size()+1, 9))
}

func TestArrayList_NegativeIndex(t *testing.T) {
	l := NewArrayList[int]()
	assert.False(t, l.Insert(-1, 1))
	_, ok := l.Set(-1, 0)
	assert.False(t, ok)
	_, ok = l.Get(-1)
	assert.False(t, ok)
}

func TestArrayList_RemoveFromEmpty(t *testing.T) {
	l := NewArrayList[int]()
	_, ok := l.RemoveAt(0)
	assert.False(t, ok)
	_, ok = l.RemoveFirst()
	require.False(t, ok)
	_, ok = l.RemoveLast()
	require.False(t, ok)
}

func TestArrayList_SubListBoundaries(t *testing.T) {
	l := NewArrayListFrom(1, 2, 3)
	// invalid ranges return empty list
	assert.Equal(t, 0, l.SubList(-1, 2).Size(), "Negative from should return empty list")
	assert.Equal(t, 0, l.SubList(0, 4).Size(), "To beyond size should return empty list")
	assert.Equal(t, 0, l.SubList(3, 2).Size(), "From > to should return empty list")
	// valid boundary [0,0) -> empty
	assert.Equal(t, 0, l.SubList(0, 0).Size(), "[0,0) should be empty")
	// full
	assert.Equal(t, 3, l.SubList(0, 3).Size(), "Full slice size mismatch")
}

func TestArrayList_Clear(t *testing.T) {
	l := NewArrayListFrom(1, 2, 3, 4, 5)
	assert.False(t, l.IsEmpty())
	l.Clear()
	assert.True(t, l.IsEmpty())
	assert.Equal(t, 0, l.Size())
}

func TestArrayList_ForEach(t *testing.T) {
	l := NewArrayListFrom(1, 2, 3, 4, 5)
	sum := 0
	l.ForEach(func(v int) bool {
		sum += v
		return true
	})
	assert.Equal(t, 15, sum)
	// Test early termination
	count := 0
	l.ForEach(func(v int) bool {
		count++
		return count < 3
	})
	assert.Equal(t, 3, count)
}

func TestArrayList_Clone(t *testing.T) {
	l := NewArrayListFrom(1, 2, 3)
	clone := l.Clone()
	clone.Add(4)
	assert.Equal(t, 3, l.Size())
	assert.Equal(t, 4, clone.Size())
}

func TestArrayList_Filter(t *testing.T) {
	l := NewArrayListFrom(1, 2, 3, 4, 5, 6)
	evens := l.Filter(func(v int) bool { return v%2 == 0 })
	slice := evens.ToSlice()
	assert.Equal(t, 3, len(slice))
	for _, v := range slice {
		assert.Equal(t, 0, v%2)
	}
}

func TestArrayList_NewWithCapacityAndString(t *testing.T) {
	l := NewArrayListWithCapacity[int](10)
	assert.True(t, l.IsEmpty())
	l.Add(1)
	assert.Equal(t, 1, l.Size())

	// String
	str := l.String()
	assert.Contains(t, str, "arrayList")
	assert.Contains(t, str, "1")
}

package collections

import (
	"slices"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func intEq(a, b int) bool { return a == b }

func TestCOWList_Basic(t *testing.T) {
	l := NewCOWList[int]()

	assert.True(t, l.IsEmpty())
	assert.Equal(t, 0, l.Size())

	l.Add(1)
	l.Add(2)
	l.Add(3)

	assert.Equal(t, 3, l.Size())

	v, ok := l.Get(1)
	require.True(t, ok)
	assert.Equal(t, 2, v)

	v, ok = l.First()
	require.True(t, ok)
	assert.Equal(t, 1, v)

	v, ok = l.Last()
	require.True(t, ok)
	assert.Equal(t, 3, v)
}

func TestCOWList_Set(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3)

	old, ok := l.Set(1, 20)
	require.True(t, ok)
	assert.Equal(t, 2, old)

	v, _ := l.Get(1)
	assert.Equal(t, 20, v)

	// Out of bounds
	_, ok = l.Set(10, 100)
	assert.False(t, ok)
}

func TestCOWList_Insert(t *testing.T) {
	l := NewCOWListFrom(1, 3)

	require.True(t, l.Insert(1, 2))

	expected := []int{1, 2, 3}
	assert.True(t, slices.Equal(l.ToSlice(), expected))

	// Insert at beginning
	l.Insert(0, 0)
	v, _ := l.First()
	assert.Equal(t, 0, v)

	// Insert at end
	l.Insert(l.Size(), 4)
	v, _ = l.Last()
	assert.Equal(t, 4, v)
}

func TestCOWList_InsertAll(t *testing.T) {
	l := NewCOWListFrom(1, 5)

	require.True(t, l.InsertAll(1, 2, 3, 4))

	expected := []int{1, 2, 3, 4, 5}
	assert.True(t, slices.Equal(l.ToSlice(), expected))
}

func TestCOWList_Remove(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3, 2, 4)

	// Remove first occurrence
	assert.True(t, l.Remove(2, intEq))

	expected := []int{1, 3, 2, 4}
	assert.True(t, slices.Equal(l.ToSlice(), expected))

	// Remove non-existent
	assert.False(t, l.Remove(99, intEq))
}

func TestCOWList_RemoveAt(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3)

	v, ok := l.RemoveAt(1)
	require.True(t, ok)
	assert.Equal(t, 2, v)

	expected := []int{1, 3}
	assert.True(t, slices.Equal(l.ToSlice(), expected))
}

func TestCOWList_RemoveFirstLast(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3)

	v, ok := l.RemoveFirst()
	require.True(t, ok)
	assert.Equal(t, 1, v)

	v, ok = l.RemoveLast()
	require.True(t, ok)
	assert.Equal(t, 3, v)

	assert.Equal(t, 1, l.Size())
}

func TestCOWList_RemoveFunc(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3, 4, 5, 6)

	removed := l.RemoveFunc(func(e int) bool { return e%2 == 0 })
	assert.Equal(t, 3, removed)

	expected := []int{1, 3, 5}
	assert.True(t, slices.Equal(l.ToSlice(), expected))
}

func TestCOWList_RetainFunc(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3, 4, 5, 6)

	removed := l.RetainFunc(func(e int) bool { return e%2 == 0 })
	assert.Equal(t, 3, removed)

	expected := []int{2, 4, 6}
	assert.True(t, slices.Equal(l.ToSlice(), expected))
}

func TestCOWList_IndexOf(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3, 2, 4)

	assert.Equal(t, 1, l.IndexOf(2, intEq))
	assert.Equal(t, 3, l.LastIndexOf(2, intEq))
	assert.Equal(t, -1, l.IndexOf(99, intEq))
}

func TestCOWList_Contains(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3)

	assert.True(t, l.Contains(2, intEq))
	assert.False(t, l.Contains(99, intEq))
}

func TestCOWList_Find(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3, 4, 5)

	v, ok := l.Find(func(e int) bool { return e > 3 })
	require.True(t, ok)
	assert.Equal(t, 4, v)

	idx := l.FindIndex(func(e int) bool { return e > 3 })
	assert.Equal(t, 3, idx)
}

func TestCOWList_SubList(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3, 4, 5)

	sub := l.SubList(1, 4)
	expected := []int{2, 3, 4}
	assert.True(t, slices.Equal(sub.ToSlice(), expected))

	// Invalid range
	empty := l.SubList(3, 1)
	assert.True(t, empty.IsEmpty())
}

func TestCOWList_Reversed(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3)

	var result []int
	for v := range l.Reversed() {
		result = append(result, v)
	}

	expected := []int{3, 2, 1}
	assert.True(t, slices.Equal(result, expected))
}

func TestCOWList_Clone(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3)
	clone := l.Clone()

	l.Add(4)
	assert.Equal(t, 3, clone.Size(), "Clone should not be affected by original modifications")
}

func TestCOWList_Filter(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3, 4, 5, 6)

	filtered := l.Filter(func(e int) bool { return e%2 == 0 })
	expected := []int{2, 4, 6}
	assert.True(t, slices.Equal(filtered.ToSlice(), expected))
}

func TestCOWList_Sort(t *testing.T) {
	l := NewCOWListFrom(3, 1, 4, 1, 5, 9, 2, 6)

	l.Sort(func(a, b int) int { return a - b })
	expected := []int{1, 1, 2, 3, 4, 5, 6, 9}
	assert.True(t, slices.Equal(l.ToSlice(), expected))
}

func TestCOWList_AnyEvery(t *testing.T) {
	l := NewCOWListFrom(2, 4, 6, 8)

	assert.True(t, l.Every(func(e int) bool { return e%2 == 0 }))
	assert.False(t, l.Any(func(e int) bool { return e%2 != 0 }))

	l.Add(3)
	assert.True(t, l.Any(func(e int) bool { return e%2 != 0 }))
}

func TestCOWList_Clear(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3)
	l.Clear()

	assert.True(t, l.IsEmpty())
}

func TestCOWList_AddAll(t *testing.T) {
	l := NewCOWList[int]()
	l.AddAll(1, 2, 3)

	expected := []int{1, 2, 3}
	assert.True(t, slices.Equal(l.ToSlice(), expected))
}

func TestCOWList_Concurrent(t *testing.T) {
	l := NewCOWList[int]()
	const goroutines = 100
	const opsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 2) // readers and writers

	// Writers
	for i := range goroutines {
		go func(base int) {
			defer wg.Done()
			for j := range opsPerGoroutine {
				l.Add(base*opsPerGoroutine + j)
			}
		}(i)
	}

	// Readers
	for range goroutines {
		go func() {
			defer wg.Done()
			for range opsPerGoroutine {
				_ = l.Size()
				_ = l.ToSlice()
				_, _ = l.First()
				_, _ = l.Last()
			}
		}()
	}

	wg.Wait()

	expectedSize := goroutines * opsPerGoroutine
	assert.Equal(t, expectedSize, l.Size())
}

func TestCOWList_SnapshotConsistency(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3, 4, 5)

	// Take a snapshot via iteration
	var snapshot []int
	for v := range l.Seq() {
		snapshot = append(snapshot, v)
		if v == 2 {
			// Modify during iteration
			l.Add(6)
		}
	}

	// Snapshot should be consistent (not see the modification)
	expected := []int{1, 2, 3, 4, 5}
	assert.True(t, slices.Equal(snapshot, expected), "Snapshot should be %v, got %v", expected, snapshot)

	// But the list should have the new element
	assert.Equal(t, 6, l.Size())
}

func TestCOWList_String(t *testing.T) {
	l := NewCOWListFrom(1, 2, 3)
	s := l.String()

	assert.Equal(t, "cowList{1, 2, 3}", s)
}

func BenchmarkCOWList_Read(b *testing.B) {
	l := NewCOWList[int]()
	for i := range 1000 {
		l.Add(i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = l.Get(500)
		}
	})
}

func BenchmarkCOWList_Write(b *testing.B) {
	l := NewCOWList[int]()

	b.ResetTimer()
	for b.Loop() {
		l.Add(0)
	}
}

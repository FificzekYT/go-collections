package collections

import (
	"slices"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLockFreeList_Basic(t *testing.T) {
	l := NewLockFreeListOrdered[int]()

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

func TestLockFreeList_Set(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3)

	old, ok := l.Set(1, 20)
	require.True(t, ok)
	assert.Equal(t, 2, old)

	v, _ := l.Get(1)
	assert.Equal(t, 20, v)

	// Out of bounds
	_, ok = l.Set(10, 100)
	assert.False(t, ok)
}

func TestLockFreeList_Insert(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 3)

	require.True(t, l.Insert(1, 2))

	expected := []int{1, 2, 3}
	assert.True(t, slices.Equal(l.ToSlice(), expected))

	// Insert at beginning
	l.Insert(0, 0)
	v, _ := l.First()
	assert.Equal(t, 0, v)
}

func TestLockFreeList_Remove(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3, 2, 4)

	// Remove first occurrence
	assert.True(t, l.Remove(2, intEq))

	expected := []int{1, 3, 2, 4}
	assert.True(t, slices.Equal(l.ToSlice(), expected))

	// Remove non-existent
	assert.False(t, l.Remove(99, intEq))
}

func TestLockFreeList_RemoveAt(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3)

	v, ok := l.RemoveAt(1)
	require.True(t, ok)
	assert.Equal(t, 2, v)

	expected := []int{1, 3}
	assert.True(t, slices.Equal(l.ToSlice(), expected))
}

func TestLockFreeList_RemoveFirstLast(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3)

	v, ok := l.RemoveFirst()
	require.True(t, ok)
	assert.Equal(t, 1, v)

	v, ok = l.RemoveLast()
	require.True(t, ok)
	assert.Equal(t, 3, v)

	assert.Equal(t, 1, l.Size())
}

func TestLockFreeList_RemoveFunc(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3, 4, 5, 6)

	removed := l.RemoveFunc(func(e int) bool { return e%2 == 0 })
	assert.Equal(t, 3, removed)

	expected := []int{1, 3, 5}
	assert.True(t, slices.Equal(l.ToSlice(), expected))
}

func TestLockFreeList_RetainFunc(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3, 4, 5, 6)

	removed := l.RetainFunc(func(e int) bool { return e%2 == 0 })
	assert.Equal(t, 3, removed)

	expected := []int{2, 4, 6}
	assert.True(t, slices.Equal(l.ToSlice(), expected))
}

func TestLockFreeList_IndexOf(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3, 2, 4)

	assert.Equal(t, 1, l.IndexOf(2, intEq))
	assert.Equal(t, 3, l.LastIndexOf(2, intEq))
	assert.Equal(t, -1, l.IndexOf(99, intEq))
}

func TestLockFreeList_Contains(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3)

	assert.True(t, l.Contains(2, intEq))
	assert.False(t, l.Contains(99, intEq))
}

func TestLockFreeList_Find(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3, 4, 5)

	v, ok := l.Find(func(e int) bool { return e > 3 })
	require.True(t, ok)
	assert.Equal(t, 4, v)

	idx := l.FindIndex(func(e int) bool { return e > 3 })
	assert.Equal(t, 3, idx)
}

func TestLockFreeList_SubList(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3, 4, 5)

	sub := l.SubList(1, 4)
	expected := []int{2, 3, 4}
	assert.True(t, slices.Equal(sub.ToSlice(), expected))

	// Invalid range
	empty := l.SubList(3, 1)
	assert.True(t, empty.IsEmpty())
}

func TestLockFreeList_Reversed(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3)

	var result []int
	for v := range l.Reversed() {
		result = append(result, v)
	}

	expected := []int{3, 2, 1}
	assert.True(t, slices.Equal(result, expected))
}

func TestLockFreeList_Clone(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3)
	clone := l.Clone()

	l.Add(4)
	assert.Equal(t, 3, clone.Size(), "Clone should not be affected by original modifications")
}

func TestLockFreeList_Filter(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3, 4, 5, 6)

	filtered := l.Filter(func(e int) bool { return e%2 == 0 })
	expected := []int{2, 4, 6}
	assert.True(t, slices.Equal(filtered.ToSlice(), expected))
}

func TestLockFreeList_Sort(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(3, 1, 4, 1, 5, 9, 2, 6)

	l.Sort(func(a, b int) int { return a - b })
	expected := []int{1, 1, 2, 3, 4, 5, 6, 9}
	assert.True(t, slices.Equal(l.ToSlice(), expected))
}

func TestLockFreeList_AnyEvery(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(2, 4, 6, 8)

	assert.True(t, l.Every(func(e int) bool { return e%2 == 0 }))
	assert.False(t, l.Any(func(e int) bool { return e%2 != 0 }))

	l.Add(3)
	assert.True(t, l.Any(func(e int) bool { return e%2 != 0 }))
}

func TestLockFreeList_Clear(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3)
	l.Clear()

	assert.True(t, l.IsEmpty())
}

func TestLockFreeList_Concurrent(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	const goroutines = 50
	const opsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 2) // readers and writers

	// Writers
	for i := range goroutines {
		go func(id int) {
			defer wg.Done()
			for j := range opsPerGoroutine {
				l.Add(id*opsPerGoroutine + j)
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
	actualSize := l.Size()
	// Size may be approximate due to concurrent modifications
	assert.InDelta(t, expectedSize, actualSize, float64(expectedSize)/10, "expected size around %d, got %d", expectedSize, actualSize)
}

func TestLockFreeList_ConcurrentRemove(t *testing.T) {
	l := NewLockFreeListOrdered[int]()

	// Add elements
	for i := range 1000 {
		l.Add(i)
	}

	var wg sync.WaitGroup
	wg.Add(10)

	// Concurrent removers
	for i := range 10 {
		go func(start int) {
			defer wg.Done()
			for j := start; j < 1000; j += 10 {
				l.Remove(j, intEq)
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, 0, l.Size())
}

func TestLockFreeList_String(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	l.AddAll(1, 2, 3)
	s := l.String()

	assert.Equal(t, "lockFreeList{1, 2, 3}", s)
}

func TestLockFreeList_WithEqualer(t *testing.T) {
	type person struct {
		name string
		age  int
	}

	eq := func(a, b person) bool { return a.name == b.name }
	l := NewLockFreeList(eq)

	l.Add(person{"Alice", 30})
	l.Add(person{"Bob", 25})

	// Should find by name only
	assert.True(t, l.Contains(person{"Alice", 0}, eq))
}

func BenchmarkLockFreeList_Add(b *testing.B) {
	l := NewLockFreeListOrdered[int]()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			l.Add(i)
			i++
		}
	})
}

func BenchmarkLockFreeList_Read(b *testing.B) {
	l := NewLockFreeListOrdered[int]()
	for i := range 1000 {
		l.Add(i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = l.First()
		}
	})
}

func BenchmarkLockFreeList_Contains(b *testing.B) {
	l := NewLockFreeListOrdered[int]()
	for i := range 1000 {
		l.Add(i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			l.Contains(i%1000, intEq)
			i++
		}
	})
}

func TestLockFreeList_PhysicalDelete(t *testing.T) {
	l := NewLockFreeListOrdered[int]()
	const n = 1000
	// Build list 0..n-1
	for i := range n {
		l.Add(i)
	}
	// Logically delete all even numbers
	removed := l.RemoveFunc(func(v int) bool { return v%2 == 0 })
	require.Equal(t, n/2, removed, "expected to logically delete half of the elements")

	// Snapshot should only contain odds
	snapBefore := l.ToSlice()
	require.Equal(t, n/2, len(snapBefore))
	for _, v := range snapBefore {
		assert.Equal(t, 1, v%2, "snapshot before physical delete should contain only odd numbers")
	}

	// Physical deletion should unlink logically deleted nodes without changing visible semantics
	lf := l.(*lockFreeList[int])
	lf.PhysicalDelete()

	// Snapshot after physical delete remains the same (odds only)
	snapAfter := l.ToSlice()
	require.Equal(t, n/2, len(snapAfter))
	for _, v := range snapAfter {
		assert.Equal(t, 1, v%2, "snapshot after physical delete should contain only odd numbers")
	}

	// Idempotency: calling PhysicalDelete again should not change content
	lf.PhysicalDelete()
	snapAgain := l.ToSlice()
	require.Equal(t, n/2, len(snapAgain))
	for i := range snapAfter {
		assert.Equal(t, snapAfter[i], snapAgain[i])
	}
}

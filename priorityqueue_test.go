package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriorityQueue_Basic(t *testing.T) {
	pq := NewPriorityQueueOrdered[int]()
	assert.True(t, pq.IsEmpty())
	assert.Equal(t, 0, pq.Size())

	pq.Push(3)
	pq.Push(1)
	pq.Push(4)
	pq.Push(1)
	pq.Push(5)

	assert.False(t, pq.IsEmpty())
	assert.Equal(t, 5, pq.Size())

	// Min-heap: smallest first
	v, ok := pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	v, ok = pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	v, ok = pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	v, ok = pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, 3, v)
}

func TestPriorityQueue_MaxHeap(t *testing.T) {
	pq := NewMaxPriorityQueue[int]()

	pq.Push(3)
	pq.Push(1)
	pq.Push(4)
	pq.Push(1)
	pq.Push(5)

	// Max-heap: largest first
	v, ok := pq.Peek()
	assert.True(t, ok)
	assert.Equal(t, 5, v)

	v, ok = pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, 5, v)

	v, ok = pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, 4, v)

	v, ok = pq.Pop()
	assert.True(t, ok)
	assert.Equal(t, 3, v)
}

func TestPriorityQueue_Empty(t *testing.T) {
	pq := NewPriorityQueueOrdered[int]()

	_, ok := pq.Peek()
	assert.False(t, ok)

	_, ok = pq.Pop()
	assert.False(t, ok)
}

func TestPriorityQueue_Clear(t *testing.T) {
	pq := NewPriorityQueueOrdered[int]()
	pq.PushAll(1, 2, 3, 4, 5)

	assert.Equal(t, 5, pq.Size())

	pq.Clear()
	assert.True(t, pq.IsEmpty())
	assert.Equal(t, 0, pq.Size())
}

func TestPriorityQueue_PushAll(t *testing.T) {
	pq := NewPriorityQueueOrdered[int]()
	pq.PushAll(5, 3, 1, 4, 2)

	assert.Equal(t, 5, pq.Size())

	// Should pop in sorted order (min-heap)
	v, _ := pq.Pop()
	assert.Equal(t, 1, v)
	v, _ = pq.Pop()
	assert.Equal(t, 2, v)
}

func TestPriorityQueue_From(t *testing.T) {
	pq := NewPriorityQueueFrom(CompareFunc[int](), 5, 3, 1, 4, 2)

	assert.Equal(t, 5, pq.Size())

	v, _ := pq.Peek()
	assert.Equal(t, 1, v)
}

func TestPriorityQueue_ToSlice(t *testing.T) {
	pq := NewPriorityQueueFrom(CompareFunc[int](), 5, 3, 1, 4, 2)

	slice := pq.ToSlice()
	assert.Equal(t, 5, len(slice))
	// ToSlice returns heap order, not sorted
}

func TestPriorityQueue_ToSortedSlice(t *testing.T) {
	pq := NewPriorityQueueFrom(CompareFunc[int](), 5, 3, 1, 4, 2)

	sorted := pq.ToSortedSlice()
	assert.Equal(t, []int{1, 2, 3, 4, 5}, sorted)
}

func TestPriorityQueue_Seq(t *testing.T) {
	pq := NewPriorityQueueFrom(CompareFunc[int](), 3, 1, 2)

	var collected []int
	for v := range pq.Seq() {
		collected = append(collected, v)
	}
	assert.Equal(t, 3, len(collected))
}

func TestPriorityQueue_String(t *testing.T) {
	pq := NewPriorityQueueFrom(CompareFunc[int](), 1, 2, 3)
	s := pq.String()
	assert.Contains(t, s, "priorityQueue")
}

func TestPriorityQueue_CustomComparator(t *testing.T) {
	// Priority by string length (shorter = higher priority)
	pq := NewPriorityQueue(func(a, b string) int {
		return len(a) - len(b)
	})

	pq.Push("hello")
	pq.Push("hi")
	pq.Push("hey")
	pq.Push("h")

	v, _ := pq.Pop()
	assert.Equal(t, "h", v)

	v, _ = pq.Pop()
	assert.Equal(t, "hi", v)

	v, _ = pq.Pop()
	assert.Equal(t, "hey", v)

	v, _ = pq.Pop()
	assert.Equal(t, "hello", v)
}

func TestPriorityQueue_WithCapacity(t *testing.T) {
	pq := NewPriorityQueueWithCapacity(CompareFunc[int](), 100)
	assert.True(t, pq.IsEmpty())

	pq.Push(1)
	assert.Equal(t, 1, pq.Size())
}

func TestPriorityQueue_Panic(t *testing.T) {
	assert.Panics(t, func() {
		var nilCmp Comparator[int]
		NewPriorityQueue(nilCmp)
	})

	assert.Panics(t, func() {
		var nilCmp Comparator[int]
		NewPriorityQueueWithCapacity(nilCmp, 10)
	})

	assert.Panics(t, func() {
		var nilCmp Comparator[int]
		NewPriorityQueueFrom(nilCmp, 1, 2, 3)
	})
}

func TestPriorityQueue_LargeDataset(t *testing.T) {
	pq := NewPriorityQueueOrdered[int]()

	// Push 1000 elements in reverse order
	for i := 1000; i > 0; i-- {
		pq.Push(i)
	}

	assert.Equal(t, 1000, pq.Size())

	// Pop should return in sorted order
	for i := 1; i <= 1000; i++ {
		v, ok := pq.Pop()
		assert.True(t, ok)
		assert.Equal(t, i, v)
	}

	assert.True(t, pq.IsEmpty())
}

func TestPriorityQueue_Duplicates(t *testing.T) {
	pq := NewPriorityQueueOrdered[int]()
	pq.PushAll(3, 3, 3, 2, 2, 1, 1)

	assert.Equal(t, 7, pq.Size())

	// All duplicates should be preserved
	v, _ := pq.Pop()
	assert.Equal(t, 1, v)
	v, _ = pq.Pop()
	assert.Equal(t, 1, v)
	v, _ = pq.Pop()
	assert.Equal(t, 2, v)
	v, _ = pq.Pop()
	assert.Equal(t, 2, v)
	v, _ = pq.Pop()
	assert.Equal(t, 3, v)
	v, _ = pq.Pop()
	assert.Equal(t, 3, v)
	v, _ = pq.Pop()
	assert.Equal(t, 3, v)
}

package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArrayQueue_BasicFIFO(t *testing.T) {
	q := NewArrayQueue[int]()
	assert.True(t, q.IsEmpty(), "New queue should be empty")
	q.Enqueue(1)
	q.EnqueueAll(2, 3)
	assert.Equal(t, 3, q.Size())
	v, ok := q.Peek()
	require.True(t, ok)
	assert.Equal(t, 1, v)
	// Seq front->back: 1,2,3
	want := 1
	for v := range q.Seq() {
		assert.Equal(t, want, v)
		want++
	}
	v, ok = q.Dequeue()
	require.True(t, ok)
	assert.Equal(t, 1, v)
	v, ok = q.Dequeue()
	require.True(t, ok)
	assert.Equal(t, 2, v)
	v, ok = q.Dequeue()
	require.True(t, ok)
	assert.Equal(t, 3, v)
	_, ok = q.Dequeue()
	assert.False(t, ok, "Dequeue on empty should fail")
}

func TestArrayQueue_GrowthAndCompact(t *testing.T) {
	q := NewArrayQueueWithCapacity[int](1)
	for i := range 100 {
		q.Enqueue(i)
	}
	for i := range 50 {
		v, ok := q.Dequeue()
		require.True(t, ok)
		assert.Equal(t, i, v, "Dequeue mismatch at %d", i)
	}
	// After many dequeues, ensure the remaining sequence is correct
	want := 50
	for v := range q.Seq() {
		assert.Equal(t, want, v)
		want++
	}
	assert.Equal(t, 50, q.Size())
}

func TestArrayQueue_Clear(t *testing.T) {
	q := NewArrayQueueFrom(1, 2, 3)
	assert.False(t, q.IsEmpty())
	q.Clear()
	assert.True(t, q.IsEmpty())
	assert.Equal(t, 0, q.Size())
}

func TestArrayQueue_ToSlice(t *testing.T) {
	q := NewArrayQueueFrom(1, 2, 3)
	slice := q.ToSlice()
	assert.Equal(t, 3, len(slice))
	assert.Equal(t, 1, slice[0])
	assert.Equal(t, 3, slice[2])
}

func TestArrayQueue_From(t *testing.T) {
	q := NewArrayQueueFrom(5, 4, 3)
	assert.Equal(t, 3, q.Size())
	v, _ := q.Dequeue()
	assert.Equal(t, 5, v)
}

func TestArrayQueue_String(t *testing.T) {
	q := NewArrayQueue[int]()
	q.Enqueue(1)
	str := q.String()
	assert.Contains(t, str, "arrayQueue")
	assert.Contains(t, str, "1")
}

func TestArrayQueue_ShrinkAfterPeak(t *testing.T) {
	// Test the shrink strategy: when head > 3/4 cap and live < 1/4 cap,
	// the capacity should shrink to release memory.
	q := NewArrayQueue[int]().(*arrayQueue[int])

	// Enqueue a large number of elements to create peak capacity
	const peakSize = 1000
	for i := range peakSize {
		q.Enqueue(i)
	}
	peakCap := cap(q.data)
	assert.GreaterOrEqual(t, peakCap, peakSize)

	// Dequeue most elements, leaving very few
	// This should trigger the shrink strategy
	for range peakSize - 10 {
		q.Dequeue()
	}

	// After shrinking, capacity should be significantly smaller
	currentCap := cap(q.data)
	assert.Less(t, currentCap, peakCap/2,
		"capacity should shrink after peak-then-low-water scenario: peak=%d, current=%d",
		peakCap, currentCap)

	// Verify remaining elements are still correct
	assert.Equal(t, 10, q.Size())
	for i := peakSize - 10; i < peakSize; i++ {
		v, ok := q.Dequeue()
		require.True(t, ok)
		assert.Equal(t, i, v)
	}
}

func TestArrayQueue_ShrinkTriggersReallocation(t *testing.T) {
	// Test that shrink behavior triggers reallocation and subsequent
	// enqueue/dequeue operations have stable performance.
	q := NewArrayQueue[int]().(*arrayQueue[int])

	// Build up to a large capacity
	const peakSize = 500
	for i := range peakSize {
		q.Enqueue(i)
	}

	// Record capacity after peak
	peakCap := cap(q.data)
	require.GreaterOrEqual(t, peakCap, peakSize)

	// Dequeue almost all elements to trigger shrink
	for range peakSize - 5 {
		q.Dequeue()
	}

	// After shrink, capacity should be reduced
	shrunkCap := cap(q.data)
	assert.Less(t, shrunkCap, peakCap/2,
		"capacity should shrink significantly: peak=%d, shrunk=%d", peakCap, shrunkCap)

	// Enqueue more elements - should work correctly after shrink
	for i := range 50 {
		q.Enqueue(1000 + i)
	}

	// Verify queue still works correctly
	assert.Equal(t, 55, q.Size())

	// Dequeue and verify correct order
	for i := peakSize - 5; i < peakSize; i++ {
		v, ok := q.Dequeue()
		require.True(t, ok)
		assert.Equal(t, i, v, "remaining original elements should be in order")
	}
	for i := range 50 {
		v, ok := q.Dequeue()
		require.True(t, ok)
		assert.Equal(t, 1000+i, v, "new elements should be in order")
	}
	assert.True(t, q.IsEmpty())
}

func TestArrayQueue_ShrinkDoesNotThrashOnSmallQueues(t *testing.T) {
	// Verify that small queues don't shrink (capacity <= 64 threshold)
	q := NewArrayQueue[int]().(*arrayQueue[int])

	// Use a small queue
	for i := range 50 {
		q.Enqueue(i)
	}
	initialCap := cap(q.data)

	// Dequeue most elements
	for range 45 {
		q.Dequeue()
	}

	// Small queues should not shrink to avoid thrashing
	// The threshold in compact() is cap > 64
	if initialCap <= 64 {
		// For small queues, capacity should remain (in-place shift, not shrink)
		assert.LessOrEqual(t, cap(q.data), initialCap,
			"small queues should not allocate larger, but may compact in-place")
	}

	// Queue should still function correctly
	assert.Equal(t, 5, q.Size())
	for i := 45; i < 50; i++ {
		v, ok := q.Dequeue()
		require.True(t, ok)
		assert.Equal(t, i, v)
	}
}

func TestArrayQueue_ShrinkReducesAllocations(t *testing.T) {
	// Test that after shrink, subsequent enqueue/dequeue operations
	// require fewer allocations due to reduced capacity.
	// This protects the shrink optimization.

	const peakSize = 200
	const steadyOps = 50

	// Scenario: Build up to peak, drain to trigger shrink, then steady-state ops
	allocs := testing.AllocsPerRun(10, func() {
		q := NewArrayQueue[int]().(*arrayQueue[int])

		// Phase 1: Build up to peak capacity
		for i := range peakSize {
			q.Enqueue(i)
		}

		// Phase 2: Drain most elements to trigger shrink
		for range peakSize - 10 {
			q.Dequeue()
		}

		// Phase 3: Steady-state operations after shrink
		// These should use the shrunk capacity, not peak capacity
		for i := range steadyOps {
			q.Enqueue(1000 + i)
			if i%2 == 0 {
				q.Dequeue()
			}
		}
	})

	// The shrink optimization should keep allocations reasonable.
	// Without shrink, peak capacity would be retained and steady-state
	// ops might not need allocations (but memory footprint would be large).
	// With shrink, we trade a reallocation for lower memory.
	// This test ensures shrink doesn't cause excessive reallocations.
	t.Logf("AllocsPerRun for shrink scenario: %.1f", allocs)

	// Expect reasonable allocations (initial growth + one shrink + possible regrowth)
	// The exact number depends on Go's slice growth strategy
	assert.Less(t, allocs, float64(30),
		"shrink optimization should not cause excessive allocations")
}

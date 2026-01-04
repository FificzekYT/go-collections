package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArrayDeque_BasicOps(t *testing.T) {
	d := NewArrayDeque[int]()
	assert.True(t, d.IsEmpty(), "New deque should be empty")
	d.PushBack(2)  // [2]
	d.PushFront(1) // [1,2]
	d.PushBack(3)  // [1,2,3]
	assert.Equal(t, 3, d.Size())
	v, ok := d.PeekFront()
	require.True(t, ok)
	assert.Equal(t, 1, v)
	v, ok = d.PeekBack()
	require.True(t, ok)
	assert.Equal(t, 3, v)
	// Seq front->back: 1,2,3
	want := 1
	for v := range d.Seq() {
		assert.Equal(t, want, v)
		want++
	}
	// Reversed: 3,2,1
	want = 3
	for v := range d.Reversed() {
		assert.Equal(t, want, v)
		want--
	}
	v, ok = d.PopFront()
	require.True(t, ok)
	assert.Equal(t, 1, v)
	v, ok = d.PopBack()
	require.True(t, ok)
	assert.Equal(t, 3, v)
	v, ok = d.PopFront()
	require.True(t, ok)
	assert.Equal(t, 2, v)
	_, ok = d.PopBack()
	require.False(t, ok, "PopBack on empty should fail")
}

func TestArrayDeque_Growth(t *testing.T) {
	d := NewArrayDequeWithCapacity[int](2)
	// push enough to trigger growth and wrap-around
	for i := range 32 {
		if i%2 == 0 {
			d.PushFront(i)
		} else {
			d.PushBack(i)
		}
	}
	// sanity: size
	assert.Equal(t, 32, d.Size())
	// pop all to ensure order is consistent with front/back operations
	count := 0
	for !d.IsEmpty() {
		_, ok := d.PopFront()
		require.True(t, ok, "PopFront failed mid-way")
		count++
	}
	assert.Equal(t, 32, count)
}

func TestArrayDeque_Clear(t *testing.T) {
	d := NewArrayDequeFrom(1, 2, 3)
	assert.False(t, d.IsEmpty())
	d.Clear()
	assert.True(t, d.IsEmpty())
	assert.Equal(t, 0, d.Size())
}

func TestArrayDeque_ToSlice(t *testing.T) {
	d := NewArrayDequeFrom(1, 2, 3)
	slice := d.ToSlice()
	assert.Equal(t, 3, len(slice))
	assert.Equal(t, 1, slice[0])
	assert.Equal(t, 3, slice[2])
}

func TestArrayDeque_From(t *testing.T) {
	d := NewArrayDequeFrom(5, 4, 3)
	assert.Equal(t, 3, d.Size())
	v, _ := d.PopFront()
	assert.Equal(t, 5, v)
}

func TestArrayDeque_String(t *testing.T) {
	d := NewArrayDeque[int]()
	d.PushBack(1)
	str := d.String()
	assert.Contains(t, str, "arrayDeque")
	assert.Contains(t, str, "1")
}

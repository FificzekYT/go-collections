package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArrayStack_BasicAndOrder(t *testing.T) {
	s := NewArrayStack[int]()
	assert.True(t, s.IsEmpty(), "New stack should be empty")
	s.Push(1)
	s.PushAll(2, 3)
	assert.Equal(t, 3, s.Size())
	v, ok := s.Peek()
	require.True(t, ok)
	assert.Equal(t, 3, v)
	// Seq should yield top->bottom: 3,2,1
	i := 3
	for v := range s.Seq() {
		assert.Equal(t, i, v)
		i--
	}
	// ToSlice bottom->top
	ts := s.ToSlice()
	assert.Equal(t, 1, ts[0])
	assert.Equal(t, 3, ts[2])
	v, ok = s.Pop()
	require.True(t, ok)
	assert.Equal(t, 3, v)
	v, ok = s.Pop()
	require.True(t, ok)
	assert.Equal(t, 2, v)
	v, ok = s.Pop()
	require.True(t, ok)
	assert.Equal(t, 1, v)
	_, ok = s.Pop()
	assert.False(t, ok, "Pop on empty should fail")
}

func TestArrayStack_Clear(t *testing.T) {
	s := NewArrayStackFrom(1, 2, 3)
	assert.False(t, s.IsEmpty())
	s.Clear()
	assert.True(t, s.IsEmpty())
	assert.Equal(t, 0, s.Size())
}

func TestArrayStack_WithCapacity(t *testing.T) {
	s := NewArrayStackWithCapacity[int](10)
	assert.True(t, s.IsEmpty())
	s.Push(1)
	assert.Equal(t, 1, s.Size())
}

func TestArrayStack_From(t *testing.T) {
	s := NewArrayStackFrom(1, 2, 3)
	assert.Equal(t, 3, s.Size())
	v, _ := s.Pop()
	assert.Equal(t, 3, v)
}

func TestArrayStack_String(t *testing.T) {
	s := NewArrayStack[int]()
	s.Push(1)
	str := s.String()
	assert.Contains(t, str, "arrayStack")
	assert.Contains(t, str, "1")
}

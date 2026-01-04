package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedList_Basic(t *testing.T) {
	l := NewLinkedList[int]()
	assert.True(t, l.IsEmpty())
	assert.Equal(t, 0, l.Size())

	l.Add(1)
	l.Add(2)
	l.Add(3)

	assert.False(t, l.IsEmpty())
	assert.Equal(t, 3, l.Size())

	v, ok := l.Get(0)
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	v, ok = l.Get(1)
	assert.True(t, ok)
	assert.Equal(t, 2, v)

	v, ok = l.Get(2)
	assert.True(t, ok)
	assert.Equal(t, 3, v)

	_, ok = l.Get(3)
	assert.False(t, ok)

	_, ok = l.Get(-1)
	assert.False(t, ok)
}

func TestLinkedList_AddFirst(t *testing.T) {
	l := NewLinkedList[int]().(*linkedList[int])
	l.AddFirst(3)
	l.AddFirst(2)
	l.AddFirst(1)

	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
}

func TestLinkedList_FirstLast(t *testing.T) {
	l := NewLinkedList[int]()

	_, ok := l.First()
	assert.False(t, ok)

	_, ok = l.Last()
	assert.False(t, ok)

	l.Add(1)
	l.Add(2)
	l.Add(3)

	v, ok := l.First()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	v, ok = l.Last()
	assert.True(t, ok)
	assert.Equal(t, 3, v)
}

func TestLinkedList_Set(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3)

	old, ok := l.Set(1, 20)
	assert.True(t, ok)
	assert.Equal(t, 2, old)

	v, _ := l.Get(1)
	assert.Equal(t, 20, v)

	_, ok = l.Set(10, 100)
	assert.False(t, ok)
}

func TestLinkedList_Insert(t *testing.T) {
	l := NewLinkedListFrom(1, 3)

	ok := l.Insert(1, 2)
	assert.True(t, ok)
	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())

	ok = l.Insert(0, 0)
	assert.True(t, ok)
	assert.Equal(t, []int{0, 1, 2, 3}, l.ToSlice())

	ok = l.Insert(4, 4)
	assert.True(t, ok)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, l.ToSlice())

	ok = l.Insert(-1, -1)
	assert.False(t, ok)

	ok = l.Insert(100, 100)
	assert.False(t, ok)
}

func TestLinkedList_InsertAll(t *testing.T) {
	l := NewLinkedListFrom(1, 5)

	ok := l.InsertAll(1, 2, 3, 4)
	assert.True(t, ok)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, l.ToSlice())
}

func TestLinkedList_RemoveAt(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3, 4, 5)

	v, ok := l.RemoveAt(2)
	assert.True(t, ok)
	assert.Equal(t, 3, v)
	assert.Equal(t, []int{1, 2, 4, 5}, l.ToSlice())

	v, ok = l.RemoveAt(0)
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	assert.Equal(t, []int{2, 4, 5}, l.ToSlice())

	v, ok = l.RemoveAt(2)
	assert.True(t, ok)
	assert.Equal(t, 5, v)
	assert.Equal(t, []int{2, 4}, l.ToSlice())

	_, ok = l.RemoveAt(10)
	assert.False(t, ok)
}

func TestLinkedList_Remove(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3, 2, 1)
	eq := EqualFunc[int]()

	ok := l.Remove(2, eq)
	assert.True(t, ok)
	assert.Equal(t, []int{1, 3, 2, 1}, l.ToSlice())

	ok = l.Remove(100, eq)
	assert.False(t, ok)
}

func TestLinkedList_RemoveFirstLast(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3)

	v, ok := l.RemoveFirst()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	v, ok = l.RemoveLast()
	assert.True(t, ok)
	assert.Equal(t, 3, v)

	assert.Equal(t, []int{2}, l.ToSlice())

	l.Clear()
	_, ok = l.RemoveFirst()
	assert.False(t, ok)

	_, ok = l.RemoveLast()
	assert.False(t, ok)
}

func TestLinkedList_RemoveFunc(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3, 4, 5, 6)

	removed := l.RemoveFunc(func(v int) bool { return v%2 == 0 })
	assert.Equal(t, 3, removed)
	assert.Equal(t, []int{1, 3, 5}, l.ToSlice())
}

func TestLinkedList_RetainFunc(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3, 4, 5, 6)

	removed := l.RetainFunc(func(v int) bool { return v%2 == 0 })
	assert.Equal(t, 3, removed)
	assert.Equal(t, []int{2, 4, 6}, l.ToSlice())
}

func TestLinkedList_IndexOf(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3, 2, 1)
	eq := EqualFunc[int]()

	assert.Equal(t, 1, l.IndexOf(2, eq))
	assert.Equal(t, 3, l.LastIndexOf(2, eq))
	assert.Equal(t, -1, l.IndexOf(100, eq))
}

func TestLinkedList_Contains(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3)
	eq := EqualFunc[int]()

	assert.True(t, l.Contains(2, eq))
	assert.False(t, l.Contains(100, eq))
}

func TestLinkedList_Find(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3, 4, 5)

	v, ok := l.Find(func(n int) bool { return n > 3 })
	assert.True(t, ok)
	assert.Equal(t, 4, v)

	_, ok = l.Find(func(n int) bool { return n > 100 })
	assert.False(t, ok)

	idx := l.FindIndex(func(n int) bool { return n > 3 })
	assert.Equal(t, 3, idx)

	idx = l.FindIndex(func(n int) bool { return n > 100 })
	assert.Equal(t, -1, idx)
}

func TestLinkedList_SubList(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3, 4, 5)

	sub := l.SubList(1, 4)
	assert.Equal(t, []int{2, 3, 4}, sub.ToSlice())

	sub = l.SubList(-1, 3)
	assert.Equal(t, 0, sub.Size())

	sub = l.SubList(3, 2)
	assert.Equal(t, 0, sub.Size())
}

func TestLinkedList_Reversed(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3)

	var reversed []int
	for v := range l.Reversed() {
		reversed = append(reversed, v)
	}
	assert.Equal(t, []int{3, 2, 1}, reversed)
}

func TestLinkedList_Clone(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3)
	clone := l.Clone()

	assert.Equal(t, l.ToSlice(), clone.ToSlice())

	l.Add(4)
	assert.NotEqual(t, l.Size(), clone.Size())
}

func TestLinkedList_Filter(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3, 4, 5, 6)

	evens := l.Filter(func(v int) bool { return v%2 == 0 })
	assert.Equal(t, []int{2, 4, 6}, evens.ToSlice())
}

func TestLinkedList_Sort(t *testing.T) {
	l := NewLinkedListFrom(5, 2, 8, 1, 9, 3)
	l.Sort(CompareFunc[int]())

	assert.Equal(t, []int{1, 2, 3, 5, 8, 9}, l.ToSlice())
}

func TestLinkedList_SortEmpty(t *testing.T) {
	l := NewLinkedList[int]()
	l.Sort(CompareFunc[int]())
	assert.Equal(t, 0, l.Size())
}

func TestLinkedList_SortSingle(t *testing.T) {
	l := NewLinkedListFrom(42)
	l.Sort(CompareFunc[int]())
	assert.Equal(t, []int{42}, l.ToSlice())
}

func TestLinkedList_AnyEvery(t *testing.T) {
	l := NewLinkedListFrom(2, 4, 6, 8)

	assert.True(t, l.Any(func(v int) bool { return v > 5 }))
	assert.False(t, l.Any(func(v int) bool { return v > 100 }))

	assert.True(t, l.Every(func(v int) bool { return v%2 == 0 }))
	assert.False(t, l.Every(func(v int) bool { return v > 5 }))
}

func TestLinkedList_Seq(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3)

	var collected []int
	for v := range l.Seq() {
		collected = append(collected, v)
	}
	assert.Equal(t, []int{1, 2, 3}, collected)
}

func TestLinkedList_ForEach(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3, 4, 5)

	var sum int
	l.ForEach(func(v int) bool {
		sum += v
		return v < 3
	})
	assert.Equal(t, 6, sum) // 1 + 2 + 3
}

func TestLinkedList_String(t *testing.T) {
	l := NewLinkedListFrom(1, 2, 3)
	s := l.String()
	assert.Contains(t, s, "linkedList")
	assert.Contains(t, s, "1")
	assert.Contains(t, s, "2")
	assert.Contains(t, s, "3")
}

func TestLinkedList_AddAll(t *testing.T) {
	l := NewLinkedList[int]()
	l.AddAll(1, 2, 3)
	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
}

func TestLinkedList_AddSeq(t *testing.T) {
	l := NewLinkedList[int]()
	other := NewLinkedListFrom(1, 2, 3)
	l.AddSeq(other.Seq())
	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
}

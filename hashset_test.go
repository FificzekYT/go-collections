package collections

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper to make a Seq from slice
func seqOf[T any](vals []T) func(func(T) bool) {
	return func(yield func(T) bool) {
		for _, v := range vals {
			if !yield(v) {
				return
			}
		}
	}
}

func TestHashSet_BasicCRUD(t *testing.T) {
	s := NewHashSetFrom(1, 2, 2, 3)
	assert.Equal(t, 3, s.Size())
	assert.True(t, s.Contains(1) && s.Contains(2) && s.Contains(3))
	assert.False(t, s.Contains(9))
	added := s.Add(2)
	assert.False(t, added)
	assert.True(t, s.Add(4))
	assert.True(t, s.Remove(2))
	assert.False(t, s.Remove(2))
	// Pop
	_, ok := s.Pop()
	require.True(t, ok)
}

func TestHashSet_AddRemoveAllSeq(t *testing.T) {
	s := NewHashSet[int]()
	added := s.AddAll(1, 2, 3, 3)
	assert.Equal(t, 3, added)
	added = s.AddSeq(seqOf([]int{3, 4, 5}))
	assert.Equal(t, 2, added)
	removed := s.RemoveAll(5, 6)
	assert.Equal(t, 1, removed)
	removed = s.RemoveSeq(seqOf([]int{1, 4, 7}))
	assert.Equal(t, 2, removed)
}

func TestHashSet_RemoveRetainFunc(t *testing.T) {
	s := NewHashSetFrom(1, 2, 3, 4, 5)
	removed := s.RemoveFunc(func(v int) bool { return v%2 == 0 })
	require.Equal(t, 2, removed)
	removed = s.RetainFunc(func(v int) bool { return v > 3 }) // keep >3
	assert.Equal(t, 2, removed)                               // removed 1,3
	assert.True(t, s.Contains(5))
	assert.Equal(t, 1, s.Size(), "Size should be 1 after RetainFunc, state=%v", s.ToSlice())
}

func TestHashSet_SeqForEach(t *testing.T) {
	s := NewHashSetFrom(1, 2, 3, 4)
	collected := make([]int, 0, s.Size())
	for v := range s.Seq() {
		collected = append(collected, v)
	}
	slices.Sort(collected)
	assert.True(t, slices.Equal(collected, []int{1, 2, 3, 4}), "Seq collected=%v", collected)
	n := 0
	s.ForEach(func(v int) bool {
		n++
		return false // stop immediately
	})
	assert.Equal(t, 1, n)
}

func TestHashSet_Algebra(t *testing.T) {
	a := NewHashSetFrom(1, 2, 3, 4)
	b := NewHashSetFrom(3, 4, 5)
	u := a.Union(b).ToSlice()
	i := a.Intersection(b).ToSlice()
	d := a.Difference(b).ToSlice()
	sd := a.SymmetricDifference(b).ToSlice()
	slices.Sort(u)
	slices.Sort(i)
	slices.Sort(d)
	slices.Sort(sd)
	assert.True(t, slices.Equal(u, []int{1, 2, 3, 4, 5}), "Union=%v", u)
	assert.True(t, slices.Equal(i, []int{3, 4}), "Intersection=%v", i)
	assert.True(t, slices.Equal(d, []int{1, 2}), "Difference=%v", d)
	assert.True(t, slices.Equal(sd, []int{1, 2, 5}), "SymmetricDifference=%v", sd)
}

func TestHashSet_RelationsAndSearch(t *testing.T) {
	a := NewHashSetFrom(1, 2, 3)
	b := NewHashSetFrom(1, 2, 3, 4)
	c := NewHashSetFrom(4, 5)
	assert.True(t, a.IsSubsetOf(b))
	assert.True(t, b.IsSupersetOf(a))
	assert.True(t, a.IsProperSubsetOf(b))
	assert.True(t, b.IsProperSupersetOf(a))
	assert.True(t, a.IsDisjoint(c))
	assert.False(t, c.IsDisjoint(b))
	assert.True(t, a.Equals(NewHashSetFrom(2, 1, 3)))
	v, ok := a.Find(func(x int) bool { return x%2 == 0 })
	require.True(t, ok)
	assert.Equal(t, 0, v%2)
	assert.True(t, a.Any(func(x int) bool { return x == 2 }))
	assert.True(t, a.Every(func(x int) bool { return x >= 1 && x <= 3 }))
}

func TestHashSet_EmptySingle(t *testing.T) {
	s := NewHashSet[int]()
	assert.True(t, s.IsEmpty())
	assert.Equal(t, 0, s.Size())
	_, ok := s.Pop()
	assert.False(t, ok)
	assert.True(t, s.Add(1))
	assert.False(t, s.IsEmpty())
	assert.Equal(t, 1, s.Size())
	assert.True(t, s.Contains(1))
}

func TestHashSet_Large(t *testing.T) {
	s := NewHashSet[int]()
	const n = 20000
	for i := range n {
		s.Add(i)
	}
	assert.Equal(t, n, s.Size())
	for _, k := range []int{0, n / 2, n - 1} {
		require.Truef(t, s.Contains(k), "Missing %d", k)
	}
}

func TestHashSet_Clear(t *testing.T) {
	s := NewHashSetFrom(1, 2, 3)
	assert.False(t, s.IsEmpty())
	s.Clear()
	assert.True(t, s.IsEmpty())
	assert.Equal(t, 0, s.Size())
}

func TestHashSet_ToSlice(t *testing.T) {
	s := NewHashSetFrom(5, 3, 1, 4, 2)
	slice := s.ToSlice()
	slices.Sort(slice)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, slice)
}

func TestHashSet_Clone(t *testing.T) {
	s := NewHashSetFrom(1, 2, 3)
	clone := s.Clone()
	clone.Remove(1)
	assert.True(t, s.Contains(1))
	assert.Equal(t, 3, s.Size())
	assert.Equal(t, 2, clone.Size())
}

func TestHashSet_Filter(t *testing.T) {
	s := NewHashSetFrom(1, 2, 3, 4, 5, 6)
	filtered := s.Filter(func(v int) bool { return v%2 == 0 })
	assert.Equal(t, 3, filtered.Size())
	assert.True(t, filtered.Contains(2))
	assert.True(t, filtered.Contains(4))
	assert.True(t, filtered.Contains(6))
}

func TestHashSet_ContainsAllAny(t *testing.T) {
	s := NewHashSetFrom(1, 2, 3, 4, 5)
	assert.True(t, s.ContainsAll(1, 3, 5))
	assert.False(t, s.ContainsAll(1, 6))
	assert.True(t, s.ContainsAny(1, 10, 20))
	assert.False(t, s.ContainsAny(10, 20, 30))
}

func TestHashSet_NewWithCapacity(t *testing.T) {
	s := NewHashSetWithCapacity[int](10)
	require.True(t, s.IsEmpty())
	s.Add(1)
	require.Equal(t, 1, s.Size())
}

func TestHashSet_String(t *testing.T) {
	s := NewHashSet[int]()
	s.Add(1)
	str := s.String()
	require.Contains(t, str, "hashSet")
	require.Contains(t, str, "1")
}

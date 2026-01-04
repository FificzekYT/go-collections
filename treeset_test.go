package collections

import (
	"cmp"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Panic tests colocated with tree set unit tests.
func TestTreeSet_PanicOnNilComparator(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "expected panic on nil comparator")
		}
	}()
	_ = NewTreeSet[int](nil)
}

func TestTreeSet_BasicAndOrder(t *testing.T) {
	s := NewTreeSetOrdered[int]()
	s.AddAll(3, 1, 2, 2)
	got := s.ToSlice()
	slices.Sort(got)
	assert.True(t, slices.Equal(got, []int{1, 2, 3}), "ToSlice=%v", got)
	asc := make([]int, 0, 3)
	for v := range s.Seq() {
		asc = append(asc, v)
	}
	assert.True(t, slices.IsSorted(asc), "Seq not ascending: %v", asc)
	// Descending via Reversed
	dec := make([]int, 0, 3)
	for v := range s.Reversed() {
		dec = append(dec, v)
	}
	assert.Equal(t, 3, len(dec))
	assert.Equal(t, 3, dec[0])
	assert.Equal(t, 1, dec[2])
}

func TestTreeSet_NavigationAndExtremes(t *testing.T) {
	s := NewTreeSet(func(a, b int) int { return cmp.Compare(a, b) })
	s.AddAll(10, 20, 30, 40)
	v, ok := s.First()
	require.True(t, ok)
	assert.Equal(t, 10, v)
	v, ok = s.Last()
	require.True(t, ok)
	assert.Equal(t, 40, v)
	v, ok = s.Min()
	require.True(t, ok)
	assert.Equal(t, 10, v)
	v, ok = s.Max()
	require.True(t, ok)
	assert.Equal(t, 40, v)
	v, ok = s.Floor(25)
	require.True(t, ok)
	require.Equal(t, 20, v)
	v, ok = s.Ceiling(25)
	require.True(t, ok)
	require.Equal(t, 30, v)
	_, ok = s.Lower(10)
	require.False(t, ok, "Lower(10) should be none")
	_, ok = s.Higher(40)
	require.False(t, ok, "Higher(40) should be none")
	// PopFirst/PopLast
	v, ok = s.PopFirst()
	require.True(t, ok)
	require.Equal(t, 10, v)
	v, ok = s.PopLast()
	require.True(t, ok)
	require.Equal(t, 40, v)
}

func TestTreeSet_RangeAndRank(t *testing.T) {
	s := NewTreeSetOrdered[int]()
	s.AddAll(1, 2, 3, 4, 5)
	r := make([]int, 0)
	s.Range(2, 4, func(e int) bool {
		r = append(r, e)
		return true
	})
	require.True(t, slices.Equal(r, []int{2, 3, 4}), "Range=%v", r)
	r2 := make([]int, 0)
	for v := range s.RangeSeq(3, 5) {
		r2 = append(r2, v)
	}
	require.True(t, slices.Equal(r2, []int{3, 4, 5}), "RangeSeq=%v", r2)
	require.Equal(t, 2, s.Rank(3))
	v, ok := s.GetByRank(0)
	require.True(t, ok)
	require.Equal(t, 1, v)
}

func TestTreeSet_Views(t *testing.T) {
	s := NewTreeSetOrdered[int]()
	s.AddAll(1, 2, 3, 4, 5)
	head := s.HeadSet(3, false).ToSlice()
	tail := s.TailSet(3, true).ToSlice()
	sub := s.SubSet(2, 4).ToSlice()
	slices.Sort(head)
	slices.Sort(tail)
	slices.Sort(sub)
	require.True(t, slices.Equal(head, []int{1, 2}), "HeadSet=%v", head)
	require.True(t, slices.Equal(tail, []int{3, 4, 5}), "TailSet=%v", tail)
	require.True(t, slices.Equal(sub, []int{2, 3, 4}), "SubSet=%v", sub)
}

// Basic Set ops (emptiness, add/contains/remove) and algebra on TreeSet.
func TestTreeSet_SetOpsBasic(t *testing.T) {
	s := NewTreeSetOrdered[int]()
	// emptiness
	assert.True(t, s.IsEmpty())
	assert.Equal(t, 0, s.Size())
	// add & contains
	assert.True(t, s.Add(1))
	assert.False(t, s.Add(1)) // duplicate
	assert.True(t, s.Contains(1))
	// AddAll count (at least 2 because 3 is duplicate later)
	added := s.AddAll(2, 3, 3)
	assert.GreaterOrEqual(t, added, 2)
	// Remove existing / non-existing
	assert.True(t, s.Remove(2))
	assert.False(t, s.Remove(99))
	// Algebra with another TreeSet
	other := NewTreeSetOrdered[int]()
	other.AddAll(1, 3, 4)
	union := s.Union(other)
	assert.True(t, union.ContainsAll(1, 3, 4))
	inter := s.Intersection(other)
	assert.True(t, inter.ContainsAll(1, 3))
	diff := other.Difference(s)
	assert.True(t, diff.Contains(4))
	sym := s.SymmetricDifference(other)
	assert.True(t, sym.Contains(4) && !sym.Contains(1))
}

func TestTreeSet_AscendDescend(t *testing.T) {
	s := NewTreeSetOrdered[int]()
	s.AddAll(5, 1, 3, 2, 4)
	// Ascend - should iterate in ascending order
	asc := make([]int, 0)
	s.Ascend(func(e int) bool {
		asc = append(asc, e)
		return true
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, asc)

	// Descend - should iterate in descending order
	desc := make([]int, 0)
	s.Descend(func(e int) bool {
		desc = append(desc, e)
		return true
	})
	assert.Equal(t, []int{5, 4, 3, 2, 1}, desc)

	// AscendFrom - iterate elements >= 3 ascending
	ascFrom := make([]int, 0)
	s.AscendFrom(3, func(e int) bool {
		ascFrom = append(ascFrom, e)
		return true
	})
	assert.Equal(t, []int{3, 4, 5}, ascFrom)

	// DescendFrom - iterate elements <= 3 descending
	descFrom := make([]int, 0)
	s.DescendFrom(3, func(e int) bool {
		descFrom = append(descFrom, e)
		return true
	})
	assert.Equal(t, []int{3, 2, 1}, descFrom)

	// Test early termination
	count := 0
	s.Ascend(func(e int) bool {
		count++
		return count < 3
	})
	assert.Equal(t, 3, count)
}

func TestTreeSet_CloneSorted(t *testing.T) {
	s := NewTreeSetOrdered[int]()
	s.AddAll(1, 2, 3)
	clone := s.CloneSorted()
	clone.Remove(1)
	assert.True(t, s.Contains(1))
	assert.Equal(t, 3, s.Size())
	assert.Equal(t, 2, clone.Size())
}

func TestTreeSet_CoverageSupplement(t *testing.T) {
	s := NewTreeSetFrom(cmp.Compare[int], 1, 2, 3)
	assert.Equal(t, 3, s.Size())

	// String
	str := s.String()
	assert.Contains(t, str, "treeSet")
	assert.Contains(t, str, "1")

	// ForEach
	cnt := 0
	s.ForEach(func(e int) bool {
		cnt++
		return true
	})
	assert.Equal(t, 3, cnt)

	// Pop
	v, ok := s.Pop()
	require.True(t, ok)
	assert.Contains(t, []int{1, 2, 3}, v)
	assert.Equal(t, 2, s.Size())

	// AddSeq
	added := s.AddSeq(func(yield func(int) bool) {
		if !yield(4) {
			return
		}
		yield(5)
	})
	assert.Equal(t, 2, added)
	assert.True(t, s.ContainsAll(4, 5))

	// RemoveSeq
	removed := s.RemoveSeq(func(yield func(int) bool) {
		if !yield(4) {
			return
		}
		yield(99) // not present
	})
	assert.Equal(t, 1, removed)
	assert.Equal(t, 3, s.Size())

	// RemoveFunc
	removed = s.RemoveFunc(func(e int) bool { return e == 5 })
	assert.Equal(t, 1, removed)
	assert.False(t, s.Contains(5))

	// RetainFunc
	s.AddAll(10, 11)
	// Only retain > 5 (should remove existing small elements)
	removed = s.RetainFunc(func(e int) bool { return e > 5 })
	assert.Greater(t, removed, 0)
	assert.True(t, s.ContainsAll(10, 11))
	assert.False(t, s.Contains(2)) // was present before

	// Set Relations
	setA := NewTreeSetFrom(cmp.Compare[int], 1, 2)
	setB := NewTreeSetFrom(cmp.Compare[int], 1, 2, 3)
	setC := NewTreeSetFrom(cmp.Compare[int], 4, 5)

	assert.True(t, setA.IsSubsetOf(setB))
	assert.True(t, setB.IsSupersetOf(setA))
	assert.True(t, setA.IsProperSubsetOf(setB))
	assert.True(t, setB.IsProperSupersetOf(setA))
	assert.True(t, setA.IsDisjoint(setC))
	assert.False(t, setA.IsSubsetOf(setC))

	// Clone
	c := setA.Clone()
	assert.True(t, c.Equals(setA))

	// Filter
	f := setB.Filter(func(e int) bool { return e%2 == 0 })
	assert.Equal(t, 1, f.Size())
	assert.True(t, f.Contains(2))

	// Any, Every, Find
	assert.True(t, setB.Any(func(e int) bool { return e == 3 }))
	assert.True(t, setB.Every(func(e int) bool { return e > 0 }))
	val, ok := setB.Find(func(e int) bool { return e == 2 })
	assert.True(t, ok)
	assert.Equal(t, 2, val)
}

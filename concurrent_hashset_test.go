package collections

import (
	"runtime"
	"slices"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConcurrentHashSet_Basic(t *testing.T) {
	s := NewConcurrentHashSet[int]()
	assert.True(t, s.IsEmpty(), "New set should be empty")
	assert.True(t, s.Add(1))
	assert.False(t, s.Add(1), "Add duplicate should be false")
	assert.True(t, s.Contains(1))
	r, ok := s.RemoveAndGet(1)
	require.True(t, ok)
	assert.Equal(t, 1, r)
	assert.False(t, s.Contains(1))
	// RemoveAndGet non-existing
	_, ok = s.RemoveAndGet(99)
	assert.False(t, ok)
}

func TestConcurrentHashSet_ConcurrentAddIfAbsent(t *testing.T) {
	s := NewConcurrentHashSet[int]()
	n := 1000
	workers := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := range workers {
		go func(id int) {
			defer wg.Done()
			for i := range n {
				s.AddIfAbsent(i)
			}
		}(w)
	}
	wg.Wait()
	// validate presence of a few keys deterministically
	for _, k := range []int{0, n / 2, n - 1} {
		assert.Truef(t, s.Contains(k), "Missing key %d", k)
	}
}

func TestConcurrentHashSet_Algebra(t *testing.T) {
	a := NewConcurrentHashSetFrom(1, 2, 3, 4)
	b := NewConcurrentHashSetFrom(3, 4, 5)

	// Union
	u := a.Union(b).ToSlice()

	// Intersection
	i := a.Intersection(b).ToSlice()

	// Difference
	d := a.Difference(b).ToSlice()

	// SymmetricDifference
	sd := a.SymmetricDifference(b).ToSlice()

	expectContains := func(slice []int, x int) bool {
		return slices.Contains(slice, x)
	}

	assert.Truef(t, expectContains(u, 1) && expectContains(u, 5), "Union unexpected: %v", u)
	assert.Len(t, i, 2, "Intersection unexpected: %v", i)
	assert.True(t, expectContains(i, 3) && expectContains(i, 4))

	assert.Len(t, d, 2, "Difference unexpected: %v", d)
	assert.True(t, expectContains(d, 1) && expectContains(d, 2))

	assert.Len(t, sd, 3, "SymmetricDifference unexpected: %v", sd)
	assert.True(t, expectContains(sd, 1) && expectContains(sd, 2) && expectContains(sd, 5))
}

func TestConcurrentHashSet_ClearAndSize(t *testing.T) {
	s := NewConcurrentHashSetFrom(1, 2, 3)
	assert.False(t, s.IsEmpty())
	assert.Equal(t, 3, s.Size())

	// ToSlice
	slice := s.ToSlice()
	assert.Len(t, slice, 3)

	// String
	str := s.String()
	assert.Contains(t, str, "concurrentHashSet")
	assert.Contains(t, str, "1")

	s.Clear()
	assert.True(t, s.IsEmpty())
	assert.Equal(t, 0, s.Size())
}

func TestConcurrentHashSet_PopAndRemove(t *testing.T) {
	s := NewConcurrentHashSetFrom(1, 2, 3)
	// Pop returns arbitrary element that is removed
	v, ok := s.Pop()
	require.True(t, ok)
	assert.Contains(t, []int{1, 2, 3}, v)
	assert.Equal(t, 2, s.Size())

	// Pop on empty
	s.Clear()
	_, ok = s.Pop()
	assert.False(t, ok)

	s.Add(4)
	// Remove
	assert.True(t, s.Remove(4))
	assert.False(t, s.Remove(4))
}

func TestConcurrentHashSet_AddAndRemoveAll(t *testing.T) {
	s := NewConcurrentHashSet[int]()

	// AddAll
	added := s.AddAll(1, 2, 3)
	assert.Equal(t, 3, added)
	added = s.AddAll(2, 4) // 2 exists, 4 is new
	assert.Equal(t, 1, added)
	assert.True(t, s.Contains(4))

	// AddSeq
	added = s.AddSeq(func(yield func(int) bool) {
		if !yield(5) {
			return
		}
		if !yield(6) {
			return
		}
		yield(1) // exists
	})
	assert.Equal(t, 2, added)

	// RemoveAll
	removed := s.RemoveAll(1, 2, 99)
	assert.Equal(t, 2, removed)
	assert.False(t, s.Contains(1))

	// RemoveSeq
	removed = s.RemoveSeq(func(yield func(int) bool) {
		if !yield(3) {
			return
		}
		if !yield(4) {
			return
		}
		yield(88)
	})
	assert.Equal(t, 2, removed)
	assert.False(t, s.Contains(3))
}

func TestConcurrentHashSet_ContainsOperations(t *testing.T) {
	s := NewConcurrentHashSetFrom(1, 2, 3)

	// ContainsAll
	assert.True(t, s.ContainsAll(1, 2))
	assert.False(t, s.ContainsAll(1, 4))

	// ContainsAny
	assert.True(t, s.ContainsAny(1, 4))
	assert.False(t, s.ContainsAny(4, 5))
}

func TestConcurrentHashSet_SubsetSuperset(t *testing.T) {
	s1 := NewConcurrentHashSetFrom(1, 2)
	s2 := NewConcurrentHashSetFrom(1, 2, 3)
	s3 := NewConcurrentHashSetFrom(2, 3)
	s4 := NewConcurrentHashSetFrom(4, 5)

	// IsSubsetOf
	assert.True(t, s1.IsSubsetOf(s2))
	assert.False(t, s2.IsSubsetOf(s1))

	// IsSupersetOf
	assert.True(t, s2.IsSupersetOf(s1))

	// IsProperSubsetOf
	assert.True(t, s1.IsProperSubsetOf(s2))
	assert.False(t, s1.IsProperSubsetOf(s1))

	// IsProperSupersetOf
	assert.True(t, s2.IsProperSupersetOf(s1))

	// IsDisjoint
	assert.False(t, s1.IsDisjoint(s2))
	assert.False(t, s2.IsDisjoint(s3))
	assert.True(t, s1.IsDisjoint(s4))
}

func TestConcurrentHashSet_CloneFilterEquals(t *testing.T) {
	s := NewConcurrentHashSetFrom(1, 2, 3)

	// Clone
	c := s.Clone()
	assert.True(t, s.Equals(c))

	c.Add(4)
	assert.False(t, s.Equals(c))

	// Filter
	even := s.Filter(func(e int) bool { return e%2 == 0 })
	assert.Equal(t, 1, even.Size())
	assert.True(t, even.Contains(2))

	// Equals comparison
	s2 := NewHashSetFrom(1, 2, 3)
	assert.True(t, s.Equals(s2))
}

func TestConcurrentHashSet_IteratorsAndPredicates(t *testing.T) {
	s := NewConcurrentHashSetFrom(1, 2, 3, 4, 5)

	// ForEach
	cnt := 0
	s.ForEach(func(e int) bool {
		cnt++
		return true
	})
	assert.Equal(t, 5, cnt)

	// Seq
	cnt = 0
	for range s.Seq() {
		cnt++
	}
	assert.Equal(t, 5, cnt)

	// Find
	v, ok := s.Find(func(e int) bool { return e > 3 })
	assert.True(t, ok)
	assert.Contains(t, []int{4, 5}, v)

	_, ok = s.Find(func(e int) bool { return e > 10 })
	assert.False(t, ok)

	// Any
	assert.True(t, s.Any(func(e int) bool { return e == 3 }))
	assert.False(t, s.Any(func(e int) bool { return e == 10 }))

	// Every
	assert.True(t, s.Every(func(e int) bool { return e > 0 }))
	assert.False(t, s.Every(func(e int) bool { return e < 5 }))

	// RemoveFunc
	removed := s.RemoveFunc(func(e int) bool { return e%2 == 0 }) // remove 2, 4
	assert.Equal(t, 2, removed)
	assert.False(t, s.Contains(2))

	// RetainFunc
	s.Add(2)
	s.Add(4)
	removed = s.RetainFunc(func(e int) bool { return e%2 != 0 }) // keep odds: 1, 3, 5; remove 2, 4
	assert.Equal(t, 2, removed)
	assert.False(t, s.Contains(2))
}

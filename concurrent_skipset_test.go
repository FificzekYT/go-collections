package collections

import (
	"runtime"
	"slices"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConcurrentSkipSet_BasicOrdered(t *testing.T) {
	s := NewConcurrentSkipSet[int]()
	s.AddAll(3, 1, 2, 2)
	got := make([]int, 0, s.Size())
	for v := range s.Seq() {
		got = append(got, v)
	}
	assert.True(t, slices.Equal(got, []int{1, 2, 3}), "Seq=%v", got)
	rev := make([]int, 0, 3)
	for v := range s.Reversed() {
		rev = append(rev, v)
	}
	assert.True(t, slices.Equal(rev, []int{3, 2, 1}), "Reversed=%v", rev)
	v, ok := s.First()
	require.True(t, ok)
	assert.Equal(t, 1, v)
	v, ok = s.Last()
	require.True(t, ok)
	assert.Equal(t, 3, v)
}

func TestConcurrentSkipSet_ConcurrentAddIfAbsent(t *testing.T) {
	s := NewConcurrentSkipSet[int]()
	n := 1000
	workers := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for i := range n {
				s.AddIfAbsent(i)
			}
		}()
	}
	wg.Wait()
	for _, k := range []int{0, n / 2, n - 1} {
		assert.Truef(t, s.Contains(k), "Missing key %d", k)
	}
}

func TestConcurrentSkipSet_NavigationAndExtremes(t *testing.T) {
	s := NewConcurrentSkipSet[int]()
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
}

func TestConcurrentSkipSet_PopFirstLast(t *testing.T) {
	s := NewConcurrentSkipSet[int]()
	s.AddAll(10, 20, 30, 40)
	v, ok := s.PopFirst()
	require.True(t, ok)
	require.Equal(t, 10, v)
	v, ok = s.PopLast()
	require.True(t, ok)
	require.Equal(t, 40, v)
	require.Equal(t, 2, s.Size())
}

func TestConcurrentSkipSet_RangeAndRank(t *testing.T) {
	s := NewConcurrentSkipSet[int]()
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

func TestConcurrentSkipSet_Views(t *testing.T) {
	s := NewConcurrentSkipSet[int]()
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

func TestConcurrentSkipSet_SetOpsBasic(t *testing.T) {
	s := NewConcurrentSkipSet[int]()
	assert.True(t, s.IsEmpty())
	assert.Equal(t, 0, s.Size())
	assert.True(t, s.Add(1))
	assert.False(t, s.Add(1))
	assert.True(t, s.Contains(1))
	added := s.AddAll(2, 3, 3)
	assert.GreaterOrEqual(t, added, 2)
	assert.True(t, s.Remove(2))
	assert.False(t, s.Remove(99))
}

func TestConcurrentSkipSet_Algebra(t *testing.T) {
	a := NewConcurrentSkipSetFrom(1, 2, 3, 4)
	b := NewConcurrentSkipSetFrom(3, 4, 5)
	u := a.Union(b).ToSlice()
	i := a.Intersection(b).ToSlice()
	d := a.Difference(b).ToSlice()
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

func TestConcurrentSkipSet_AscendDescend(t *testing.T) {
	s := NewConcurrentSkipSet[int]()
	s.AddAll(5, 1, 3, 2, 4)
	asc := make([]int, 0)
	s.Ascend(func(e int) bool {
		asc = append(asc, e)
		return true
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, asc)
	desc := make([]int, 0)
	s.Descend(func(e int) bool {
		desc = append(desc, e)
		return true
	})
	assert.Equal(t, []int{5, 4, 3, 2, 1}, desc)
	af := make([]int, 0)
	s.AscendFrom(3, func(e int) bool {
		af = append(af, e)
		return true
	})
	assert.Equal(t, []int{3, 4, 5}, af)
	df := make([]int, 0)
	s.DescendFrom(3, func(e int) bool {
		df = append(df, e)
		return true
	})
	assert.Equal(t, []int{3, 2, 1}, df)
}

func TestConcurrentSkipSet_CloneSorted(t *testing.T) {
	s := NewConcurrentSkipSet[int]()
	s.AddAll(1, 2, 3)
	clone := s.CloneSorted()
	clone.Remove(1)
	assert.True(t, s.Contains(1))
	assert.Equal(t, 3, s.Size())
	assert.Equal(t, 2, clone.Size())
}

func TestConcurrentSkipSet_FilterFindAnyEvery(t *testing.T) {
	s := NewConcurrentSkipSet[int]()
	s.AddAll(1, 2, 3, 4, 5, 6)
	filtered := s.Filter(func(e int) bool { return e%2 == 0 })
	assert.Equal(t, 3, filtered.Size())
	v, ok := filtered.Find(func(e int) bool { return e > 4 })
	require.True(t, ok)
	assert.Equal(t, 6, v)
	assert.True(t, s.Any(func(e int) bool { return e == 3 }))
	assert.True(t, s.Every(func(e int) bool { return e >= 1 }))
	assert.False(t, s.Every(func(e int) bool { return e > 3 }))
}

func TestConcurrentSkipSet_RemoveFuncRetainFunc(t *testing.T) {
	s := NewConcurrentSkipSetFrom(1, 2, 3, 4, 5, 6)
	count := s.RemoveFunc(func(e int) bool { return e%2 == 0 })
	assert.Equal(t, 3, count)
	assert.Equal(t, 3, s.Size())
	assert.True(t, s.Contains(1))
	assert.False(t, s.Contains(2))
	s.AddAll(7, 8, 9)
	count = s.RetainFunc(func(e int) bool { return e > 5 })
	// Keep 7, 8, 9 (all > 5), remove 1, 3, 5 = 3 removed
	assert.Equal(t, 3, count)
	assert.Equal(t, 3, s.Size())
	assert.False(t, s.Contains(1))
	assert.False(t, s.Contains(3))
	assert.False(t, s.Contains(5))
	assert.True(t, s.Contains(7))
	assert.True(t, s.Contains(8))
	assert.True(t, s.Contains(9))
}

func TestConcurrentSkipSet_ClearAndSize(t *testing.T) {
	s := NewConcurrentSkipSetFrom(1, 2, 3)
	assert.False(t, s.IsEmpty())
	assert.Equal(t, 3, s.Size())
	s.Clear()
	assert.True(t, s.IsEmpty())
	assert.Equal(t, 0, s.Size())
}

func TestConcurrentSkipSet_ToSliceAndPop(t *testing.T) {
	s := NewConcurrentSkipSetFrom(5, 3, 1, 4, 2)
	slice := s.ToSlice()
	slices.Sort(slice)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, slice)
	v, ok := s.Pop()
	require.True(t, ok)
	assert.Contains(t, []int{1, 2, 3, 4, 5}, v)
	assert.Equal(t, 4, s.Size())
}

func TestConcurrentSkipSet_AddSeqRemoveSeq(t *testing.T) {
	s := NewConcurrentSkipSet[int]()
	s.AddSeq(func(yield func(int) bool) {
		if !yield(1) {
			return
		}
		if !yield(2) {
			return
		}
		yield(3)
	})
	assert.Equal(t, 3, s.Size())
	removed := s.RemoveSeq(func(yield func(int) bool) {
		if !yield(1) {
			return
		}
		yield(2)
	})
	assert.Equal(t, 2, removed)
	assert.Equal(t, 1, s.Size())
}

func TestConcurrentSkipSet_ContainsAllAny(t *testing.T) {
	s := NewConcurrentSkipSetFrom(1, 2, 3, 4, 5)
	assert.True(t, s.ContainsAll(1, 3, 5))
	assert.False(t, s.ContainsAll(1, 6))
	assert.True(t, s.ContainsAny(1, 10, 20))
	assert.False(t, s.ContainsAny(10, 20, 30))
}

func TestConcurrentSkipSet_IsSubsetSupersetDisjoint(t *testing.T) {
	a := NewConcurrentSkipSetFrom(1, 2, 3)
	b := NewConcurrentSkipSetFrom(1, 2, 3, 4, 5)
	c := NewConcurrentSkipSetFrom(4, 5, 6)
	assert.True(t, a.IsSubsetOf(b))
	assert.False(t, b.IsSubsetOf(a))
	assert.True(t, b.IsSupersetOf(a))
	assert.True(t, a.IsProperSubsetOf(b))
	assert.False(t, a.IsProperSubsetOf(a))
	assert.True(t, a.IsDisjoint(c))
	assert.False(t, a.IsDisjoint(b))
}

func TestConcurrentSkipSet_Equals(t *testing.T) {
	a := NewConcurrentSkipSetFrom(1, 2, 3)
	b := NewConcurrentSkipSetFrom(1, 2, 3)
	c := NewConcurrentSkipSetFrom(1, 2, 4)
	assert.True(t, a.Equals(b))
	assert.False(t, a.Equals(c))
}

func TestConcurrentSkipSet_Races(t *testing.T) {
	s := NewConcurrentSkipSet[int]()
	workers := runtime.GOMAXPROCS(0) * 2
	iters := 500
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := range workers {
		go func(id int) {
			defer wg.Done()
			for i := range iters {
				k := (id * 997) ^ i
				switch i % 3 {
				case 0:
					s.Add(k)
				case 1:
					s.Remove(k)
				default:
					_ = s.Contains(k)
				}
				if i%100 == 0 {
					count := 0
					s.Range(-1<<31, 1<<31-1, func(e int) bool {
						if count > 10 {
							return false
						}
						count++
						return true
					})
				}
			}
		}(w)
	}
	wg.Wait()
}

func TestConcurrentSkipSet_RemoveAll(t *testing.T) {
	s := NewConcurrentSkipSetFrom(1, 2, 3, 4, 5)
	removed := s.RemoveAll(1, 3, 5)
	assert.Equal(t, 3, removed)
	assert.Equal(t, 2, s.Size())
	assert.False(t, s.Contains(1))
	assert.False(t, s.Contains(3))
	assert.False(t, s.Contains(5))
}

func TestConcurrentSkipSet_AddIfAbsentAndRemoveAndGet(t *testing.T) {
	s := NewConcurrentSkipSet[int]()
	assert.True(t, s.AddIfAbsent(1))
	assert.False(t, s.AddIfAbsent(1))
	v, ok := s.RemoveAndGet(1)
	require.True(t, ok)
	assert.Equal(t, 1, v)
	_, ok = s.RemoveAndGet(1)
	assert.False(t, ok)
}

func TestConcurrentSkipSet_String(t *testing.T) {
	s := NewConcurrentSkipSet[int]()
	s.Add(1)
	str := s.String()
	require.Contains(t, str, "concurrentSkipSet")
	require.Contains(t, str, "1")
}

func TestConcurrentSkipSet_CoverageSupplement(t *testing.T) {
	s := NewConcurrentSkipSetFrom(1, 2, 3)

	// ForEach
	cnt := 0
	s.ForEach(func(e int) bool {
		cnt++
		return true
	})
	require.Equal(t, 3, cnt)

	// Clone
	c := s.Clone()
	require.True(t, s.Equals(c))

	// IsProperSupersetOf
	sub := NewConcurrentSkipSetFrom(1, 2)
	assert.True(t, s.IsProperSupersetOf(sub))

	// Find
	found, fok := s.Find(func(e int) bool { return e == 2 })
	require.True(t, fok)
	assert.Equal(t, 2, found)

	// Filter
	filter := s.Filter(func(e int) bool { return e > 1 })
	assert.Equal(t, 2, filter.Size())

	// Lower, Higher extra checks
	// s has 1, 2, 3
	_, ok := s.Lower(1)
	assert.False(t, ok)
	v, ok := s.Lower(3)
	assert.True(t, ok)
	assert.Equal(t, 2, v)

	v, ok = s.Higher(1)
	assert.True(t, ok)
	assert.Equal(t, 2, v)
	_, ok = s.Higher(3)
	assert.False(t, ok)

	// RangeSeq - range is inclusive [from, to], includes 1, 2, 3
	count := 0
	for range s.RangeSeq(1, 3) {
		count++
	}
	assert.Equal(t, 3, count, "RangeSeq(1, 3) should include elements 1, 2, 3")
}

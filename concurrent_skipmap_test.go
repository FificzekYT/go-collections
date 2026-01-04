package collections

import (
	"runtime"
	"slices"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConcurrentSkipMap_BasicOrdered(t *testing.T) {
	m := NewConcurrentSkipMap[int, string]()
	m.Put(2, "b")
	m.Put(1, "a")
	m.Put(3, "c")
	require.True(t, slices.Equal(m.Keys(), []int{1, 2, 3}), "Keys=%v", m.Keys())
	// Reversed
	rev := make([]int, 0)
	for k := range m.Reversed() {
		rev = append(rev, k)
	}
	require.True(t, slices.Equal(rev, []int{3, 2, 1}), "Reversed=%v", rev)
	// First/Last
	e, ok := m.FirstEntry()
	require.True(t, ok)
	require.Equal(t, 1, e.Key)
	require.Equal(t, "a", e.Value)
	e, ok = m.LastEntry()
	require.True(t, ok)
	require.Equal(t, 3, e.Key)
	require.Equal(t, "c", e.Value)
}

func TestConcurrentSkipMap_PutIfAbsentAndAtomics(t *testing.T) {
	m := NewConcurrentSkipMap[int, int]()
	n := 1000
	workers := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for i := range n {
				m.PutIfAbsent(i, i)
			}
		}()
	}
	wg.Wait()
	for _, k := range []int{0, n / 2, n - 1} {
		v, ok := m.Get(k)
		require.Truef(t, ok, "Missing key %d", k)
		require.Equal(t, k, v)
	}
	// Atomics
	v, computed := m.GetOrCompute(n+1, func() int { return 42 })
	require.True(t, computed)
	require.Equal(t, 42, v)
	v, loaded := m.LoadOrStore(n+1, 99)
	require.True(t, loaded)
	require.Equal(t, 42, v)
	// CAS
	require.False(t, m.CompareAndSwap(n+2, 0, 1, eqV[int]))
	m.Put(n+2, 10)
	require.True(t, m.CompareAndSwap(n+2, 10, 11, eqV[int]))
	require.True(t, m.CompareAndDelete(n+2, 11, eqV[int]))
	_, ok := m.Get(n + 2)
	require.False(t, ok, "Key should be deleted")
}

func TestConcurrentSkipMap_Navigation(t *testing.T) {
	m := NewConcurrentSkipMap[int, int]()
	for i := 1; i <= 5; i++ {
		m.Put(i*10, i)
	}
	k, ok := m.FirstKey()
	require.True(t, ok)
	require.Equal(t, 10, k)
	k, ok = m.LastKey()
	require.True(t, ok)
	require.Equal(t, 50, k)
	k, ok = m.FloorKey(33)
	require.True(t, ok)
	require.Equal(t, 30, k)
	k, ok = m.CeilingKey(33)
	require.True(t, ok)
	require.Equal(t, 40, k)
	_, ok = m.LowerKey(10)
	require.False(t, ok, "LowerKey(10) should be none (no element < 10)")
	_, ok = m.HigherKey(50)
	require.False(t, ok, "HigherKey(50) should be none (no element > 50)")
	// Test with element removed
	m.Remove(10)
	_, ok = m.LowerKey(20)
	require.False(t, ok, "LowerKey(20) should be none after removing 10")
}

func TestConcurrentSkipMap_RangeRankAndViews(t *testing.T) {
	m := NewConcurrentSkipMap[int, int]()
	m.Put(10, 1)
	m.Put(20, 2)
	m.Put(30, 3)
	m.Put(40, 4)
	// Range
	rs := make([]int, 0)
	m.Range(15, 35, func(k, v int) bool {
		rs = append(rs, k)
		return true
	})
	require.True(t, slices.Equal(rs, []int{20, 30}), "Range=%v", rs)
	// RangeSeq
	r2 := make([]int, 0)
	for k := range m.RangeSeq(15, 35) {
		r2 = append(r2, k)
	}
	require.True(t, slices.Equal(r2, []int{20, 30}), "RangeSeq=%v", r2)
	// Rank
	require.Equal(t, 2, m.RankOfKey(30))
	e, ok := m.GetByRank(1)
	require.True(t, ok)
	require.Equal(t, 20, e.Key)
	// SubMap, HeadMap, TailMap
	sub := m.SubMap(15, 35).Keys()
	head := m.HeadMap(25, false).Keys()
	tail := m.TailMap(25, true).Keys()
	require.True(t, slices.Equal(sub, []int{20, 30}), "SubMap=%v", sub)
	require.True(t, slices.Equal(head, []int{10, 20}), "HeadMap=%v", head)
	require.True(t, slices.Equal(tail, []int{30, 40}), "TailMap=%v", tail)
}

func TestConcurrentSkipMap_ClearAndClone(t *testing.T) {
	m := NewConcurrentSkipMapFrom(map[int]string{1: "a", 2: "b", 3: "c"})
	require.False(t, m.IsEmpty())
	require.Equal(t, 3, m.Size())
	// Clone
	clone := m.Clone()
	clone.Put(4, "d")
	require.Equal(t, 3, m.Size())
	require.Equal(t, 4, clone.Size())
	// CloneSorted
	sortedClone := m.CloneSorted()
	require.Equal(t, 3, sortedClone.Size())
	// Clear
	m.Clear()
	require.True(t, m.IsEmpty())
	require.Equal(t, 0, m.Size())
}

func TestConcurrentSkipMap_PutAllAndPutSeq(t *testing.T) {
	m := NewConcurrentSkipMap[int, int]()
	other := NewConcurrentSkipMap[int, int]()
	other.Put(1, 10)
	other.Put(2, 20)
	m.PutAll(other)
	require.Equal(t, 2, m.Size())
	// PutSeq
	m.PutSeq(func(yield func(int, int) bool) {
		if !yield(3, 30) {
			return
		}
		yield(4, 40)
	})
	require.Equal(t, 4, m.Size())
}

func TestConcurrentSkipMap_RemoveKeysAndFunc(t *testing.T) {
	m := NewConcurrentSkipMap[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	m.Put(3, "c")
	// RemoveKeys
	removed := m.RemoveKeys(1, 2)
	require.Equal(t, 2, removed)
	require.Equal(t, 1, m.Size())
	// RemoveFunc
	m.Put(4, "d")
	m.Put(5, "e")
	count := m.RemoveFunc(func(k int, v string) bool { return k > 3 })
	require.Equal(t, 2, count)
	require.Equal(t, 1, m.Size())
}

func TestConcurrentSkipMap_ComputeAndMerge(t *testing.T) {
	m := NewConcurrentSkipMap[int, string]()
	m.Put(1, "one")
	// Compute
	v, ok := m.Compute(1, func(k int, old string, exists bool) (string, bool) {
		if exists {
			return old + "_updated", true
		}
		return "new", true
	})
	require.True(t, ok)
	require.Equal(t, "one_updated", v)
	// ComputeIfAbsent
	v = m.ComputeIfAbsent(2, func(k int) string { return "two" })
	require.Equal(t, "two", v)
	v = m.ComputeIfAbsent(2, func(k int) string { return "two_again" })
	require.Equal(t, "two", v) // should not change
	// ComputeIfPresent
	v, ok = m.ComputeIfPresent(1, func(k int, old string) (string, bool) {
		return old + "_modified", true
	})
	require.True(t, ok)
	require.Equal(t, "one_updated_modified", v)
	// Merge
	v, ok = m.Merge(3, "three", func(old, new string) (string, bool) {
		return old + "_" + new, true
	})
	require.True(t, ok)
	require.Equal(t, "three", v)
	v, ok = m.Merge(1, "one", func(old, new string) (string, bool) {
		return old + "_merged", true
	})
	require.True(t, ok)
	require.Equal(t, "one_updated_modified_merged", v)
}

func TestConcurrentSkipMap_ReplaceOperations(t *testing.T) {
	m := NewConcurrentSkipMap[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	// Replace
	v, ok := m.Replace(1, "aa")
	require.True(t, ok)
	require.Equal(t, "a", v)
	v, _ = m.Get(1)
	require.Equal(t, "aa", v)
	// ReplaceIf
	require.True(t, m.ReplaceIf(1, "aa", "aaa", eqV[string]))
	require.False(t, m.ReplaceIf(1, "wrong", "xxx", eqV[string]))
	// ReplaceAll
	m.ReplaceAll(func(k int, v string) string { return v + "_replaced" })
	v, _ = m.Get(1)
	require.Equal(t, "aaa_replaced", v)
}

func TestConcurrentSkipMap_PopFirstLast(t *testing.T) {
	m := NewConcurrentSkipMap[int, int]()
	m.Put(10, 1)
	m.Put(20, 2)
	m.Put(30, 3)
	e, ok := m.PopFirst()
	require.True(t, ok)
	require.Equal(t, 10, e.Key)
	require.Equal(t, 1, e.Value)
	e, ok = m.PopLast()
	require.True(t, ok)
	require.Equal(t, 30, e.Key)
	require.Equal(t, 3, e.Value)
	require.Equal(t, 1, m.Size())
}

func TestConcurrentSkipMap_RangeFromAndTo(t *testing.T) {
	m := NewConcurrentSkipMap[int, int]()
	for i := 1; i <= 5; i++ {
		m.Put(i, i*10)
	}
	// RangeFrom
	rf := make([]int, 0)
	m.RangeFrom(3, func(k, v int) bool {
		rf = append(rf, k)
		return true
	})
	require.Equal(t, []int{3, 4, 5}, rf)
	// RangeTo
	rt := make([]int, 0)
	m.RangeTo(3, func(k, v int) bool {
		rt = append(rt, k)
		return true
	})
	require.Equal(t, []int{1, 2, 3}, rt)
}

func TestConcurrentSkipMap_AscendDescend(t *testing.T) {
	m := NewConcurrentSkipMap[int, int]()
	m.Put(5, 50)
	m.Put(1, 10)
	m.Put(3, 30)
	m.Put(2, 20)
	m.Put(4, 40)
	// Ascend
	asc := make([]int, 0)
	m.Ascend(func(k, v int) bool {
		asc = append(asc, k)
		return true
	})
	require.Equal(t, []int{1, 2, 3, 4, 5}, asc)
	// Descend
	desc := make([]int, 0)
	m.Descend(func(k, v int) bool {
		desc = append(desc, k)
		return true
	})
	require.Equal(t, []int{5, 4, 3, 2, 1}, desc)
	// AscendFrom
	af := make([]int, 0)
	m.AscendFrom(3, func(k, v int) bool {
		af = append(af, k)
		return true
	})
	require.Equal(t, []int{3, 4, 5}, af)
	// DescendFrom
	df := make([]int, 0)
	m.DescendFrom(3, func(k, v int) bool {
		df = append(df, k)
		return true
	})
	require.Equal(t, []int{3, 2, 1}, df)
}

func TestConcurrentSkipMap_EntryNavigation(t *testing.T) {
	m := NewConcurrentSkipMap[int, int]()
	for i := 1; i <= 5; i++ {
		m.Put(i*10, i)
	}
	e, ok := m.FloorEntry(35)
	require.True(t, ok)
	require.Equal(t, 30, e.Key)
	e, ok = m.CeilingEntry(35)
	require.True(t, ok)
	require.Equal(t, 40, e.Key)
	e, ok = m.LowerEntry(25)
	require.True(t, ok)
	require.Equal(t, 20, e.Key)
	e, ok = m.HigherEntry(25)
	require.True(t, ok)
	require.Equal(t, 30, e.Key)
}

func TestConcurrentSkipMap_FilterAndEquals(t *testing.T) {
	m := NewConcurrentSkipMap[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	m.Put(3, "c")
	// Filter
	filtered := m.Filter(func(k int, v string) bool { return k > 1 })
	require.Equal(t, 2, filtered.Size())
	// Equals
	m2 := NewConcurrentSkipMap[int, string]()
	m2.Put(1, "a")
	m2.Put(2, "b")
	m2.Put(3, "c")
	require.True(t, m.Equals(m2, eqV[string]))
	m3 := NewConcurrentSkipMap[int, string]()
	m3.Put(1, "a")
	m3.Put(2, "b")
	m3.Put(3, "different")
	require.False(t, m.Equals(m3, eqV[string]))
}

func TestConcurrentSkipMap_Races(t *testing.T) {
	m := NewConcurrentSkipMap[int, int]()
	workers := runtime.GOMAXPROCS(0) * 2
	iters := 500
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := range workers {
		go func(id int) {
			defer wg.Done()
			for i := range iters {
				k := (id * 131) + i
				switch i % 4 {
				case 0:
					m.Put(k, i)
				case 1:
					m.Remove(k)
				case 2:
					m.Get(k)
				default:
					m.ReplaceAll(func(key, val int) int { return val })
				}
				if i%100 == 0 {
					count := 0
					m.Range(-1<<31, 1<<31-1, func(k, v int) bool {
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

func TestConcurrentSkipMap_ContainsValueAndRemoveIf(t *testing.T) {
	m := NewConcurrentSkipMap[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	m.Put(3, "a")
	require.True(t, m.ContainsValue("a", eqV[string]))
	require.False(t, m.ContainsValue("z", eqV[string]))
	// RemoveIf
	require.True(t, m.RemoveIf(2, "b", eqV[string]))
	require.False(t, m.RemoveIf(2, "wrong", eqV[string]))
	require.Equal(t, 2, m.Size())
}

func TestConcurrentSkipMap_LoadAndDelete(t *testing.T) {
	m := NewConcurrentSkipMap[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	v, ok := m.LoadAndDelete(1)
	require.True(t, ok)
	require.Equal(t, "a", v)
	_, ok = m.Get(1)
	require.False(t, ok)
}

func TestConcurrentSkipMap_RemoveKeysSeq(t *testing.T) {
	m := NewConcurrentSkipMap[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(5, "e")
	removed := m.RemoveKeysSeq(func(yield func(int) bool) {
		if !yield(1) {
			return
		}
		if !yield(3) {
			return
		}
		if !yield(5) {
			return
		}
	})
	require.Equal(t, 3, removed)
	require.Equal(t, 2, m.Size())
}

func TestConcurrentSkipMap_CoverageSupplement(t *testing.T) {
	m := NewConcurrentSkipMap[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")

	// ForEach
	cnt := 0
	m.ForEach(func(k int, v string) bool {
		cnt++
		return true
	})
	require.Equal(t, 2, cnt)

	// SeqKeys, SeqValues
	kc := 0
	for range m.SeqKeys() {
		kc++
	}
	require.Equal(t, 2, kc)

	vc := 0
	for range m.SeqValues() {
		vc++
	}
	require.Equal(t, 2, vc)

	// String
	str := m.String()
	require.Contains(t, str, "concurrentSkipMap")
}

package collections

import (
	"cmp"
	"runtime"
	"slices"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConcurrentTreeMap_BasicOrdered(t *testing.T) {
	m := NewConcurrentTreeMapOrdered[int, string]()
	m.Put(2, "b")
	m.Put(1, "a")
	m.Put(3, "c")
	require.True(t, slices.Equal(m.Keys(), []int{1, 2, 3}), "Keys=%v", m.Keys())
	rev := make([]int, 0)
	for k := range m.Reversed() {
		rev = append(rev, k)
	}
	require.True(t, slices.Equal(rev, []int{3, 2, 1}), "Reversed=%v", rev)
	e, ok := m.FirstEntry()
	require.True(t, ok)
	require.Equal(t, 1, e.Key)
	require.Equal(t, "a", e.Value)
	e, ok = m.LastEntry()
	require.True(t, ok)
	require.Equal(t, 3, e.Key)
	require.Equal(t, "c", e.Value)
}

func TestConcurrentTreeMap_CustomComparator(t *testing.T) {
	cmpRev := func(a, b int) int { return -cmp.Compare(a, b) }
	m := NewConcurrentTreeMap[int, int](cmpRev)
	for i := 1; i <= 3; i++ {
		m.Put(i, i)
	}
	require.True(t, slices.Equal(m.Keys(), []int{3, 2, 1}), "Keys with reverse comparator=%v", m.Keys())
}

func TestConcurrentTreeMap_ConcurrentPutIfAbsentAndAtomics(t *testing.T) {
	m := NewConcurrentTreeMapOrdered[int, int]()
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
	// first CAS should fail (no key)
	require.False(t, m.CompareAndSwap(n+2, 0, 1, eqV[int]))
	m.Put(n+2, 10)
	require.True(t, m.CompareAndSwap(n+2, 10, 11, eqV[int]))
	require.True(t, m.CompareAndDelete(n+2, 11, eqV[int]))
	_, ok := m.Get(n + 2)
	require.False(t, ok, "Key should be deleted")

	// LoadAndDelete
	m.Put(n+3, 123)
	v, ok = m.LoadAndDelete(n + 3)
	require.True(t, ok)
	require.Equal(t, 123, v)
	_, ok = m.Get(n + 3)
	require.False(t, ok)
}

// Race test colocated with concurrent tree map tests.
func TestConcurrentTreeMap_Races(t *testing.T) {
	m := NewConcurrentTreeMapOrdered[int, int]()
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

func TestConcurrentTreeMap_NavigationAndRange(t *testing.T) {
	m := NewConcurrentTreeMapOrdered[int, int]()
	for i := 1; i <= 5; i++ {
		m.Put(i*10, i)
	}
	k, ok := m.FirstKey()
	require.True(t, ok)
	require.Equal(t, 10, k)
	k, ok = m.LastKey()
	require.True(t, ok)
	require.Equal(t, 50, k)
	k, ok = m.FloorKey(35)
	require.True(t, ok)
	require.Equal(t, 30, k)
	k, ok = m.CeilingKey(35)
	require.True(t, ok)
	require.Equal(t, 40, k)
	k, ok = m.LowerKey(30)
	require.True(t, ok)
	require.Equal(t, 20, k)
	k, ok = m.HigherKey(30)
	require.True(t, ok)
	require.Equal(t, 40, k)

	e, ok := m.FloorEntry(35)
	require.True(t, ok)
	require.Equal(t, 30, e.Key)
	e, ok = m.CeilingEntry(35)
	require.True(t, ok)
	require.Equal(t, 40, e.Key)
	e, ok = m.LowerEntry(30)
	require.True(t, ok)
	require.Equal(t, 20, e.Key)
	e, ok = m.HigherEntry(30)
	require.True(t, ok)
	require.Equal(t, 40, e.Key)

	// Range
	rs := make([]int, 0)
	m.Range(15, 45, func(k, v int) bool {
		rs = append(rs, k)
		return true
	})
	require.True(t, slices.Equal(rs, []int{20, 30, 40}), "Range=%v", rs)
	// Rank and GetByRank
	require.Equal(t, 2, m.RankOfKey(30))
	e, ok = m.GetByRank(1)
	require.True(t, ok)
	require.Equal(t, 20, e.Key)
}

func TestConcurrentTreeMap_Views(t *testing.T) {
	m := NewConcurrentTreeMapOrdered[int, int]()
	m.Put(10, 1)
	m.Put(20, 2)
	m.Put(30, 3)
	m.Put(40, 4)
	sub := m.SubMap(15, 35).Keys()
	head := m.HeadMap(25, false).Keys()
	tail := m.TailMap(25, true).Keys()
	require.True(t, slices.Equal(sub, []int{20, 30}), "SubMap=%v", sub)
	require.True(t, slices.Equal(head, []int{10, 20}), "HeadMap=%v", head)
	require.True(t, slices.Equal(tail, []int{30, 40}), "TailMap=%v", tail)
}

func TestConcurrentTreeMap_ClearAndClone(t *testing.T) {
	m := NewConcurrentTreeMapFrom(func(a, b int) int { return cmp.Compare(a, b) }, map[int]string{1: "a", 2: "b", 3: "c"})
	require.False(t, m.IsEmpty())
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

func TestConcurrentTreeMap_PutAllAndPutSeq(t *testing.T) {
	m := NewConcurrentTreeMapOrdered[int, int]()
	m.Put(1, 10)
	m2 := NewConcurrentTreeMapOrdered[int, int]()
	m2.Put(1, 100)
	m2.Put(2, 200)
	m.PutAll(m2)
	require.Equal(t, 2, m.Size())
	v, _ := m.Get(1)
	require.Equal(t, 100, v)
	// PutSeq
	m.PutSeq(func(yield func(int, int) bool) {
		if !yield(3, 300) {
			return
		}
		yield(4, 400)
	})
	require.Equal(t, 4, m.Size())
}

func TestConcurrentTreeMap_ComputeAndMerge(t *testing.T) {
	m := NewConcurrentTreeMapOrdered[int, string]()
	m.Put(1, "a")
	v, ok := m.Compute(1, func(k int, old string, exists bool) (string, bool) {
		if exists {
			return old + "_updated", true
		}
		return "new", true
	})
	require.True(t, ok)
	require.Equal(t, "a_updated", v)
	// ComputeIfAbsent
	v = m.ComputeIfAbsent(2, func(k int) string { return "b" })
	require.Equal(t, "b", v)
	// ComputeIfPresent
	v, ok = m.ComputeIfPresent(1, func(k int, old string) (string, bool) {
		return old + "_again", true
	})
	require.True(t, ok)
	require.Equal(t, "a_updated_again", v)

	// Merge
	v, ok = m.Merge(3, "c", func(old, new string) (string, bool) {
		return old + "_" + new, true
	})
	require.True(t, ok)
	require.Equal(t, "c", v)
}

func TestConcurrentTreeMap_ReplaceOperations(t *testing.T) {
	m := NewConcurrentTreeMapOrdered[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	old, ok := m.Replace(1, "aa")
	require.True(t, ok)
	require.Equal(t, "a", old)
	v, _ := m.Get(1)
	require.Equal(t, "aa", v)
	// ReplaceIf
	require.True(t, m.ReplaceIf(1, "aa", "aaa", eqV[string]))
	require.False(t, m.ReplaceIf(1, "wrong", "xxx", eqV[string]))
	// ReplaceAll
	m.ReplaceAll(func(k int, v string) string { return v + "!" })
	v, _ = m.Get(1)
	require.Equal(t, "aaa!", v)
}

func TestConcurrentTreeMap_PopFirstLast(t *testing.T) {
	m := NewConcurrentTreeMapOrdered[int, int]()
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

func TestConcurrentTreeMap_FilterAndEquals(t *testing.T) {
	m := NewConcurrentTreeMapOrdered[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	m.Put(3, "c")
	filtered := m.Filter(func(k int, v string) bool { return k > 1 })
	require.Equal(t, 2, filtered.Size())
	// Equals
	m2 := NewConcurrentTreeMapOrdered[int, string]()
	m2.Put(1, "a")
	m2.Put(2, "b")
	m2.Put(3, "c")
	require.True(t, m.Equals(m2, eqV[string]))
}

func TestConcurrentTreeMap_CoverageSupplement(t *testing.T) {
	m := NewConcurrentTreeMapOrdered[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")

	// GetOrDefault
	require.Equal(t, "a", m.GetOrDefault(1, "x"))
	require.Equal(t, "x", m.GetOrDefault(99, "x"))

	// RemoveIf
	require.True(t, m.RemoveIf(1, "a", eqV[string]))
	require.False(t, m.RemoveIf(1, "a", eqV[string]))
	require.False(t, m.RemoveIf(2, "wrong", eqV[string]))

	// ContainsKey/Value
	require.True(t, m.ContainsKey(2))
	require.False(t, m.ContainsKey(1))
	require.True(t, m.ContainsValue("b", eqV[string]))
	require.False(t, m.ContainsValue("z", eqV[string]))

	// RemoveKeys
	m.Put(3, "c")
	m.Put(4, "d")
	// keys: 2, 3, 4
	require.Equal(t, 2, m.RemoveKeys(2, 3))
	require.Equal(t, 1, m.Size())

	// RemoveKeysSeq
	m.Put(5, "e")
	removed := m.RemoveKeysSeq(func(yield func(int) bool) {
		if !yield(4) {
			return
		}
		yield(99)
	})
	require.Equal(t, 1, removed)

	// RemoveFunc
	m.Put(6, "f")
	removed = m.RemoveFunc(func(k int, v string) bool { return k == 6 })
	require.Equal(t, 1, removed)

	// Iteration
	m.Put(10, "ten")
	m.Put(20, "twenty")

	// ForEach
	cnt := 0
	m.ForEach(func(k int, v string) bool {
		cnt++
		return true
	})
	require.Equal(t, 3, cnt)

	// SeqKeys, SeqValues
	kc := 0
	for range m.SeqKeys() {
		kc++
	}
	require.Equal(t, 3, kc)

	vc := 0
	for range m.SeqValues() {
		vc++
	}
	require.Equal(t, 3, vc)

	// String
	str := m.String()
	require.Contains(t, str, "concurrentTreeMap")

	// LowerKey, HigherKey
	m.Put(1, "one")
	m.Put(3, "three")
	// keys: 1, 5, 10, 20 (Wait, 5, 10, 20 were remaining? Yes)
	// keys currently: 5, 10, 20. +1, +3 -> 1, 3, 5, 10, 20

	lk, lok := m.LowerKey(4)
	require.True(t, lok)
	require.Equal(t, 3, lk)
	hk, hok := m.HigherKey(4)
	require.True(t, hok)
	require.Equal(t, 5, hk)

	// FloorEntry, CeilingEntry, LowerEntry, HigherEntry
	fe, fok := m.FloorEntry(4)
	require.True(t, fok)
	require.Equal(t, 3, fe.Key)

	// Ascend, Descend
	ascCnt := 0
	m.Ascend(func(k int, v string) bool {
		ascCnt++
		return true
	})
	require.Equal(t, 5, ascCnt) // 1, 3, 5, 10, 20
	descCnt := 0
	m.Descend(func(k int, v string) bool {
		descCnt++
		return true
	})
	require.Equal(t, 5, descCnt)

	// AscendFrom, DescendFrom
	m.AscendFrom(5, func(k int, v string) bool { return true })
	m.DescendFrom(5, func(k int, v string) bool { return true })

	// RangeSeq, RangeFrom, RangeTo
	for range m.RangeSeq(1, 100) {
	}
	m.RangeFrom(1, func(k int, v string) bool { return true })
	m.RangeTo(100, func(k int, v string) bool { return true })

	// Reversed
	for range m.Reversed() {
	}

	// CloneSorted
	cs := m.CloneSorted()
	require.Equal(t, m.Size(), cs.Size())

	// Merge
	m.Put(25, "a")
	// Merge existing
	v, ok := m.Merge(25, "b", func(old, new string) (string, bool) {
		return old + new, true
	})
	require.True(t, ok)
	require.Equal(t, "ab", v)
	// Merge existing -> remove
	_, ok = m.Merge(25, "c", func(old, new string) (string, bool) {
		return "", false
	})
	require.False(t, ok)
	require.False(t, m.ContainsKey(25))
	// Merge absent
	v, ok = m.Merge(30, "new", func(old, new string) (string, bool) {
		return new, true // shouldn't be called
	})
	require.True(t, ok)
	require.Equal(t, "new", v)

	// ComputeIfPresent
	m.Put(10, "ten")
	// Present -> update
	v, ok = m.ComputeIfPresent(10, func(k int, old string) (string, bool) {
		return old + "_updated", true
	})
	require.True(t, ok)
	require.Equal(t, "ten_updated", v)
	// Present -> remove
	_, ok = m.ComputeIfPresent(10, func(k int, old string) (string, bool) {
		return "", false
	})
	require.False(t, ok)
	require.False(t, m.ContainsKey(10))
	// Absent
	_, ok = m.ComputeIfPresent(999, func(k int, old string) (string, bool) {
		return "new", true
	})
	require.False(t, ok)

	// ReplaceIf
	m.Put(40, "val")
	require.True(t, m.ReplaceIf(40, "val", "newVal", eqV[string]))
	require.False(t, m.ReplaceIf(40, "wrong", "xxx", eqV[string]))

	// ReplaceAll
	m.ReplaceAll(func(k int, v string) string {
		return v + "!"
	})
	val, _ := m.Get(40)
	require.Equal(t, "newVal!", val)
	t.Logf("CoverageSupplement finished")
	// require.Fail(t, "Force fail to see log")
}

package collections

import (
	"cmp"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Panic tests colocated with tree map unit tests.
func TestTreeMap_PanicOnNilComparator(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "expected panic on nil comparator")
		}
	}()
	_ = NewTreeMap[int, int](nil)
}

func TestTreeMap_BasicAndOrder(t *testing.T) {
	m := NewTreeMapOrdered[int, string]()
	m.Put(2, "b")
	m.Put(1, "a")
	m.Put(3, "c")
	ks := m.Keys()
	assert.True(t, slices.IsSorted(ks), "Keys not sorted: %v", ks)
	dec := make([]int, 0, 3)
	for k := range m.Reversed() {
		dec = append(dec, k)
	}
	assert.Equal(t, 3, len(dec))
	assert.Equal(t, 3, dec[0])
	assert.Equal(t, 1, dec[2])
}

func TestTreeMap_NavigationAndExtremes(t *testing.T) {
	m := NewTreeMap[int, int](func(a, b int) int { return cmp.Compare(a, b) })
	for i := 1; i <= 5; i++ {
		m.Put(i*10, i)
	}
	k, ok := m.FirstKey()
	require.True(t, ok)
	assert.Equal(t, 10, k)
	k, ok = m.LastKey()
	require.True(t, ok)
	assert.Equal(t, 50, k)
	e, ok := m.FirstEntry()
	require.True(t, ok)
	assert.Equal(t, 10, e.Key)
	assert.Equal(t, 1, e.Value)
	e, ok = m.LastEntry()
	require.True(t, ok)
	assert.Equal(t, 50, e.Key)
	assert.Equal(t, 5, e.Value)
	e, ok = m.PopFirst()
	require.True(t, ok)
	assert.Equal(t, 10, e.Key)
	e, ok = m.PopLast()
	require.True(t, ok)
	assert.Equal(t, 50, e.Key)
	k, ok = m.FloorKey(33)
	require.True(t, ok)
	assert.Equal(t, 30, k)
	k, ok = m.CeilingKey(33)
	require.True(t, ok)
	assert.Equal(t, 40, k)
	_, ok = m.LowerKey(20)
	require.False(t, ok, "LowerKey(20) should be none after PopFirst")
	_, ok = m.HigherKey(40)
	require.False(t, ok, "HigherKey(40) should be none after PopLast")
}

func TestTreeMap_RangeRankAndViews(t *testing.T) {
	m := NewTreeMapOrdered[int, int]()
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
	assert.True(t, slices.Equal(rs, []int{20, 30}), "Range=%v", rs)
	// Rank
	assert.Equal(t, 2, m.RankOfKey(30))
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

func TestTreeMap_AscendDescend(t *testing.T) {
	m := NewTreeMapOrdered[int, int]()
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
	assert.Equal(t, []int{1, 2, 3, 4, 5}, asc)
	// Descend
	desc := make([]int, 0)
	m.Descend(func(k, v int) bool {
		desc = append(desc, k)
		return true
	})
	assert.Equal(t, []int{5, 4, 3, 2, 1}, desc)
	// AscendFrom
	af := make([]int, 0)
	m.AscendFrom(3, func(k, v int) bool {
		af = append(af, k)
		return true
	})
	assert.Equal(t, []int{3, 4, 5}, af)
	// DescendFrom
	df := make([]int, 0)
	m.DescendFrom(3, func(k, v int) bool {
		df = append(df, k)
		return true
	})
	assert.Equal(t, []int{3, 2, 1}, df)
	// Test early termination
	count := 0
	m.Ascend(func(k, v int) bool {
		count++
		return count < 3
	})
	assert.Equal(t, 3, count)
}

func TestTreeMap_RangeFromAndTo(t *testing.T) {
	m := NewTreeMapOrdered[int, int]()
	for i := 1; i <= 5; i++ {
		m.Put(i, i*10)
	}
	// RangeFrom
	rf := make([]int, 0)
	m.RangeFrom(3, func(k, v int) bool {
		rf = append(rf, k)
		return true
	})
	assert.Equal(t, []int{3, 4, 5}, rf)
	// RangeTo
	rt := make([]int, 0)
	m.RangeTo(3, func(k, v int) bool {
		rt = append(rt, k)
		return true
	})
	assert.Equal(t, []int{1, 2, 3}, rt)
}

func TestTreeMap_RangeSeq(t *testing.T) {
	m := NewTreeMapOrdered[int, int]()
	m.Put(10, 1)
	m.Put(20, 2)
	m.Put(30, 3)
	m.Put(40, 4)
	m.Put(50, 5)
	rs := make([]int, 0)
	for k := range m.RangeSeq(15, 45) {
		rs = append(rs, k)
	}
	assert.Equal(t, []int{20, 30, 40}, rs)
}

func TestTreeMap_EntryNavigation(t *testing.T) {
	m := NewTreeMapOrdered[int, int]()
	for i := 1; i <= 5; i++ {
		m.Put(i*10, i)
	}
	e, ok := m.FloorEntry(35)
	require.True(t, ok)
	assert.Equal(t, 30, e.Key)
	e, ok = m.CeilingEntry(35)
	require.True(t, ok)
	assert.Equal(t, 40, e.Key)
	e, ok = m.LowerEntry(25)
	require.True(t, ok)
	assert.Equal(t, 20, e.Key)
	e, ok = m.HigherEntry(25)
	require.True(t, ok)
	assert.Equal(t, 30, e.Key)
}

func TestTreeMap_CloneAndFilter(t *testing.T) {
	m := NewTreeMapOrdered[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	m.Put(3, "c")
	// Clone
	clone := m.Clone()
	clone.Put(4, "d")
	assert.Equal(t, 3, m.Size())
	assert.Equal(t, 4, clone.Size())
	// CloneSorted
	sortedClone := m.CloneSorted()
	assert.Equal(t, 3, sortedClone.Size())
	// Filter
	filtered := m.Filter(func(k int, v string) bool { return k > 1 })
	assert.Equal(t, 2, filtered.Size())
}

func TestTreeMap_Equals(t *testing.T) {
	m := NewTreeMapOrdered[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	m.Put(3, "c")
	m2 := NewTreeMapOrdered[int, string]()
	m2.Put(1, "a")
	m2.Put(2, "b")
	m2.Put(3, "c")
	assert.True(t, m.Equals(m2, eqV[string]))
	m3 := NewTreeMapOrdered[int, string]()
	m3.Put(1, "a")
	m3.Put(2, "b")
	m3.Put(3, "different")
	assert.False(t, m.Equals(m3, eqV[string]))
}

func TestTreeMap_Clear(t *testing.T) {
	m := NewTreeMapOrdered[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	m.Put(3, "c")
	assert.False(t, m.IsEmpty())
	m.Clear()
	assert.True(t, m.IsEmpty())
	assert.Equal(t, 0, m.Size())
}

func TestTreeMap_PutAllAndPutSeq(t *testing.T) {
	m := NewTreeMapOrdered[int, int]()
	m2 := NewTreeMapOrdered[int, int]()
	m2.Put(1, 10)
	m2.Put(2, 20)
	m.PutAll(m2)
	assert.Equal(t, 2, m.Size())
	// PutSeq
	m.PutSeq(func(yield func(int, int) bool) {
		if !yield(3, 30) {
			return
		}
		yield(4, 40)
	})
	assert.Equal(t, 4, m.Size())
}

func TestTreeMap_RemoveFuncAndCompute(t *testing.T) {
	m := NewTreeMapOrdered[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	m.Put(3, "c")
	// RemoveFunc
	count := m.RemoveFunc(func(k int, v string) bool { return k > 1 })
	assert.Equal(t, 2, count)
	assert.Equal(t, 1, m.Size())
	// Compute
	m.Put(4, "d")
	v, ok := m.Compute(4, func(k int, old string, exists bool) (string, bool) {
		if exists {
			return old + "_updated", true
		}
		return "new", true
	})
	require.True(t, ok)
	assert.Equal(t, "d_updated", v)
}

func TestTreeMap_CoverageSupplement(t *testing.T) {
	// NewTreeMapFrom
	m := NewTreeMapFrom(func(a, b int) int { return cmp.Compare(a, b) }, map[int]string{1: "a", 2: "b"})
	assert.Equal(t, 2, m.Size())
	assert.True(t, m.ContainsKey(1))

	// String
	str := m.String()
	assert.Contains(t, str, "treeMap")
	assert.Contains(t, str, "1:a")

	// GetOrDefault
	assert.Equal(t, "b", m.GetOrDefault(2, "x"))
	assert.Equal(t, "x", m.GetOrDefault(99, "x"))

	// RemoveIf
	assert.True(t, m.RemoveIf(1, "a", eqV[string]))
	assert.False(t, m.RemoveIf(1, "a", eqV[string])) // already removed
	assert.False(t, m.RemoveIf(2, "wrong", eqV[string]))

	// ContainsValue
	m.Put(3, "c")
	assert.True(t, m.ContainsValue("c", eqV[string]))
	assert.False(t, m.ContainsValue("z", eqV[string]))

	// RemoveKeys, RemoveKeysSeq
	m.PutAll(NewTreeMapFrom(cmp.Compare[int], map[int]string{4: "d", 5: "e", 6: "f"}))
	// keys: 2, 3, 4, 5, 6
	removed := m.RemoveKeys(2, 3)
	assert.Equal(t, 2, removed)
	assert.Equal(t, 3, m.Size()) // 4, 5, 6 remaining
	removed = m.RemoveKeysSeq(func(yield func(int) bool) {
		if !yield(4) {
			return
		}
		yield(99) // not present
	})
	assert.Equal(t, 1, removed)
	assert.Equal(t, 2, m.Size()) // 5, 6 remaining

	// Values, Keys, Entries
	ks := m.Keys()
	vs := m.Values()
	es := m.Entries()
	assert.Equal(t, 2, len(ks))
	assert.Equal(t, 2, len(vs))
	assert.Equal(t, 2, len(es))

	// ForEach
	cnt := 0
	m.ForEach(func(k int, v string) bool {
		cnt++
		return true
	})
	assert.Equal(t, 2, cnt)

	// SeqKeys, SeqValues
	kSeq := make([]int, 0)
	for k := range m.SeqKeys() {
		kSeq = append(kSeq, k)
	}
	assert.Len(t, kSeq, 2)
	vSeq := make([]string, 0)
	for v := range m.SeqValues() {
		vSeq = append(vSeq, v)
	}
	assert.Len(t, vSeq, 2)

	// ComputeIfPresent
	m.Put(10, "ten")
	// Present -> update
	v, ok := m.ComputeIfPresent(10, func(k int, old string) (string, bool) {
		return old + "_updated", true
	})
	assert.True(t, ok)
	assert.Equal(t, "ten_updated", v)
	// Present -> remove
	_, ok = m.ComputeIfPresent(10, func(k int, old string) (string, bool) {
		return "", false
	})
	assert.False(t, ok)
	assert.False(t, m.ContainsKey(10))
	// Absent
	_, ok = m.ComputeIfPresent(999, func(k int, old string) (string, bool) {
		return "new", true
	})
	assert.False(t, ok)

	// Merge
	m.Put(20, "a")
	// Merge existing
	v, ok = m.Merge(20, "b", func(old, new string) (string, bool) {
		return old + new, true
	})
	assert.True(t, ok)
	assert.Equal(t, "ab", v)
	// Merge existing -> remove
	_, ok = m.Merge(20, "c", func(old, new string) (string, bool) {
		return "", false
	})
	assert.False(t, ok)
	assert.False(t, m.ContainsKey(20))
	// Merge absent
	v, ok = m.Merge(30, "new", func(old, new string) (string, bool) {
		return new, true // shouldn't be called
	})
	assert.True(t, ok)
	assert.Equal(t, "new", v)

	// ReplaceIf
	m.Put(40, "val")
	assert.True(t, m.ReplaceIf(40, "val", "newVal", eqV[string]))
	assert.False(t, m.ReplaceIf(40, "wrong", "xxx", eqV[string]))

	// ReplaceAll
	m.ReplaceAll(func(k int, v string) string {
		return v + "!"
	})
	val, _ := m.Get(40)
	assert.Equal(t, "newVal!", val)
}

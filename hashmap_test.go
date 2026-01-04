package collections

import (
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

// helper to make a Seq2 from slices of keys/values (same length)
func seq2Of[K any, V any](keys []K, values []V) func(func(K, V) bool) {
	return func(yield func(K, V) bool) {
		for i := range keys {
			if !yield(keys[i], values[i]) {
				return
			}
		}
	}
}

func eqV[T comparable](a, b T) bool { return a == b }

func TestHashMap_BasicCRUD(t *testing.T) {
	m := NewHashMap[int, string]()
	require.True(t, m.IsEmpty())
	require.Equal(t, 0, m.Size())
	_, ok := m.Get(1)
	require.False(t, ok, "Unexpected key should not be present")
	_, ok = m.Remove(1)
	require.False(t, ok)
	require.Equal(t, "x", m.GetOrDefault(1, "x"))
	old, replaced := m.Put(1, "a")
	require.False(t, replaced)
	require.Equal(t, "", old)
	old, replaced = m.Put(1, "b")
	require.True(t, replaced)
	require.Equal(t, "a", old)
	val, ok := m.Get(1)
	require.True(t, ok)
	require.Equal(t, "b", val)
	// PutIfAbsent
	v, inserted := m.PutIfAbsent(1, "c")
	require.False(t, inserted)
	require.Equal(t, "b", v)
	v, inserted = m.PutIfAbsent(2, "x")
	require.True(t, inserted)
	require.Equal(t, "x", v)
	// RemoveIf
	require.True(t, m.RemoveIf(2, "x", eqV[string]))
	require.False(t, m.RemoveIf(2, "x", eqV[string]))
}

func TestHashMap_PutSeq_RemoveKeysSeq(t *testing.T) {
	m := NewHashMap[int, int]()
	changed := m.PutSeq(seq2Of([]int{1, 2, 3, 2}, []int{10, 20, 30, 200}))
	// Keys touched: 1,2,3 -> 3 unique
	require.Equal(t, 3, changed)
	require.Equal(t, 3, m.Size())
	removed := m.RemoveKeysSeq(seqOf([]int{2, 4}))
	require.Equal(t, 1, removed)
	require.Equal(t, 2, m.Size())
}

func TestHashMap_Contains_Views_Iter(t *testing.T) {
	m := NewHashMap[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	require.True(t, m.ContainsKey(1))
	require.False(t, m.ContainsKey(3))
	require.True(t, m.ContainsValue("b", eqV[string]))
	require.False(t, m.ContainsValue("x", eqV[string]))
	ks := m.Keys()
	slices.Sort(ks)
	require.True(t, slices.Equal(ks, []int{1, 2}), "Keys=%v", ks)
	vs := m.Values()
	slices.Sort(vs)
	require.True(t, slices.Equal(vs, []string{"a", "b"}), "Values=%v", vs)
	es := m.Entries()
	sort.Slice(es, func(i, j int) bool { return es[i].Key < es[j].Key })
	require.Len(t, es, 2)
	require.Equal(t, 1, es[0].Key)
	require.Equal(t, 2, es[1].Key)
	// ForEach early stop
	cnt := 0
	m.ForEach(func(k int, v string) bool {
		cnt++
		return false // stop immediately
	})
	require.Equal(t, 1, cnt)
	// Seq/SeqKeys/SeqValues
	gotK := make([]int, 0, 2)
	for k := range m.SeqKeys() {
		gotK = append(gotK, k)
	}
	slices.Sort(gotK)
	require.True(t, slices.Equal(gotK, []int{1, 2}), "SeqKeys=%v", gotK)
	gotV := make([]string, 0, 2)
	for _, v := range m.Seq() {
		gotV = append(gotV, v)
	}
	slices.Sort(gotV)
	require.True(t, slices.Equal(gotV, []string{"a", "b"}), "Seq Values should match, got=%v", gotV)
	gotValues := make([]string, 0, 2)
	for v := range m.SeqValues() {
		gotValues = append(gotValues, v)
	}
	slices.Sort(gotValues)
	require.True(t, slices.Equal(gotValues, []string{"a", "b"}), "SeqValues=%v", gotValues)
}

func TestHashMap_ComputeVariants(t *testing.T) {
	m := NewHashMap[int, int]()
	// ComputeIfAbsent
	v := m.ComputeIfAbsent(1, func(k int) int { return k * 10 })
	require.Equal(t, 10, v)
	require.Equal(t, 1, m.Size())
	// Compute existing -> update
	nv, ok := m.Compute(1, func(k, old int, exists bool) (int, bool) {
		require.True(t, exists, "Exists expected true")
		return old + 1, true
	})
	require.True(t, ok)
	require.Equal(t, 11, nv)
	// Compute remove
	_, ok = m.Compute(1, func(k, old int, exists bool) (int, bool) {
		return 0, false
	})
	require.False(t, ok)
	require.False(t, m.ContainsKey(1))
	// ComputeIfPresent on absent
	_, ok = m.ComputeIfPresent(1, func(k, old int) (int, bool) { return 0, true })
	require.False(t, ok, "ComputeIfPresent should be false on absent")
	// Merge
	m.Put(2, 5)
	nv, ok = m.Merge(2, 7, func(old, new int) (int, bool) { return old + new, true })
	require.True(t, ok)
	require.Equal(t, 12, nv)
	_, ok = m.Merge(2, 0, func(old, new int) (int, bool) { return 0, false })
	require.False(t, ok)
	require.False(t, m.ContainsKey(2))
	// Merge on absent
	v, ok = m.Merge(3, 30, func(old, new int) (int, bool) { return old + new, true })
	require.True(t, ok)
	require.Equal(t, 30, v)
}

func TestHashMap_ReplaceAndFilterAndEquals(t *testing.T) {
	m := NewHashMap[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	old, replaced := m.Replace(1, "aa")
	require.True(t, replaced)
	require.Equal(t, "a", old)
	require.False(t, m.ReplaceIf(2, "x", "bb", eqV[string]), "ReplaceIf should fail on mismatched oldValue")
	require.True(t, m.ReplaceIf(2, "b", "bb", eqV[string]))
	// ReplaceAll
	m.ReplaceAll(func(k int, v string) string { return v + "!" })
	require.True(t, m.ContainsValue("aa!", eqV[string]))
	require.True(t, m.ContainsValue("bb!", eqV[string]))
	// Clone
	cp := m.Clone()
	require.True(t, m.Equals(cp, eqV[string]))
	require.True(t, cp.Equals(m, eqV[string]))
	// Filter
	f := m.Filter(func(k int, v string) bool { return k == 1 })
	require.Equal(t, 1, f.Size())
	require.True(t, f.ContainsKey(1))
	// ToGoMap is a snapshot (available via GoMapView)
	gm := m.(GoMapView[int, string]).ToGoMap()
	gm[3] = "zzz"
	require.False(t, m.ContainsKey(3), "ToGoMap must return snapshot, not live map")
}

func TestHashMap_EmptySingle(t *testing.T) {
	m := NewHashMap[int, string]()
	require.True(t, m.IsEmpty())
	require.Equal(t, 0, m.Size())
	require.Equal(t, "x", m.GetOrDefault(1, "x"))
	_, ok := m.Remove(1)
	require.False(t, ok)
	old, replaced := m.Put(1, "a")
	require.False(t, replaced)
	require.Equal(t, "", old)
	v, ok := m.Get(1)
	require.True(t, ok)
	require.Equal(t, "a", v)
}

func TestHashMap_Clear(t *testing.T) {
	m := NewHashMapFrom(map[int]string{1: "a", 2: "b", 3: "c"})
	require.False(t, m.IsEmpty())
	m.Clear()
	require.True(t, m.IsEmpty())
	require.Equal(t, 0, m.Size())
}

func TestHashMap_RemoveKeys(t *testing.T) {
	m := NewHashMapFrom(map[int]string{1: "a", 2: "b", 3: "c"})
	removed := m.RemoveKeys(1, 2)
	require.Equal(t, 2, removed)
	require.Equal(t, 1, m.Size())
	require.False(t, m.ContainsKey(1))
	require.False(t, m.ContainsKey(2))
	require.True(t, m.ContainsKey(3))
}

func TestHashMap_PutAll(t *testing.T) {
	m := NewHashMap[int, string]()
	m.Put(1, "a")
	m2 := NewHashMapFrom(map[int]string{1: "x", 2: "y"})
	m.PutAll(m2)
	require.Equal(t, 2, m.Size())
	v, _ := m.Get(1)
	require.Equal(t, "x", v)
	v, _ = m.Get(2)
	require.Equal(t, "y", v)
}

func TestHashMap_RemoveFunc(t *testing.T) {
	m := NewHashMapFrom(map[int]string{1: "a", 2: "b", 3: "c"})
	count := m.RemoveFunc(func(k int, v string) bool { return k > 1 })
	require.Equal(t, 2, count)
	require.Equal(t, 1, m.Size())
	require.True(t, m.ContainsKey(1))
	require.False(t, m.ContainsKey(2))
	require.False(t, m.ContainsKey(3))
}

func TestHashMap_From(t *testing.T) {
	m := NewHashMapFrom(map[int]string{1: "a", 2: "b"})
	require.Equal(t, 2, m.Size())
	v, ok := m.Get(1)
	require.True(t, ok)
	require.Equal(t, "a", v)
}

func TestHashMap_NewWithCapacity(t *testing.T) {
	m := NewHashMapWithCapacity[int, string](10)
	require.True(t, m.IsEmpty())
	m.Put(1, "a")
	require.Equal(t, 1, m.Size())
}

func TestHashMap_String(t *testing.T) {
	m := NewHashMap[int, string]()
	m.Put(1, "a")
	s := m.String()
	require.Contains(t, s, "hashMap")
	require.Contains(t, s, "1:a")
}

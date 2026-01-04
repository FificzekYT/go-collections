package collections

import (
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConcurrentHashMap_Basic(t *testing.T) {
	m := NewConcurrentHashMap[int, string]()
	assert.True(t, m.IsEmpty(), "New map should be empty")
	_, ok := m.Get(1)
	assert.False(t, ok, "Unexpected key should not be present")
	old, replaced := m.Put(1, "a")
	require.False(t, replaced)
	assert.Equal(t, "", old)
	v, ok := m.Get(1)
	require.True(t, ok)
	assert.Equal(t, "a", v)
	v, inserted := m.PutIfAbsent(1, "x")
	require.False(t, inserted)
	assert.Equal(t, "a", v)
	v, computed := m.GetOrCompute(2, func() string { return "b" })
	assert.True(t, computed)
	assert.Equal(t, "b", v)
	v, loaded := m.LoadOrStore(2, "bb")
	assert.True(t, loaded)
	assert.Equal(t, "b", v)
	v, ok = m.LoadAndDelete(2)
	require.True(t, ok)
	assert.Equal(t, "b", v)
}

func TestConcurrentHashMap_CompareAndOps(t *testing.T) {
	m := NewConcurrentHashMap[int, int]()
	m.Put(1, 10)
	// CAS success
	require.True(t, m.CompareAndSwap(1, 10, 11, eqV[int]))
	// CAS fail
	require.False(t, m.CompareAndSwap(1, 10, 12, eqV[int]))
	// CompareAndDelete success
	require.True(t, m.CompareAndDelete(1, 11, eqV[int]))
	_, ok := m.Get(1)
	assert.False(t, ok, "Key 1 should be deleted")
	// CompareAndDelete fail
	m.Put(2, 20)
	require.False(t, m.CompareAndDelete(2, 21, eqV[int]))
	require.True(t, m.ContainsKey(2))
}

func TestConcurrentHashMap_ConcurrentPutIfAbsent(t *testing.T) {
	m := NewConcurrentHashMap[int, int]()
	n := 1000
	workers := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := range workers {
		go func(id int) {
			defer wg.Done()
			for i := range n {
				m.PutIfAbsent(i, i)
			}
		}(w)
	}
	wg.Wait()
	// Validate presence
	for _, k := range []int{0, n / 2, n - 1} {
		v, ok := m.Get(k)
		require.Truef(t, ok, "Missing key %d", k)
		assert.Equal(t, k, v)
	}
}

func TestConcurrentHashMap_ClearAndString(t *testing.T) {
	m := NewConcurrentHashMapFrom(map[int]string{1: "a", 2: "b"})
	assert.False(t, m.IsEmpty())

	// String
	str := m.String()
	assert.Contains(t, str, "concurrentHashMap")
	assert.Contains(t, str, "1:a")
	assert.Contains(t, str, "2:b")

	m.Clear()
	assert.True(t, m.IsEmpty())
	assert.Equal(t, 0, m.Size())
}

func TestConcurrentHashMap_PutAllAndPutSeq(t *testing.T) {
	m := NewConcurrentHashMap[int, int]()
	other := NewConcurrentHashMap[int, int]()
	other.Put(1, 10)
	other.Put(2, 20)
	m.PutAll(other)
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

func TestConcurrentHashMap_GetOrDefaultAndRemove(t *testing.T) {
	m := NewConcurrentHashMap[int, string]()
	m.Put(1, "a")
	assert.Equal(t, "a", m.GetOrDefault(1, "default"))
	assert.Equal(t, "default", m.GetOrDefault(999, "default"))

	v, ok := m.Remove(1)
	require.True(t, ok)
	assert.Equal(t, "a", v)
	_, ok = m.Get(1)
	assert.False(t, ok)
}

func TestConcurrentHashMap_ContainsAndRemoveKeys(t *testing.T) {
	m := NewConcurrentHashMap[int, string]()
	m.Put(1, "a")
	m.Put(2, "b")
	m.Put(3, "c")

	assert.True(t, m.ContainsKey(1))
	assert.False(t, m.ContainsKey(4))

	assert.True(t, m.ContainsValue("a", eqV[string]))
	assert.False(t, m.ContainsValue("z", eqV[string]))

	// RemoveIf
	removed := m.RemoveIf(1, "a", eqV[string])
	assert.True(t, removed)
	assert.False(t, m.ContainsKey(1))

	removed = m.RemoveIf(2, "xxx", eqV[string])
	assert.False(t, removed)
	assert.True(t, m.ContainsKey(2))

	// RemoveKeys
	count := m.RemoveKeys(2, 3, 99)
	assert.Equal(t, 2, count)
	assert.True(t, m.IsEmpty())

	// RemoveKeysSeq
	m.Put(4, "d")
	m.Put(5, "e")
	count = m.RemoveKeysSeq(func(yield func(int) bool) {
		if !yield(4) {
			return
		}
		if !yield(5) {
			return
		}
		yield(100)
	})
	assert.Equal(t, 2, count)
	assert.True(t, m.IsEmpty())

	// RemoveFunc
	m.Put(1, "a")
	m.Put(2, "b")
	m.Put(3, "c")
	count = m.RemoveFunc(func(k int, v string) bool {
		return k%2 != 0 // remove odd keys: 1, 3
	})
	assert.Equal(t, 2, count)
	assert.Equal(t, 1, m.Size())
	assert.True(t, m.ContainsKey(2))
}

func TestConcurrentHashMap_ComputeAndMerge(t *testing.T) {
	m := NewConcurrentHashMap[int, int]()

	// Compute
	v, ok := m.Compute(1, func(k int, old int, exists bool) (int, bool) {
		assert.False(t, exists)
		return 10, true
	})
	require.True(t, ok)
	assert.Equal(t, 10, v)

	v, ok = m.Compute(1, func(k int, old int, exists bool) (int, bool) {
		assert.True(t, exists)
		return old + 5, true
	})
	require.True(t, ok)
	assert.Equal(t, 15, v)

	// Compute: delete
	_, ok = m.Compute(1, func(k int, old int, exists bool) (int, bool) {
		return 0, false
	})
	assert.False(t, ok)
	assert.False(t, m.ContainsKey(1))

	// ComputeIfAbsent
	v = m.ComputeIfAbsent(2, func(k int) int { return 20 })
	assert.Equal(t, 20, v)
	v = m.ComputeIfAbsent(2, func(k int) int { return 999 }) // should not update
	assert.Equal(t, 20, v)

	// ComputeIfPresent
	v, ok = m.ComputeIfPresent(2, func(k int, old int) (int, bool) {
		return old * 2, true
	})
	require.True(t, ok)
	assert.Equal(t, 40, v)

	_, ok = m.ComputeIfPresent(3, func(k int, old int) (int, bool) { return 0, true })
	assert.False(t, ok)

	// ComputeIfPresent: delete
	_, ok = m.ComputeIfPresent(2, func(k int, old int) (int, bool) {
		return 0, false
	})
	assert.False(t, ok)
	assert.False(t, m.ContainsKey(2))

	// Merge
	m.Put(1, 10)
	v, ok = m.Merge(1, 20, func(old, newV int) (int, bool) {
		return old + newV, true
	})
	require.True(t, ok)
	assert.Equal(t, 30, v)

	v, ok = m.Merge(5, 50, func(old, newV int) (int, bool) { return 0, true }) // key absent
	require.True(t, ok)
	assert.Equal(t, 50, v)

	// Merge: delete
	_, ok = m.Merge(1, 999, func(old, newV int) (int, bool) {
		return 0, false
	})
	assert.False(t, ok)
	assert.False(t, m.ContainsKey(1))
}

func TestConcurrentHashMap_ReplaceOperations(t *testing.T) {
	m := NewConcurrentHashMap[int, string]()
	m.Put(1, "a")

	// Replace
	old, ok := m.Replace(1, "b")
	require.True(t, ok)
	assert.Equal(t, "a", old)
	assert.Equal(t, "b", m.GetOrDefault(1, ""))

	_, ok = m.Replace(99, "z")
	require.False(t, ok)

	// ReplaceIf
	ok = m.ReplaceIf(1, "b", "c", eqV[string])
	require.True(t, ok)
	assert.Equal(t, "c", m.GetOrDefault(1, ""))

	ok = m.ReplaceIf(1, "x", "y", eqV[string])
	require.False(t, ok)

	// ReplaceAll
	m.Put(2, "hello")
	m.ReplaceAll(func(k int, v string) string {
		return strings.ToUpper(v)
	})
	assert.Equal(t, "C", m.GetOrDefault(1, ""))
	assert.Equal(t, "HELLO", m.GetOrDefault(2, ""))
}

func TestConcurrentHashMap_ViewsAndIterations(t *testing.T) {
	m := NewConcurrentHashMap[int, int]()
	m.Put(1, 10)
	m.Put(2, 20)
	m.Put(3, 30)

	// Keys
	keys := m.Keys()
	assert.Len(t, keys, 3)
	assert.Contains(t, keys, 1)

	// Values
	values := m.Values()
	assert.Len(t, values, 3)
	assert.Contains(t, values, 10)

	// Entries
	entries := m.Entries()
	assert.Len(t, entries, 3)

	// ForEach
	count := 0
	m.ForEach(func(k, v int) bool {
		count++
		return true
	})
	assert.Equal(t, 3, count)

	// Seq
	count = 0
	for k, v := range m.Seq() {
		assert.Equal(t, k*10, v)
		count++
	}
	assert.Equal(t, 3, count)

	// SeqKeys
	count = 0
	for range m.SeqKeys() {
		count++
	}
	assert.Equal(t, 3, count)

	// SeqValues
	count = 0
	for range m.SeqValues() {
		count++
	}
	assert.Equal(t, 3, count)
}

func TestConcurrentHashMap_CloneFilterEquals(t *testing.T) {
	m := NewConcurrentHashMap[int, int]()
	m.Put(1, 1)
	m.Put(2, 2)
	m.Put(3, 3)

	// Clone
	c := m.Clone()
	assert.Equal(t, m.Size(), c.Size())
	assert.True(t, m.Equals(c, eqV[int]))

	// Filter
	even := m.Filter(func(k, v int) bool {
		return k%2 == 0
	})
	assert.Equal(t, 1, even.Size())
	assert.True(t, even.ContainsKey(2))

	// Equals
	m2 := NewHashMap[int, int]()
	m2.Put(1, 1)
	m2.Put(2, 2)
	m2.Put(3, 3)
	assert.True(t, m.Equals(m2, eqV[int]))

	m2.Put(4, 4)
	assert.False(t, m.Equals(m2, eqV[int]))
}

func TestConcurrentHashMap_CoverageSupplement(t *testing.T) {
	// Cover cases where ForEach/Range stops early
	m := NewConcurrentHashMap[int, int]()
	for i := range 10 {
		m.Put(i, i)
	}

	count := 0
	m.ForEach(func(k, v int) bool {
		count++
		return count < 5
	})
	assert.Equal(t, 5, count)

	// ContainsValue short-circuit
	assert.True(t, m.ContainsValue(0, eqV[int]))
}

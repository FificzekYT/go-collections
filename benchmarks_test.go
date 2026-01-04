package collections

import (
	"strconv"
	"testing"
)

func BenchmarkHashSet_AddContains(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 5e4} {
		b.Run("n="+strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			s := NewHashSet[int]()
			vals := make([]int, n)
			for i := range n {
				vals[i] = i
			}
			b.ResetTimer()
			for b.Loop() {
				for _, v := range vals {
					s.Add(v)
				}
				for _, v := range vals {
					if !s.Contains(v) {
						b.Fatal("missing")
					}
				}
			}
		})
	}
}

func BenchmarkTreeMap_PutGet(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 5e4} {
		b.Run("n="+strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			m := NewTreeMapOrdered[int, int]()
			keys := make([]int, n)
			for i := range n {
				keys[i] = i
			}
			b.ResetTimer()
			for b.Loop() {
				for _, k := range keys {
					m.Put(k, k)
				}
				for _, k := range keys {
					if v, ok := m.Get(k); !ok || v != k {
						b.Fatal("missing")
					}
				}
			}
		})
	}
}

// Additional Set benchmarks consolidated here to keep benchmark tests in one place.
func BenchmarkSet_AddRemoveContainsIter_HashSet(b *testing.B) {
	benchSetAddRemoveContainsIter(b, "HashSet", func() Set[int] { return NewHashSet[int]() })
}

func BenchmarkSet_AddRemoveContainsIter_TreeSet(b *testing.B) {
	benchSetAddRemoveContainsIter(b, "TreeSet", func() Set[int] { return NewTreeSetOrdered[int]() })
}

func benchSetAddRemoveContainsIter(b *testing.B, name string, newSet func() Set[int]) {
	for _, n := range []int{1e3, 1e4} {
		b.Run(name+"-n="+strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			vals := make([]int, n)
			for i := range n {
				vals[i] = i
			}
			b.ResetTimer()
			for b.Loop() {
				s := newSet()
				for _, v := range vals {
					s.Add(v)
				}
				for _, v := range vals {
					if !s.Contains(v) {
						b.Fatal("contains failed")
					}
				}
				for range s.Seq() {
					// iterate
				}
				for _, v := range vals {
					s.Remove(v)
				}
			}
		})
	}
}

package collections

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ==========================
// 1. HashSet Serialization Tests
// ==========================

type HashSetSerializationTestSuite struct {
	suite.Suite
}

func (suite *HashSetSerializationTestSuite) TestEmptyHashSet() {
	suite.Run("JSON", func() {
		original := NewHashSet[int]()
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewHashSet[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})

	suite.Run("Gob", func() {
		original := NewHashSet[int]()
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewHashSet[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})
}

func (suite *HashSetSerializationTestSuite) TestSingleElement() {
	suite.Run("JSON", func() {
		original := NewHashSetFrom(42)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewHashSet[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 1, restored.Size(), "Size should be one")
		assert.True(suite.T(), restored.Contains(42), "Should contain element")
	})

	suite.Run("Gob", func() {
		original := NewHashSetFrom(42)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewHashSet[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 1, restored.Size(), "Size should be one")
		assert.True(suite.T(), restored.Contains(42), "Should contain element")
	})
}

func (suite *HashSetSerializationTestSuite) TestMultipleElements() {
	suite.Run("JSON", func() {
		original := NewHashSetFrom(1, 2, 3, 4, 5)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewHashSet[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.True(suite.T(), restored.ContainsAll(1, 2, 3, 4, 5), "Should contain all elements")
	})

	suite.Run("Gob", func() {
		original := NewHashSetFrom(1, 2, 3, 4, 5)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewHashSet[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.True(suite.T(), restored.ContainsAll(1, 2, 3, 4, 5), "Should contain all elements")
	})
}

func (suite *HashSetSerializationTestSuite) TestRoundTrip() {
	suite.Run("JSON", func() {
		original := NewHashSetFrom(-10, 0, 42, 100, 999)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewHashSet[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.True(suite.T(), original.Equals(restored), "Should be equal after round-trip")
	})

	suite.Run("Gob", func() {
		original := NewHashSetFrom(-10, 0, 42, 100, 999)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewHashSet[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.True(suite.T(), original.Equals(restored), "Should be equal after round-trip")
	})
}

func TestHashSetSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(HashSetSerializationTestSuite))
}

// ==========================
// 2. ArrayList Serialization Tests
// ==========================

type ArrayListSerializationTestSuite struct {
	suite.Suite
}

func (suite *ArrayListSerializationTestSuite) TestEmptyList() {
	suite.Run("JSON", func() {
		original := NewArrayList[int]()
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewArrayList[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})

	suite.Run("Gob", func() {
		original := NewArrayList[int]()
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewArrayList[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})
}

func (suite *ArrayListSerializationTestSuite) TestSingleElement() {
	suite.Run("JSON", func() {
		original := NewArrayListFrom(42)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewArrayList[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 1, restored.Size(), "Size should be one")
		val, ok := restored.Get(0)
		require.True(suite.T(), ok, "Should retrieve element")
		assert.Equal(suite.T(), 42, val, "Value should match")
	})

	suite.Run("Gob", func() {
		original := NewArrayListFrom(42)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewArrayList[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 1, restored.Size(), "Size should be one")
		val, ok := restored.Get(0)
		require.True(suite.T(), ok, "Should retrieve element")
		assert.Equal(suite.T(), 42, val, "Value should match")
	})
}

func (suite *ArrayListSerializationTestSuite) TestMultipleElements() {
	suite.Run("JSON", func() {
		original := NewArrayListFrom(1, 2, 3, 4, 5)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewArrayList[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		for i := range original.Size() {
			origVal, _ := original.Get(i)
			restVal, _ := restored.Get(i)
			assert.Equal(suite.T(), origVal, restVal, "Element at index %d should match", i)
		}
	})

	suite.Run("Gob", func() {
		original := NewArrayListFrom(1, 2, 3, 4, 5)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewArrayList[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		for i := range original.Size() {
			origVal, _ := original.Get(i)
			restVal, _ := restored.Get(i)
			assert.Equal(suite.T(), origVal, restVal, "Element at index %d should match", i)
		}
	})
}

func (suite *ArrayListSerializationTestSuite) TestRoundTripPreservesOrder() {
	suite.Run("JSON", func() {
		original := NewArrayListFrom(5, 3, 1, 4, 2)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewArrayList[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})

	suite.Run("Gob", func() {
		original := NewArrayListFrom(5, 3, 1, 4, 2)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewArrayList[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})
}

func TestArrayListSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(ArrayListSerializationTestSuite))
}

// ==========================
// 3. LinkedList Serialization Tests
// ==========================

type LinkedListSerializationTestSuite struct {
	suite.Suite
}

func (suite *LinkedListSerializationTestSuite) TestEmptyList() {
	suite.Run("JSON", func() {
		original := NewLinkedList[int]()
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewLinkedList[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})

	suite.Run("Gob", func() {
		original := NewLinkedList[int]()
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewLinkedList[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})
}

func (suite *LinkedListSerializationTestSuite) TestMultipleElements() {
	suite.Run("JSON", func() {
		original := NewLinkedListFrom(10, 20, 30, 40, 50)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewLinkedList[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})

	suite.Run("Gob", func() {
		original := NewLinkedListFrom(10, 20, 30, 40, 50)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewLinkedList[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})
}

func TestLinkedListSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(LinkedListSerializationTestSuite))
}

// ==========================
// 4. SegmentedList Serialization Tests
// ==========================

type SegmentedListSerializationTestSuite struct {
	suite.Suite
}

func (suite *SegmentedListSerializationTestSuite) TestEmptyList() {
	suite.Run("JSON", func() {
		original := NewSegmentedList[int]()
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewSegmentedList[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})

	suite.Run("Gob", func() {
		original := NewSegmentedList[int]()
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewSegmentedList[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})
}

func (suite *SegmentedListSerializationTestSuite) TestMultipleElements() {
	suite.Run("JSON", func() {
		original := NewSegmentedListFrom(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewSegmentedList[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})

	suite.Run("Gob", func() {
		original := NewSegmentedListFrom(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewSegmentedList[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})
}

func TestSegmentedListSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(SegmentedListSerializationTestSuite))
}

// ==========================
// 5. CowList Serialization Tests
// ==========================

type CowListSerializationTestSuite struct {
	suite.Suite
}

func (suite *CowListSerializationTestSuite) TestEmptyList() {
	suite.Run("JSON", func() {
		original := NewCOWList[int]()
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewCOWList[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})

	suite.Run("Gob", func() {
		original := NewCOWList[int]()
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewCOWList[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})
}

func (suite *CowListSerializationTestSuite) TestMultipleElements() {
	suite.Run("JSON", func() {
		original := NewCOWListFrom(100, 200, 300)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewCOWList[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})

	suite.Run("Gob", func() {
		original := NewCOWListFrom(100, 200, 300)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewCOWList[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})
}

func TestCowListSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(CowListSerializationTestSuite))
}

// ==========================
// 6. LockFreeList Serialization Tests
// ==========================

type LockFreeListSerializationTestSuite struct {
	suite.Suite
}

func (suite *LockFreeListSerializationTestSuite) TestEmptyList() {
	suite.Run("JSON", func() {
		original := NewLockFreeListOrdered[int]()
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewLockFreeListOrdered[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})

	suite.Run("Gob", func() {
		original := NewLockFreeListOrdered[int]()
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewLockFreeListOrdered[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})
}

func (suite *LockFreeListSerializationTestSuite) TestMultipleElements() {
	suite.Run("JSON", func() {
		original := NewLockFreeListFrom(func(a, b int) bool { return a == b }, 7, 8, 9)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewLockFreeList(func(a, b int) bool { return a == b })
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})

	suite.Run("Gob", func() {
		original := NewLockFreeListFrom(func(a, b int) bool { return a == b }, 7, 8, 9)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewLockFreeList(func(a, b int) bool { return a == b })
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})
}

func TestLockFreeListSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(LockFreeListSerializationTestSuite))
}

// ==========================
// 7. HashMap Serialization Tests
// ==========================

type HashMapSerializationTestSuite struct {
	suite.Suite
}

func (suite *HashMapSerializationTestSuite) TestEmptyMap() {
	suite.Run("JSON", func() {
		original := NewHashMap[string, int]()
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewHashMap[string, int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})

	suite.Run("Gob", func() {
		original := NewHashMap[string, int]()
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewHashMap[string, int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})
}

func (suite *HashMapSerializationTestSuite) TestSingleEntry() {
	suite.Run("JSON", func() {
		original := NewHashMap[string, int]()
		original.Put("answer", 42)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewHashMap[string, int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 1, restored.Size(), "Size should be one")
		val, ok := restored.Get("answer")
		require.True(suite.T(), ok, "Key should exist")
		assert.Equal(suite.T(), 42, val, "Value should match")
	})

	suite.Run("Gob", func() {
		original := NewHashMap[string, int]()
		original.Put("answer", 42)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewHashMap[string, int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 1, restored.Size(), "Size should be one")
		val, ok := restored.Get("answer")
		require.True(suite.T(), ok, "Key should exist")
		assert.Equal(suite.T(), 42, val, "Value should match")
	})
}

func (suite *HashMapSerializationTestSuite) TestMultipleEntries() {
	suite.Run("JSON", func() {
		original := NewHashMap[string, int]()
		original.Put("one", 1)
		original.Put("two", 2)
		original.Put("three", 3)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewHashMap[string, int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		for _, key := range []string{"one", "two", "three"} {
			origVal, _ := original.Get(key)
			restVal, ok := restored.Get(key)
			require.True(suite.T(), ok, "Key %s should exist", key)
			assert.Equal(suite.T(), origVal, restVal, "Value for key %s should match", key)
		}
	})

	suite.Run("Gob", func() {
		original := NewHashMap[string, int]()
		original.Put("one", 1)
		original.Put("two", 2)
		original.Put("three", 3)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewHashMap[string, int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		for _, key := range []string{"one", "two", "three"} {
			origVal, _ := original.Get(key)
			restVal, ok := restored.Get(key)
			require.True(suite.T(), ok, "Key %s should exist", key)
			assert.Equal(suite.T(), origVal, restVal, "Value for key %s should match", key)
		}
	})
}

func TestHashMapSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(HashMapSerializationTestSuite))
}

// ==========================
// 8. ArrayStack Serialization Tests
// ==========================

type ArrayStackSerializationTestSuite struct {
	suite.Suite
}

func (suite *ArrayStackSerializationTestSuite) TestEmptyStack() {
	suite.Run("JSON", func() {
		original := NewArrayStack[int]()
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewArrayStack[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})

	suite.Run("Gob", func() {
		original := NewArrayStack[int]()
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewArrayStack[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})
}

func (suite *ArrayStackSerializationTestSuite) TestMultipleElements() {
	suite.Run("JSON", func() {
		original := NewArrayStack[int]()
		original.PushAll(1, 2, 3, 4, 5)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewArrayStack[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})

	suite.Run("Gob", func() {
		original := NewArrayStack[int]()
		original.PushAll(1, 2, 3, 4, 5)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewArrayStack[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})
}

func TestArrayStackSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(ArrayStackSerializationTestSuite))
}

// ==========================
// 9. ArrayQueue Serialization Tests
// ==========================

type ArrayQueueSerializationTestSuite struct {
	suite.Suite
}

func (suite *ArrayQueueSerializationTestSuite) TestEmptyQueue() {
	suite.Run("JSON", func() {
		original := NewArrayQueue[int]()
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewArrayQueue[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})

	suite.Run("Gob", func() {
		original := NewArrayQueue[int]()
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewArrayQueue[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})
}

func (suite *ArrayQueueSerializationTestSuite) TestMultipleElements() {
	suite.Run("JSON", func() {
		original := NewArrayQueueFrom(10, 20, 30, 40, 50)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewArrayQueue[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})

	suite.Run("Gob", func() {
		original := NewArrayQueueFrom(10, 20, 30, 40, 50)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewArrayQueue[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})
}

func TestArrayQueueSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(ArrayQueueSerializationTestSuite))
}

// ==========================
// 10. ArrayDeque Serialization Tests
// ==========================

type ArrayDequeSerializationTestSuite struct {
	suite.Suite
}

func (suite *ArrayDequeSerializationTestSuite) TestEmptyDeque() {
	suite.Run("JSON", func() {
		original := NewArrayDeque[int]()
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewArrayDeque[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})

	suite.Run("Gob", func() {
		original := NewArrayDeque[int]()
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewArrayDeque[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})
}

func (suite *ArrayDequeSerializationTestSuite) TestMultipleElements() {
	suite.Run("JSON", func() {
		original := NewArrayDequeFrom(5, 10, 15, 20)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewArrayDeque[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})

	suite.Run("Gob", func() {
		original := NewArrayDequeFrom(5, 10, 15, 20)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewArrayDeque[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Order should be preserved")
	})
}

func TestArrayDequeSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(ArrayDequeSerializationTestSuite))
}

// ==========================
// 11. ConcurrentHashSet Serialization Tests
// ==========================

type ConcurrentHashSetSerializationTestSuite struct {
	suite.Suite
}

func (suite *ConcurrentHashSetSerializationTestSuite) TestEmptySet() {
	suite.Run("JSON", func() {
		original := NewConcurrentHashSet[int]()
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewConcurrentHashSet[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})

	suite.Run("Gob", func() {
		original := NewConcurrentHashSet[int]()
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewConcurrentHashSet[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})
}

func (suite *ConcurrentHashSetSerializationTestSuite) TestMultipleElements() {
	suite.Run("JSON", func() {
		original := NewConcurrentHashSetFrom(100, 200, 300, 400)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewConcurrentHashSet[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.True(suite.T(), restored.ContainsAll(100, 200, 300, 400), "Should contain all elements")
	})

	suite.Run("Gob", func() {
		original := NewConcurrentHashSetFrom(100, 200, 300, 400)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewConcurrentHashSet[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.True(suite.T(), restored.ContainsAll(100, 200, 300, 400), "Should contain all elements")
	})
}

func TestConcurrentHashSetSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(ConcurrentHashSetSerializationTestSuite))
}

// ==========================
// 12. ConcurrentHashMap Serialization Tests
// ==========================

type ConcurrentHashMapSerializationTestSuite struct {
	suite.Suite
}

func (suite *ConcurrentHashMapSerializationTestSuite) TestEmptyMap() {
	suite.Run("JSON", func() {
		original := NewConcurrentHashMap[string, int]()
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewConcurrentHashMap[string, int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})

	suite.Run("Gob", func() {
		original := NewConcurrentHashMap[string, int]()
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewConcurrentHashMap[string, int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})
}

func (suite *ConcurrentHashMapSerializationTestSuite) TestMultipleEntries() {
	suite.Run("JSON", func() {
		original := NewConcurrentHashMap[string, int]()
		original.Put("alpha", 1)
		original.Put("beta", 2)
		original.Put("gamma", 3)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewConcurrentHashMap[string, int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		for _, key := range []string{"alpha", "beta", "gamma"} {
			origVal, _ := original.Get(key)
			restVal, ok := restored.Get(key)
			require.True(suite.T(), ok, "Key %s should exist", key)
			assert.Equal(suite.T(), origVal, restVal, "Value for key %s should match", key)
		}
	})

	suite.Run("Gob", func() {
		original := NewConcurrentHashMap[string, int]()
		original.Put("alpha", 1)
		original.Put("beta", 2)
		original.Put("gamma", 3)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewConcurrentHashMap[string, int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		for _, key := range []string{"alpha", "beta", "gamma"} {
			origVal, _ := original.Get(key)
			restVal, ok := restored.Get(key)
			require.True(suite.T(), ok, "Key %s should exist", key)
			assert.Equal(suite.T(), origVal, restVal, "Value for key %s should match", key)
		}
	})
}

func TestConcurrentHashMapSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(ConcurrentHashMapSerializationTestSuite))
}

// ==========================
// 13. ConcurrentSkipSet Serialization Tests
// ==========================

type ConcurrentSkipSetSerializationTestSuite struct {
	suite.Suite
}

func (suite *ConcurrentSkipSetSerializationTestSuite) TestEmptySet() {
	suite.Run("JSON", func() {
		original := NewConcurrentSkipSet[int]()
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewConcurrentSkipSet[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})

	suite.Run("Gob", func() {
		original := NewConcurrentSkipSet[int]()
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewConcurrentSkipSet[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})
}

func (suite *ConcurrentSkipSetSerializationTestSuite) TestMultipleElements() {
	suite.Run("JSON", func() {
		original := NewConcurrentSkipSetFrom(5, 2, 8, 1, 9)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewConcurrentSkipSet[int]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.True(suite.T(), restored.ContainsAll(1, 2, 5, 8, 9), "Should contain all elements")
	})

	suite.Run("Gob", func() {
		original := NewConcurrentSkipSetFrom(5, 2, 8, 1, 9)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewConcurrentSkipSet[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.True(suite.T(), restored.ContainsAll(1, 2, 5, 8, 9), "Should contain all elements")
	})
}

func TestConcurrentSkipSetSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(ConcurrentSkipSetSerializationTestSuite))
}

// ==========================
// 14. ConcurrentSkipMap Serialization Tests
// ==========================

type ConcurrentSkipMapSerializationTestSuite struct {
	suite.Suite
}

func (suite *ConcurrentSkipMapSerializationTestSuite) TestEmptyMap() {
	suite.Run("JSON", func() {
		original := NewConcurrentSkipMap[int, string]()
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewConcurrentSkipMap[int, string]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})

	suite.Run("Gob", func() {
		original := NewConcurrentSkipMap[int, string]()
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewConcurrentSkipMap[int, string]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), 0, restored.Size(), "Size should be zero")
	})
}

func (suite *ConcurrentSkipMapSerializationTestSuite) TestMultipleEntries() {
	suite.Run("JSON", func() {
		original := NewConcurrentSkipMap[int, string]()
		original.Put(1, "one")
		original.Put(2, "two")
		original.Put(3, "three")
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewConcurrentSkipMap[int, string]()
		err = json.Unmarshal(data, restored)
		require.NoError(suite.T(), err, "Unmarshal should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		for _, key := range []int{1, 2, 3} {
			origVal, _ := original.Get(key)
			restVal, ok := restored.Get(key)
			require.True(suite.T(), ok, "Key %d should exist", key)
			assert.Equal(suite.T(), origVal, restVal, "Value for key %d should match", key)
		}
	})

	suite.Run("Gob", func() {
		original := NewConcurrentSkipMap[int, string]()
		original.Put(1, "one")
		original.Put(2, "two")
		original.Put(3, "three")
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewConcurrentSkipMap[int, string]()
		err = gob.NewDecoder(&buf).Decode(restored)
		require.NoError(suite.T(), err, "Gob decode should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		for _, key := range []int{1, 2, 3} {
			origVal, _ := original.Get(key)
			restVal, ok := restored.Get(key)
			require.True(suite.T(), ok, "Key %d should exist", key)
			assert.Equal(suite.T(), origVal, restVal, "Value for key %d should match", key)
		}
	})
}

func TestConcurrentSkipMapSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(ConcurrentSkipMapSerializationTestSuite))
}

// ==========================
// 15. TreeSet Serialization Tests (requires comparator)
// ==========================

type TreeSetSerializationTestSuite struct {
	suite.Suite
}

func (suite *TreeSetSerializationTestSuite) TestDirectUnmarshalReturnsError() {
	suite.Run("JSON", func() {
		original := NewTreeSetFrom(CompareFunc[int](), 5, 2, 8)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewTreeSet(CompareFunc[int]())
		err = json.Unmarshal(data, restored)
		assert.Error(suite.T(), err, "Direct unmarshal should return error")
		assert.Contains(suite.T(), err.Error(), "cannot unmarshal TreeSet directly", "Error message should indicate use of helper functions")
	})

	suite.Run("Gob", func() {
		original := NewTreeSetFrom(CompareFunc[int](), 5, 2, 8)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewTreeSet(CompareFunc[int]())
		err = gob.NewDecoder(&buf).Decode(restored)
		assert.Error(suite.T(), err, "Direct gob decode should return error")
		assert.Contains(suite.T(), err.Error(), "cannot unmarshal TreeSet directly", "Error message should indicate use of helper functions")
	})
}

func (suite *TreeSetSerializationTestSuite) TestOrderedTypeWithHelper() {
	suite.Run("JSON", func() {
		original := NewTreeSetFrom(CompareFunc[int](), 5, 2, 8, 1, 9)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored, err := UnmarshalTreeSetOrderedJSON[int](data)
		require.NoError(suite.T(), err, "Unmarshal with helper should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Sorted order should match")
	})

	suite.Run("Gob", func() {
		original := NewTreeSetFrom(CompareFunc[int](), 5, 2, 8, 1, 9)

		// Encode the underlying data slice directly for helper functions
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(original.ToSlice())
		require.NoError(suite.T(), err, "Gob encode data should succeed")

		restored, err := UnmarshalTreeSetOrderedGob[int](buf.Bytes())
		require.NoError(suite.T(), err, "Unmarshal with helper should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Sorted order should match")
	})
}

func (suite *TreeSetSerializationTestSuite) TestCustomComparatorWithHelper() {
	suite.Run("JSON", func() {
		reverseCompare := func(a, b int) int {
			return CompareFunc[int]()(b, a)
		}
		original := NewTreeSetFrom(reverseCompare, 5, 2, 8, 1, 9)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored, err := UnmarshalTreeSetJSON(data, reverseCompare)
		require.NoError(suite.T(), err, "Unmarshal with custom comparator should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Reverse sorted order should match")
	})

	suite.Run("Gob", func() {
		reverseCompare := func(a, b int) int {
			return CompareFunc[int]()(b, a)
		}
		original := NewTreeSetFrom(reverseCompare, 5, 2, 8, 1, 9)
		// Encode the underlying data slice directly for helper functions
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(original.ToSlice())
		require.NoError(suite.T(), err, "Gob encode data should succeed")

		restored, err := UnmarshalTreeSetGob(buf.Bytes(), reverseCompare)
		require.NoError(suite.T(), err, "Unmarshal with custom comparator should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Reverse sorted order should match")
	})
}

func TestTreeSetSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(TreeSetSerializationTestSuite))
}

// ==========================
// 16. TreeMap Serialization Tests (requires comparator)
// ==========================

type TreeMapSerializationTestSuite struct {
	suite.Suite
}

func (suite *TreeMapSerializationTestSuite) TestDirectUnmarshalReturnsError() {
	suite.Run("JSON", func() {
		original := NewTreeMap[int, string](CompareFunc[int]())
		original.Put(1, "one")
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewTreeMap[int, string](CompareFunc[int]())
		err = json.Unmarshal(data, restored)
		assert.Error(suite.T(), err, "Direct unmarshal should return error")
		assert.Contains(suite.T(), err.Error(), "cannot unmarshal TreeMap directly", "Error message should indicate use of helper functions")
	})

	suite.Run("Gob", func() {
		original := NewTreeMap[int, string](CompareFunc[int]())
		original.Put(1, "one")
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewTreeMap[int, string](CompareFunc[int]())
		err = gob.NewDecoder(&buf).Decode(restored)
		assert.Error(suite.T(), err, "Direct gob decode should return error")
		assert.Contains(suite.T(), err.Error(), "cannot unmarshal TreeMap directly", "Error message should indicate use of helper functions")
	})
}

func (suite *TreeMapSerializationTestSuite) TestOrderedKeyTypeWithHelper() {
	suite.Run("JSON", func() {
		original := NewTreeMapOrdered[int, string]()
		original.Put(3, "three")
		original.Put(1, "one")
		original.Put(2, "two")
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored, err := UnmarshalTreeMapOrderedJSON[int, string](data)
		require.NoError(suite.T(), err, "Unmarshal with helper should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		for _, key := range []int{1, 2, 3} {
			origVal, _ := original.Get(key)
			restVal, ok := restored.Get(key)
			require.True(suite.T(), ok, "Key %d should exist", key)
			assert.Equal(suite.T(), origVal, restVal, "Value for key %d should match", key)
		}
	})

	suite.Run("Gob", func() {
		original := NewTreeMapOrdered[int, string]()
		original.Put(3, "three")
		original.Put(1, "one")
		original.Put(2, "two")
		// Encode the underlying entries directly for helper functions
		entries := make([]serializableEntry[int, string], 0)
		for k, v := range original.Seq() {
			entries = append(entries, serializableEntry[int, string]{Key: k, Value: v})
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(entries)
		require.NoError(suite.T(), err, "Gob encode entries should succeed")

		restored, err := UnmarshalTreeMapOrderedGob[int, string](buf.Bytes())
		require.NoError(suite.T(), err, "Unmarshal with helper should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		for _, key := range []int{1, 2, 3} {
			origVal, _ := original.Get(key)
			restVal, ok := restored.Get(key)
			require.True(suite.T(), ok, "Key %d should exist", key)
			assert.Equal(suite.T(), origVal, restVal, "Value for key %d should match", key)
		}
	})
}

func (suite *TreeMapSerializationTestSuite) TestCustomComparatorWithHelper() {
	suite.Run("JSON", func() {
		reverseCompare := func(a, b int) int {
			return CompareFunc[int]()(b, a)
		}
		original := NewTreeMap[int, string](reverseCompare)
		original.Put(3, "three")
		original.Put(1, "one")
		original.Put(2, "two")
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored, err := UnmarshalTreeMapJSON[int, string](data, reverseCompare)
		require.NoError(suite.T(), err, "Unmarshal with custom comparator should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.Keys(), restored.Keys(), "Reverse sorted key order should match")
	})

	suite.Run("Gob", func() {
		reverseCompare := func(a, b int) int {
			return CompareFunc[int]()(b, a)
		}
		original := NewTreeMap[int, string](reverseCompare)
		original.Put(3, "three")
		original.Put(1, "one")
		original.Put(2, "two")

		// Encode the underlying entries directly for helper functions
		entries := make([]serializableEntry[int, string], 0)
		for k, v := range original.Seq() {
			entries = append(entries, serializableEntry[int, string]{Key: k, Value: v})
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(entries)
		require.NoError(suite.T(), err, "Gob encode entries should succeed")

		restored, err := UnmarshalTreeMapGob[int, string](buf.Bytes(), reverseCompare)
		require.NoError(suite.T(), err, "Unmarshal with custom comparator should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.Keys(), restored.Keys(), "Reverse sorted key order should match")
	})
}

func TestTreeMapSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(TreeMapSerializationTestSuite))
}

// ==========================
// 17. PriorityQueue Serialization Tests (requires comparator)
// ==========================

type PriorityQueueSerializationTestSuite struct {
	suite.Suite
}

func (suite *PriorityQueueSerializationTestSuite) TestDirectUnmarshalReturnsError() {
	suite.Run("JSON", func() {
		original := NewPriorityQueueOrdered[int]()
		original.Push(5)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewPriorityQueueOrdered[int]()
		err = json.Unmarshal(data, restored)
		assert.Error(suite.T(), err, "Direct unmarshal should return error")
		assert.Contains(suite.T(), err.Error(), "cannot unmarshal PriorityQueue directly", "Error message should indicate use of helper functions")
	})

	suite.Run("Gob", func() {
		original := NewPriorityQueueOrdered[int]()
		original.Push(5)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewPriorityQueueOrdered[int]()
		err = gob.NewDecoder(&buf).Decode(restored)
		assert.Error(suite.T(), err, "Direct gob decode should return error")
		assert.Contains(suite.T(), err.Error(), "cannot unmarshal PriorityQueue directly", "Error message should indicate use of helper functions")
	})
}

func (suite *PriorityQueueSerializationTestSuite) TestOrderedTypeWithHelper() {
	suite.Run("JSON", func() {
		original := NewPriorityQueueFrom(CompareFunc[int](), 5, 2, 8, 1, 9)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored, err := UnmarshalPriorityQueueOrderedJSON[int](data)
		require.NoError(suite.T(), err, "Unmarshal with helper should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		// Verify all elements are present by popping them
		origSlice := original.ToSortedSlice()
		restSlice := restored.ToSortedSlice()
		assert.Equal(suite.T(), origSlice, restSlice, "Elements should match")
	})

	suite.Run("Gob", func() {
		original := NewPriorityQueueFrom(CompareFunc[int](), 5, 2, 8, 1, 9)

		// Encode the underlying data slice directly for helper functions
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(original.ToSlice())
		require.NoError(suite.T(), err, "Gob encode data should succeed")

		restored, err := UnmarshalPriorityQueueOrderedGob[int](buf.Bytes())
		require.NoError(suite.T(), err, "Unmarshal with helper should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		origSlice := original.ToSortedSlice()
		restSlice := restored.ToSortedSlice()
		assert.Equal(suite.T(), origSlice, restSlice, "Elements should match")
	})
}

func (suite *PriorityQueueSerializationTestSuite) TestCustomComparatorWithHelper() {
	suite.Run("JSON", func() {
		maxHeapCompare := func(a, b int) int {
			return CompareFunc[int]()(b, a)
		}
		original := NewPriorityQueueFrom(maxHeapCompare, 5, 2, 8, 1, 9)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored, err := UnmarshalPriorityQueueJSON(data, maxHeapCompare)
		require.NoError(suite.T(), err, "Unmarshal with custom comparator should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		origSlice := original.ToSortedSlice()
		restSlice := restored.ToSortedSlice()
		assert.Equal(suite.T(), origSlice, restSlice, "Max-heap order should match")
	})

	suite.Run("Gob", func() {
		maxHeapCompare := func(a, b int) int {
			return CompareFunc[int]()(b, a)
		}
		original := NewPriorityQueueFrom(maxHeapCompare, 5, 2, 8, 1, 9)

		// Encode the underlying data slice directly for helper functions
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(original.ToSlice())
		require.NoError(suite.T(), err, "Gob encode data should succeed")

		restored, err := UnmarshalPriorityQueueGob(buf.Bytes(), maxHeapCompare)
		require.NoError(suite.T(), err, "Unmarshal with custom comparator should succeed")
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		origSlice := original.ToSortedSlice()
		restSlice := restored.ToSortedSlice()
		assert.Equal(suite.T(), origSlice, restSlice, "Max-heap order should match")
	})
}

func TestPriorityQueueSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(PriorityQueueSerializationTestSuite))
}

// ==========================
// 18. ConcurrentTreeSet Serialization Tests (requires comparator)
// ==========================

type ConcurrentTreeSetSerializationTestSuite struct {
	suite.Suite
}

func (suite *ConcurrentTreeSetSerializationTestSuite) TestDirectUnmarshalReturnsError() {
	suite.Run("JSON", func() {
		original := NewConcurrentTreeSetFrom(CompareFunc[int](), 5, 2, 8)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewConcurrentTreeSet(CompareFunc[int]())
		err = json.Unmarshal(data, restored)
		assert.Error(suite.T(), err, "Direct unmarshal should return error")
		assert.Contains(suite.T(), err.Error(), "cannot unmarshal ConcurrentTreeSet directly", "Error message should indicate use of helper functions")
	})

	suite.Run("Gob", func() {
		original := NewConcurrentTreeSetFrom(CompareFunc[int](), 5, 2, 8)
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewConcurrentTreeSet(CompareFunc[int]())
		err = gob.NewDecoder(&buf).Decode(restored)
		assert.Error(suite.T(), err, "Direct gob decode should return error")
		assert.Contains(suite.T(), err.Error(), "cannot unmarshal ConcurrentTreeSet directly", "Error message should indicate use of helper functions")
	})
}

func (suite *ConcurrentTreeSetSerializationTestSuite) TestOrderedTypeWithHelper() {
	suite.Run("JSON", func() {
		original := NewConcurrentTreeSetFrom(CompareFunc[int](), 5, 2, 8, 1, 9)
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		// Use TreeSet helper then wrap with ConcurrentTreeSet
		treeSet, err := UnmarshalTreeSetOrderedJSON[int](data)
		require.NoError(suite.T(), err, "Unmarshal TreeSet with helper should succeed")

		restored := NewConcurrentTreeSetFrom(CompareFunc[int](), treeSet.ToSlice()...)
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Sorted order should match")
	})

	suite.Run("Gob", func() {
		original := NewConcurrentTreeSetFrom(CompareFunc[int](), 5, 2, 8, 1, 9)

		// Encode the underlying data slice directly for helper functions
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(original.ToSlice())
		require.NoError(suite.T(), err, "Gob encode data should succeed")

		treeSet, err := UnmarshalTreeSetOrderedGob[int](buf.Bytes())
		require.NoError(suite.T(), err, "Unmarshal TreeSet with helper should succeed")

		restored := NewConcurrentTreeSetFrom(CompareFunc[int](), treeSet.ToSlice()...)
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")
		assert.Equal(suite.T(), original.ToSlice(), restored.ToSlice(), "Sorted order should match")
	})
}

func TestConcurrentTreeSetSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(ConcurrentTreeSetSerializationTestSuite))
}

// ==========================
// 19. ConcurrentTreeMap Serialization Tests (requires comparator)
// ==========================

type ConcurrentTreeMapSerializationTestSuite struct {
	suite.Suite
}

func (suite *ConcurrentTreeMapSerializationTestSuite) TestDirectUnmarshalReturnsError() {
	suite.Run("JSON", func() {
		original := NewConcurrentTreeMapOrdered[int, string]()
		original.Put(1, "one")
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		restored := NewConcurrentTreeMapOrdered[int, string]()
		err = json.Unmarshal(data, restored)
		assert.Error(suite.T(), err, "Direct unmarshal should return error")
		assert.Contains(suite.T(), err.Error(), "cannot unmarshal ConcurrentTreeMap directly", "Error message should indicate use of helper functions")
	})

	suite.Run("Gob", func() {
		original := NewConcurrentTreeMapOrdered[int, string]()
		original.Put(1, "one")
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(original)
		require.NoError(suite.T(), err, "Gob encode should succeed")

		restored := NewConcurrentTreeMapOrdered[int, string]()
		err = gob.NewDecoder(&buf).Decode(restored)
		assert.Error(suite.T(), err, "Direct gob decode should return error")
		assert.Contains(suite.T(), err.Error(), "cannot unmarshal ConcurrentTreeMap directly", "Error message should indicate use of helper functions")
	})
}

func (suite *ConcurrentTreeMapSerializationTestSuite) TestOrderedKeyTypeWithHelper() {
	suite.Run("JSON", func() {
		original := NewConcurrentTreeMapOrdered[int, string]()
		original.Put(3, "three")
		original.Put(1, "one")
		original.Put(2, "two")
		data, err := json.Marshal(original)
		require.NoError(suite.T(), err, "Marshal should succeed")

		// Use TreeMap helper then wrap with ConcurrentTreeMap
		treeMap, err := UnmarshalTreeMapOrderedJSON[int, string](data)
		require.NoError(suite.T(), err, "Unmarshal TreeMap with helper should succeed")

		restored := NewConcurrentTreeMapOrdered[int, string]()
		for k, v := range treeMap.Seq() {
			restored.Put(k, v)
		}
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		for _, key := range []int{1, 2, 3} {
			origVal, _ := original.Get(key)
			restVal, ok := restored.Get(key)
			require.True(suite.T(), ok, "Key %d should exist", key)
			assert.Equal(suite.T(), origVal, restVal, "Value for key %d should match", key)
		}
	})

	suite.Run("Gob", func() {
		original := NewConcurrentTreeMapOrdered[int, string]()
		original.Put(3, "three")
		original.Put(1, "one")
		original.Put(2, "two")

		// Encode the underlying entries directly for helper functions
		entries := make([]serializableEntry[int, string], 0)
		for k, v := range original.Seq() {
			entries = append(entries, serializableEntry[int, string]{Key: k, Value: v})
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(entries)
		require.NoError(suite.T(), err, "Gob encode entries should succeed")

		treeMap, err := UnmarshalTreeMapOrderedGob[int, string](buf.Bytes())
		require.NoError(suite.T(), err, "Unmarshal TreeMap with helper should succeed")

		restored := NewConcurrentTreeMapOrdered[int, string]()
		for k, v := range treeMap.Seq() {
			restored.Put(k, v)
		}
		assert.Equal(suite.T(), original.Size(), restored.Size(), "Size should match")

		for _, key := range []int{1, 2, 3} {
			origVal, _ := original.Get(key)
			restVal, ok := restored.Get(key)
			require.True(suite.T(), ok, "Key %d should exist", key)
			assert.Equal(suite.T(), origVal, restVal, "Value for key %d should match", key)
		}
	})
}

func TestConcurrentTreeMapSerializationTestSuite(t *testing.T) {
	suite.Run(t, new(ConcurrentTreeMapSerializationTestSuite))
}

// ==========================
// Error Handling Tests
// ==========================

type SerializationErrorHandlingTestSuite struct {
	suite.Suite
}

// Test invalid JSON data for direct serialization types
func (suite *SerializationErrorHandlingTestSuite) TestInvalidJSONData() {
	suite.Run("HashSet", func() {
		set := NewHashSet[int]()
		err := json.Unmarshal([]byte("invalid json"), set)
		assert.Error(suite.T(), err, "Should fail on invalid JSON")
	})

	suite.Run("ArrayList", func() {
		list := NewArrayList[int]()
		err := json.Unmarshal([]byte("{not an array}"), list)
		assert.Error(suite.T(), err, "Should fail on invalid JSON")
	})

	suite.Run("HashMap", func() {
		m := NewHashMap[string, int]()
		err := json.Unmarshal([]byte("not json"), m)
		assert.Error(suite.T(), err, "Should fail on invalid JSON")
	})
}

// Test invalid Gob data for direct serialization types
func (suite *SerializationErrorHandlingTestSuite) TestInvalidGobData() {
	suite.Run("HashSet", func() {
		set := NewHashSet[int]()
		err := gob.NewDecoder(bytes.NewReader([]byte("invalid gob data"))).Decode(set)
		assert.Error(suite.T(), err, "Should fail on invalid Gob data")
	})

	suite.Run("ArrayList", func() {
		list := NewArrayList[int]()
		err := gob.NewDecoder(bytes.NewReader([]byte("invalid gob data"))).Decode(list)
		assert.Error(suite.T(), err, "Should fail on invalid Gob data")
	})

	suite.Run("HashMap", func() {
		m := NewHashMap[string, int]()
		err := gob.NewDecoder(bytes.NewReader([]byte("invalid gob data"))).Decode(m)
		assert.Error(suite.T(), err, "Should fail on invalid Gob data")
	})

	suite.Run("ArrayStack", func() {
		stack := NewArrayStack[int]()
		err := gob.NewDecoder(bytes.NewReader([]byte("invalid gob data"))).Decode(stack)
		assert.Error(suite.T(), err, "Should fail on invalid Gob data")
	})

	suite.Run("ArrayQueue", func() {
		queue := NewArrayQueue[int]()
		err := gob.NewDecoder(bytes.NewReader([]byte("invalid gob data"))).Decode(queue)
		assert.Error(suite.T(), err, "Should fail on invalid Gob data")
	})

	suite.Run("ArrayDeque", func() {
		deque := NewArrayDeque[int]()
		err := gob.NewDecoder(bytes.NewReader([]byte("invalid gob data"))).Decode(deque)
		assert.Error(suite.T(), err, "Should fail on invalid Gob data")
	})
}

// Test nil comparator for helper functions
func (suite *SerializationErrorHandlingTestSuite) TestNilComparatorError() {
	suite.Run("UnmarshalTreeSetJSON", func() {
		data, _ := json.Marshal([]int{1, 2, 3})
		_, err := UnmarshalTreeSetJSON[int](data, nil)
		assert.Error(suite.T(), err, "Should fail with nil comparator")
		assert.Contains(suite.T(), err.Error(), "comparator required")
	})

	suite.Run("UnmarshalTreeMapJSON", func() {
		data, _ := json.Marshal(serializableMap[int, string]{
			Entries: []serializableEntry[int, string]{{Key: 1, Value: "one"}},
		})
		_, err := UnmarshalTreeMapJSON[int, string](data, nil)
		assert.Error(suite.T(), err, "Should fail with nil comparator")
		assert.Contains(suite.T(), err.Error(), "comparator required")
	})

	suite.Run("UnmarshalTreeSetGob", func() {
		var buf bytes.Buffer
		_ = gob.NewEncoder(&buf).Encode([]int{1, 2, 3})
		_, err := UnmarshalTreeSetGob[int](buf.Bytes(), nil)
		assert.Error(suite.T(), err, "Should fail with nil comparator")
		assert.Contains(suite.T(), err.Error(), "comparator required")
	})

	suite.Run("UnmarshalTreeMapGob", func() {
		entries := []serializableEntry[int, string]{{Key: 1, Value: "one"}}
		var buf bytes.Buffer
		_ = gob.NewEncoder(&buf).Encode(entries)
		_, err := UnmarshalTreeMapGob[int, string](buf.Bytes(), nil)
		assert.Error(suite.T(), err, "Should fail with nil comparator")
		assert.Contains(suite.T(), err.Error(), "comparator required")
	})

	suite.Run("UnmarshalPriorityQueueJSON", func() {
		data, _ := json.Marshal([]int{1, 2, 3})
		_, err := UnmarshalPriorityQueueJSON[int](data, nil)
		assert.Error(suite.T(), err, "Should fail with nil comparator")
		assert.Contains(suite.T(), err.Error(), "comparator required")
	})

	suite.Run("UnmarshalPriorityQueueGob", func() {
		var buf bytes.Buffer
		_ = gob.NewEncoder(&buf).Encode([]int{1, 2, 3})
		_, err := UnmarshalPriorityQueueGob[int](buf.Bytes(), nil)
		assert.Error(suite.T(), err, "Should fail with nil comparator")
		assert.Contains(suite.T(), err.Error(), "comparator required")
	})
}

// Test invalid JSON data for helper functions
func (suite *SerializationErrorHandlingTestSuite) TestHelperFunctionsWithInvalidJSON() {
	suite.Run("UnmarshalTreeSetJSON", func() {
		_, err := UnmarshalTreeSetJSON([]byte("invalid json"), CompareFunc[int]())
		assert.Error(suite.T(), err, "Should fail on invalid JSON")
	})

	suite.Run("UnmarshalTreeMapJSON", func() {
		_, err := UnmarshalTreeMapJSON[int, string]([]byte("invalid json"), CompareFunc[int]())
		assert.Error(suite.T(), err, "Should fail on invalid JSON")
	})

	suite.Run("UnmarshalPriorityQueueJSON", func() {
		_, err := UnmarshalPriorityQueueJSON([]byte("invalid json"), CompareFunc[int]())
		assert.Error(suite.T(), err, "Should fail on invalid JSON")
	})
}

// Test invalid Gob data for helper functions
func (suite *SerializationErrorHandlingTestSuite) TestHelperFunctionsWithInvalidGob() {
	suite.Run("UnmarshalTreeSetGob", func() {
		_, err := UnmarshalTreeSetGob([]byte("invalid gob"), CompareFunc[int]())
		assert.Error(suite.T(), err, "Should fail on invalid Gob data")
	})

	suite.Run("UnmarshalTreeMapGob", func() {
		_, err := UnmarshalTreeMapGob[int, string]([]byte("invalid gob"), CompareFunc[int]())
		assert.Error(suite.T(), err, "Should fail on invalid Gob data")
	})

	suite.Run("UnmarshalPriorityQueueGob", func() {
		_, err := UnmarshalPriorityQueueGob([]byte("invalid gob"), CompareFunc[int]())
		assert.Error(suite.T(), err, "Should fail on invalid Gob data")
	})
}

func TestSerializationErrorHandlingTestSuite(t *testing.T) {
	suite.Run(t, new(SerializationErrorHandlingTestSuite))
}

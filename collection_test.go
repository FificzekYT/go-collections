package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntry_Unpack(t *testing.T) {
	t.Parallel()
	e := Entry[int, string]{Key: 1, Value: "a"}
	k, v := e.Unpack()
	assert.Equal(t, 1, k, "Unpack should return original key")
	assert.Equal(t, "a", v, "Unpack should return original value")
}

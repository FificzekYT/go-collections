package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntry_Unpack(t *testing.T) {
	e := Entry[int, string]{Key: 1, Value: "a"}
	k, v := e.Unpack()
	assert.Equal(t, 1, k)
	assert.Equal(t, "a", v)
}

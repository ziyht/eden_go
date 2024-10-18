package esl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	set := New[int, int]()
	set.Add(1, 1)
	set.Add(2, 2)
	set.Add(3, 3)	
	assert.Equal(t, 3, set.Len())
	assert.True(t, set.Contains(1))
	assert.True(t, set.Contains(2))
	assert.True(t, set.Contains(3))
	assert.False(t, set.Contains(4))
	set.Remove(2)
	assert.Equal(t, 2, set.Len())
	assert.False(t, set.Contains(2))

}
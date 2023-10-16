package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	q := NewSyncQueue()
	q.Push(1)
	assert.Equal(t, 1, q.Length())
	q.Push(2, 2, 3, 4, 5, 6)
	assert.Equal(t, 1, q.Shift())
	assert.Equal(t, 6, q.Pop())
	dropped := q.RemoveWhere(func(r interface{}, stop *bool) bool {
		i := r.(int)
		return i&1 == 1
	})
	assert.Equal(t, []interface{}{3, 5}, dropped)
	assert.Equal(t, 3, q.Length())
}

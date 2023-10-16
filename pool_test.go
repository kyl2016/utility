package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	var pool BytesPool
	poolCapacity := 10

	bs := pool.Get(poolCapacity)
	assert.NotNil(t, bs)
	assert.Equal(t, cap(bs), poolCapacity)
	bs = append(bs, []byte("1234567890abcdefghijklmnopqrstuvwxyz")...)
	pool.Put(bs)

	bs2 := pool.Get(poolCapacity)
	assert.Equal(t, 0, len(bs2))
	assert.Greater(t, cap(bs2), poolCapacity)
}

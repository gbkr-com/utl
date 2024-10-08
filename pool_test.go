package utl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoolWithValue(t *testing.T) {

	pool := NewPool(
		2,
		func(item int) int {
			return 0
		},
	)
	assert.Equal(t, 2, len(pool.c))

	next := pool.Get()
	assert.Equal(t, 0, next)

	next = pool.Get()
	assert.Equal(t, 0, next)
	assert.Equal(t, 0, len(pool.c))

	next = pool.Get()
	assert.Equal(t, 0, next)
	assert.Equal(t, 0, len(pool.c))

	pool.Recycle(1)
	assert.Equal(t, 1, len(pool.c))

	next = pool.Get()
	assert.Equal(t, 0, next)

}

func TestPoolWithRef(t *testing.T) {

	type x struct{ i int }

	pool := NewPool(
		2,
		func(item *x) *x {
			if item == nil {
				return &x{}
			}
			item.i = 0
			return item
		},
	)
	assert.Equal(t, 2, len(pool.c))

	first := pool.Get()
	assert.NotNil(t, first)
	assert.Equal(t, 0, first.i)

	second := pool.Get()
	assert.NotNil(t, second)
	assert.Equal(t, 0, second.i)
	assert.Equal(t, 0, len(pool.c), "pool now empty")

	first.i = 1
	pool.Recycle(first)

	first = pool.Get()
	assert.NotNil(t, first)
	assert.Equal(t, 0, first.i, "item has been zeroed")

	pool.Get()
	assert.Equal(t, 0, len(pool.c), "pool empty again")

	next := pool.Get()
	assert.NotNil(t, next)
	assert.Equal(t, 0, next.i, "item was manufactured")

}

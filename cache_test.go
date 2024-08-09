package utl

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {

	ctx := context.Background()

	mini := miniredis.RunT(t)
	defer mini.Close()
	rdb := redis.NewClient(&redis.Options{
		Addr: mini.Addr(),
	})

	rdb.Set(ctx, "key", "one", 0)

	cache := NewCache(
		time.Second,
		func(key string) (string, bool) {
			v, err := rdb.Get(ctx, key).Result()
			if err != nil {
				return "", false
			}
			return v, true
		},
	)

	v, ok := cache.Get("key")
	assert.True(t, ok)
	assert.Equal(t, "one", v)

	rdb.Set(ctx, "key", "two", 0)

	v, ok = cache.Get("key")
	assert.True(t, ok)
	assert.Equal(t, "one", v)

	time.Sleep(time.Second)

	v, ok = cache.Get("key")
	assert.True(t, ok)
	assert.Equal(t, "two", v)

}

package utl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {

	t.Skip()

	r := NewRateLimiter(1, time.Second)
	assert.True(t, r.Try(), "get token immediately")

	time.Sleep(time.Second)

	assert.True(t, r.Try(), "get token after refill")
	assert.False(t, r.Try(), "no more tokens")

	time.Sleep(time.Second)

	assert.True(t, r.Try(), "get token after refill")

}

package utl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnce(t *testing.T) {

	do := NewOnce(true)

	assert.False(t, do.Reset())

	assert.True(t, do.Try())
	assert.False(t, do.Try())

	assert.True(t, do.Reset())
	assert.False(t, do.Reset())

	assert.True(t, do.Try())
	assert.False(t, do.Try())

}

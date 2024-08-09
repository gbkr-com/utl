package utl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConflatingUpdate(t *testing.T) {

	type update struct {
		a *int64
		b *string
	}

	bundle := NewConflatingUpdate(
		func(pending, update *update) *update {
			if pending == nil {
				return update
			}
			if update.a != nil {
				pending.a = update.a
			}
			if update.b != nil {
				pending.b = update.b
			}
			return pending
		},
	)

	var upd update

	upd.a = ref[int64](42)
	bundle.Push(&upd)

	upd.a = ref[int64](43)
	upd.b = ref("b")
	bundle.Push(&upd)

	u := bundle.Pop()
	assert.NotNil(t, u)
	assert.NotNil(t, u.a)
	assert.Equal(t, int64(43), *u.a)
	assert.NotNil(t, u.b)
	assert.Equal(t, "b", *u.b)

}

func ref[T any](x T) *T { return &x }

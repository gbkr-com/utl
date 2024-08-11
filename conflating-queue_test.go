package utl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConflatingQueue(t *testing.T) {

	type quote struct {
		symbol string
		price  int
	}
	key := func(q *quote) string { return q.symbol }

	queue := NewConflatingQueue(key)
	assert.Nil(t, queue.Pop(), "empty queue returns the zero value of T")

	queue.Push(&quote{symbol: "A", price: 1})
	queue.Push(&quote{symbol: "B", price: 2})
	queue.Push(&quote{symbol: "A", price: 3})

	pop := queue.Pop()
	assert.Equal(t, "A", pop.symbol, "arrival sequence preserved")
	assert.Equal(t, 3, pop.price, "second item overrides first")

	pop = queue.Pop()
	assert.Equal(t, "B", pop.symbol, "arrival sequence preserved")
	assert.Equal(t, 2, pop.price)

	pop = queue.Pop()
	assert.Nil(t, pop)

}

func TestConflatingQueueUsingChan(t *testing.T) {

	type quote struct {
		symbol string
		price  int
	}
	key := func(q *quote) string { return q.symbol }

	queue := NewConflatingQueue(key)

	select {
	case <-queue.C():
		t.Fail()
	default:
	}

	queue.Push(&quote{symbol: "A", price: 1})
	queue.Push(&quote{symbol: "B", price: 2})
	queue.Push(&quote{symbol: "A", price: 3})

	k := 0

LOOP:
	for {
		select {
		case <-queue.C():
			q := queue.Pop()
			assert.NotNil(t, q)
			switch k {
			case 0:
				assert.Equal(t, "A", q.symbol)
			case 1:
				assert.Equal(t, "B", q.symbol)
			}
			k++
		default:
			assert.Equal(t, 2, k)
			break LOOP
		}
	}

}

func TestConflatingQueuePopAndChanAreConsist(t *testing.T) {

	type quote struct {
		symbol string
		price  int
	}
	key := func(q *quote) string { return q.symbol }

	queue := NewConflatingQueue(key)

	queue.Push(&quote{symbol: "A", price: 1})
	assert.NotNil(t, queue.Pop())

	select {
	case <-queue.C():
		t.Fail()
	default:
	}

}

func TestConflatingQueueWithPoolOption(t *testing.T) {

	type quote struct {
		symbol string
		price  int
	}
	key := func(q *quote) string { return q.symbol }

	pool := NewPool(8,
		func(q *quote) *quote {
			if q == nil {
				return &quote{}
			}
			q.symbol, q.price = "", 0
			return q
		},
	)
	assert.Equal(t, 8, len(pool.c))

	queue := NewConflatingQueue(key, WithPoolOption[string](pool))

	q := pool.Get()
	assert.Equal(t, 7, len(pool.c))
	assert.NotNil(t, q)
	q.symbol = "A"
	q.price = 1

	queue.Push(q)

	q = pool.Get()
	assert.Equal(t, 6, len(pool.c))
	q.symbol = "A"
	q.price = 2

	queue.Push(q)
	assert.Equal(t, 7, len(pool.c), "conflated item recycled")

	q = queue.Pop()
	assert.NotNil(t, q)
	assert.Equal(t, 2, q.price)

	pool.Recycle(q)
	assert.Equal(t, 8, len(pool.c))

}

func TestConflatingQueueWithConflateOption(t *testing.T) {

	type instruction struct {
		id  string
		qty *float64
		px  *float64
	}
	key := func(q *instruction) string { return q.id }

	queue := NewConflatingQueue(key,
		WithConflateOption[string, *instruction](
			func(existing, item *instruction) *instruction {
				if item.qty != nil {
					existing.qty = item.qty
				}
				if item.px != nil {
					existing.px = item.px
				}
				return existing
			},
		),
	)

	queue.Push(&instruction{id: "A"})
	queue.Push(&instruction{id: "A", qty: ref(100.0)})
	queue.Push(&instruction{id: "A", px: ref(42.0)})

	popped := queue.Pop()
	assert.NotNil(t, popped)
	assert.Equal(t, "A", popped.id)
	assert.NotNil(t, popped.qty)
	assert.Equal(t, 100.0, *popped.qty)
	assert.NotNil(t, popped.px)
	assert.Equal(t, 42.0, *popped.px)

	assert.Nil(t, queue.Pop())

}

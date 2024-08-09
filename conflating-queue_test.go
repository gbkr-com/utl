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

	queue := NewConflatingQueue(func(item *quote) string { return item.symbol })
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

	queue := NewConflatingQueue(func(item *quote) string { return item.symbol })

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

	queue := NewConflatingQueue(func(item *quote) string { return item.symbol })

	queue.Push(&quote{symbol: "A", price: 1})
	assert.NotNil(t, queue.Pop())

	select {
	case <-queue.C():
		t.Fail()
	default:
	}

}

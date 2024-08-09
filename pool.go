package utl

// A Pool of reusable items.
type Pool[T any] struct {
	c chan T
	z func(T) T
}

// NewPool returns a [*Pool] of the given type with capacity 'n'. The function
// arguments is to 'zero' an item when it is pushed back to the pool.
// If T is a pointer type then the function should zero the referenced item, not
// assign nil. For example:
//
//	pool := NewPool(
//		2,
//		func(item *X) *X {
//		    if item == nil {
//		        return &X{}
//		    }
//		    item.x = 0
//		    return item
//		},
//	)
func NewPool[T any](n int, zero func(T) T) *Pool[T] {
	pool := &Pool[T]{
		c: make(chan T, n),
		z: zero,
	}
	items := make([]T, n)
	for i := range items {
		zeroed := zero(items[i])
		pool.c <- zeroed
	}
	return pool
}

// Pop returns a reusable item from the pool. If the pool is empty, a new zeroed
// item is created.
func (x *Pool[T]) Pop() T {
	select {
	case t := <-x.c:
		return t
	default:
		var t T
		return x.z(t)
	}
}

// Push returns an item to the pool. The item is zeroed as it is added.
func (x *Pool[T]) Push(t T) {
	select {
	case x.c <- x.z(t):
	default:
	}
}

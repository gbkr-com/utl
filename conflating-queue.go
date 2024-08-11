package utl

import "sync"

// ConflatingQueue is a queue of items V each having a key K. Multiple items may
// arrive for a key value but only the most recent item is of interest.
//
// For example, consider a queue of prices. Prices for 'A' then 'B' are pushed
// to the queue. Another price for 'A' is pushed and this replaces the existing
// price for 'A' in the queue. When the consumer calls [ConflatingQueue.Pop] it
// will return the latest price for A (because 'A' was the first of any price pushed),
// and a second call will return the price for 'B'.
type ConflatingQueue[K comparable, V any] struct {
	notifier
	queue    []V          // A slice of the items in arrival sequence but without duplicate keys.
	index    map[K]int    // The index of each key within the queue.
	key      func(V) K    // The function to return the key for a given item.
	lock     sync.Mutex   //
	pool     *Pool[V]     // Optional pool.
	conflate func(V, V) V // Optional conflate function.
}

// ConflatingQueueOption is an option specified when a queue is manufactured.
type ConflatingQueueOption[K comparable, V any] func(*ConflatingQueue[K, V])

// WithPoolOption specifies that conflated items need to be returned to a pool.
func WithPoolOption[K comparable, V any](pool *Pool[V]) ConflatingQueueOption[K, V] {
	return func(queue *ConflatingQueue[K, V]) {
		queue.pool = pool
	}
}

// WithConflateOption specifies that the given function will be used to conflate
// an item, rather than the default of substitution.
func WithConflateOption[K comparable, V any](conflate func(existing, next V) V) ConflatingQueueOption[K, V] {
	return func(queue *ConflatingQueue[K, V]) {
		queue.conflate = conflate
	}
}

// NewConflatingQueue returns a [*ConflatingQueue] ready to use. The function
// argument returns the key value for a given item.
func NewConflatingQueue[K comparable, V any](key func(V) K, options ...ConflatingQueueOption[K, V]) *ConflatingQueue[K, V] {
	queue := &ConflatingQueue[K, V]{
		queue: make([]V, 0, 128),
		index: make(map[K]int),
		key:   key,
	}
	queue.init()
	for _, option := range options {
		option(queue)
	}
	return queue
}

// C returns a channel on which a struct{} is sent when the queue is not empty.
// This is useful when including the queue in a select statement. Use
// [ConflatingQueue.Pop] to consume the item.
//
// This channel is only for notification; the length of the channel does not
// reflect the number of items in the queue.
func (x *ConflatingQueue[K, V]) C() chan struct{} {
	return x.c
}

// Push an item into the queue. If an item with the same key is in the queue:
//   - If no [WithConflateOption] was specified this item will replace the existing item
//   - If [WithConflateOption] was specified this item will be merged into the existing item using the provided function
func (x *ConflatingQueue[K, V]) Push(item V) {
	x.lock.Lock()
	defer x.lock.Unlock()

	key := x.key(item)

	index, present := x.index[key]
	if present {
		existing := x.queue[index]
		switch {
		case x.conflate != nil:
			//
			// Conflate the item then recycle 'item'.
			//
			item = x.conflate(existing, item)
			x.queue[index] = item

		case x.pool != nil:
			//
			// Replace and recycle the existing item.
			//
			x.queue[index] = item
			x.pool.Recycle(existing)

		default:
			//
			// Replace.
			//
			x.queue[index] = item
		}
		x.notify()
		return
	}

	index = len(x.queue)
	x.queue = append(x.queue, item)
	x.index[key] = index
	x.notify()
}

// Pop from the queue. If the queue is empty this will return the zero value
// of T.
func (x *ConflatingQueue[K, V]) Pop() V {
	x.lock.Lock()
	defer x.lock.Unlock()

	var result V

	switch len(x.queue) {
	case 0:
		x.clear()

	case 1:
		result = x.queue[0]
		x.queue = x.queue[:0]
		delete(x.index, x.key(result))
		x.clear()

	default:
		result = x.queue[0]
		x.queue = x.queue[1:]
		delete(x.index, x.key(result))
		for k := range x.index {
			x.index[k]--
		}
		x.notify()
	}

	return result
}

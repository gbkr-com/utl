package utl

import "sync"

// A ConflatingUpdate of contemporaneous data. T is a struct or struct pointer
// representing an update of one or more data points.
//
// For example, a consumer is to make decisions based upon different events and
// the consumer wants to assess all those events at once rather than to receive
// them in sequence.
type ConflatingUpdate[T any] struct {
	notifier
	pending T
	f       func(pending, update T) T
	lock    sync.Mutex
}

// NewConflatingUpdate returns a [*ConflatingUpdate] ready to use. The function
// argument is used to merge a new update with any pending update that has not
// yet been consumed.
func NewConflatingUpdate[T any](f func(pending, update T) T) *ConflatingUpdate[T] {
	bundle := &ConflatingUpdate[T]{f: f}
	bundle.init()
	return bundle
}

// C returns a channel on which a struct{} is sent when there is a pending update.
// This is useful when including the queue in a select statement. Use
// [ConflatingUpdate.Pop] to consume the update.
func (x *ConflatingUpdate[T]) C() chan struct{} {
	return x.c
}

// Push an update.
func (x *ConflatingUpdate[T]) Push(update T) {
	x.lock.Lock()
	defer x.lock.Unlock()

	x.pending = x.f(x.pending, update)
	x.notify()
}

// Pop the update.
func (x *ConflatingUpdate[T]) Pop() T {
	x.lock.Lock()
	defer x.lock.Unlock()

	result := x.pending
	var init T
	x.pending = init
	x.clear()

	return result
}

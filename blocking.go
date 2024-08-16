package utl

import (
	"context"
	"errors"
	"time"
)

// Blocking interface.
type Blocking interface {
	Done()
	Wait() error
}

// ErrTimeout is the error that can be returned by [Blocking.Wait].
var ErrTimeout = errors.New("utl: timeout")

type blocking struct {
	ctx  context.Context
	cxl  context.CancelFunc
	done chan struct{}
}

// Done signals a successful conclusion. Repeated calls to [Blocking.Done] have
// no affect.
func (x *blocking) Done() {
	select {
	case x.done <- struct{}{}:
	default:
	}
}

// Wait until [Blocking.Done] or the timeout.
func (x *blocking) Wait() error {
	select {
	case <-x.done:
		x.cxl()
		return nil
	case <-x.ctx.Done():
		return ErrTimeout
	}
}

// WithTimeout returns a [Blocking] interface ready to use..
func WithTimeout(d time.Duration) Blocking {

	ctx, cxl := context.WithTimeout(context.Background(), d)
	b := &blocking{
		ctx:  ctx,
		cxl:  cxl,
		done: make(chan struct{}, 1),
	}
	return b
}

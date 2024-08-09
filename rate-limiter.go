package utl

import "time"

// A RateLimiter for governing the rate of an action.
type RateLimiter struct {
	tokens chan struct{}
}

func (r *RateLimiter) fill() {
	for i := 0; i < cap(r.tokens); i++ {
		select {
		case r.tokens <- struct{}{}:
		default:
			return
		}
	}
}

// NewRateLimiter returns a [*RateLimiter] with the given request rate of 'n'
// requests over 'interval' time period.
func NewRateLimiter(n int, interval time.Duration) *RateLimiter {
	r := &RateLimiter{
		tokens: make(chan struct{}, n),
	}
	r.fill()
	go func() {
		for range time.Tick(interval) {
			r.fill()
		}
	}()
	return r
}

// Block until a token is available.
func (r *RateLimiter) Block() {
	<-r.tokens
}

// Try returns true if a token is available, otherwise false. Try does not block.
func (r *RateLimiter) Try() bool {
	select {
	case <-r.tokens:
		return true
	default:
		return false
	}
}

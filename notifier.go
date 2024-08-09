package utl

type notifier struct {
	c chan struct{}
}

func (x *notifier) init() {
	x.c = make(chan struct{}, 1)
}

// Try to send a struct{} on the notification channel. Return true if this was
// possible.
func (x *notifier) notify() bool {
	select {
	case x.c <- struct{}{}:
		return true
	default:
		return false
	}
}

// Try to receive from the notification channel. Return true if this was possible.
func (x *notifier) clear() bool {
	select {
	case <-x.c:
		return true
	default:
		return false
	}
}

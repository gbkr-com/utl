package utl

// Once is a repeatable [sync.Once].
type Once struct {
	notifier
}

// NewOnce returns a [*Once] ready to use. If the argument is true the returned
// [*Once] is set to return true the first time [Once.Try] is called.
func NewOnce(set bool) *Once {
	once := &Once{}
	once.init()
	if set {
		once.notify()
	}
	return once
}

// Reset Once. Return true if it was actually reset, false if ignored.
func (x *Once) Reset() bool {
	return x.notify()
}

// Try returns true if an action can be taken. If this returns true, subsequent
// calls to this method will return false until [Once.Reset] has been called.
func (x *Once) Try() bool {
	return x.clear()
}

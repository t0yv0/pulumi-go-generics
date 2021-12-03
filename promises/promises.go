package promises

const (
	PendingState = iota
	ResolvedState
	RejectedState
)

type Promise[T any] struct {
	// closed when `state != PendingState`
	complete <-chan struct{}

	// one of `{PendingState,ResolvedState,RejectedState}`
	state uint32

        // set iff `state == ResolvedState`
	value T

	// non-nil iff `state == RejectedState`
	err error

	// threads through promise chains to simplify the API
	observer Observer
}

func Resolved[T any](observer Observer, value T) *Promise[T] {
	observer.Created()
	observer.Resolved()
	c := make(chan struct{})
	close(c)
	return &Promise[T]{
		state:    ResolvedState,
		value:    value,
		observer: observer,
		complete: c,
	}
}

func Rejected[T any](observer Observer, err error) *Promise[T] {
	observer.Created()
	observer.Rejected(err)
	c := make(chan struct{})
	close(c)
	return &Promise[T]{
		state:    RejectedState,
		err:      err,
		observer: observer,
		complete: c,
	}
}

// Builds a new Promise and gives a resolve and reject callbacks to
// fullfill it. Only one of these callbacks should be called once (or
// the code panics).
func NewPromise[T any](observer Observer) (*Promise[T], func(T), func(error)) {
	complete := make(chan struct{})

	f := &Promise[T]{
		complete: complete,
		state:    PendingState,
		observer: observer,
	}
	observer.Created()

	resolve := func(value T) {
		f.value = value
		f.state = ResolvedState
		close(complete) // possible "panic: close of closed channel"
		observer.Resolved()
	}

	reject := func(err error) {
		f.err = err
		f.state = RejectedState
		close(complete) // possible "panic: close of closed channel"
		observer.Rejected(err)
	}

	return f, resolve, reject
}

func (p *Promise[T]) Await() (T, error) {
	<- p.complete
 	return p.value, p.err
}

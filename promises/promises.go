package promises

import (
	"fmt"
)

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
		fmt.Printf("reject error=%v\n", err)
		close(complete) // possible "panic: close of closed channel"
		observer.Rejected(err)
	}

	return f, resolve, reject
}

func (p *Promise[T]) Await() (T, error) {
	<- p.complete
 	return p.value, p.err
}

func MapErr[T any, U any](f func(T) (U, error)) func (*Promise[T]) *Promise[U] {
	return func(p *Promise[T]) *Promise[U] {
		var res *Promise[U]
		switch p.state {
		case ResolvedState:
			v, err := f(p.value)
			if err != nil {
				res = Rejected[U](p.observer, err)
			} else {
				res = Resolved(p.observer, v)
			}
		case RejectedState:
			return Rejected[U](p.observer, p.err)
		case PendingState:
			res, resolve, reject := NewPromise[U](p.observer)
			go func() {
				v, err := res.Await()
				if err != nil {
					reject(err)
				}
				resolve(v)
			}()
		}
		return res
	}
}

func Map[T any, U any](f func(T) U) func (*Promise[T]) *Promise[U] {
	return MapErr(func (x T) (U, error) {
		return f(x), nil
	})
}

func All[T any](observer Observer, promises []*Promise[T]) *Promise[[]T] {
	if len(promises) == 0 {
		return Resolved(observer, []T{})
	}

	allPending := true
	for _, p := range promises {
		if p.state != PendingState {
			allPending = false
		}
	}

	if allPending {
		values := []T{}
		for _, p := range promises {
			if p.state == RejectedState {
				return Rejected[[]T](observer, p.err)
			}
			values = append(values, p.value)
		}
		return Resolved(observer, values)
	}

	result, resolve, reject := NewPromise[[]T](observer)

	completed := make(chan int)

	observe := func(i int) {
		<- promises[i].complete
		completed <- i
	}

	for i := range promises {
		go observe(i)
	}

	go func() {
		done := 0
		results := make([]T, len(promises))
		for result.state == PendingState {
			i := <-completed
			if promises[i].state == RejectedState {
				reject(promises[i].err)
			} else {
				results[i] = promises[i].value
				done++
				if done == len(promises) {
					resolve(results)
				}
			}
		}
	}()

	return result
}

func Join[T any](pp *Promise[*Promise[T]]) *Promise[T] {
	switch pp.state {
	case RejectedState:
		return Rejected[T](pp.observer, pp.err)
	case ResolvedState:
		return pp.value
	default:
		result, resolve, reject := NewPromise[T](pp.observer)
		go func() {
			f, err := pp.Await()
			if err != nil {
				reject(err)
				return
			}
			v, err := f.Await()
			if err != nil {
				reject(err)
				return
			}
			resolve(v)
		}()
		return result
	}
}

package promises

import (
	"sync"
)

type Waiter struct {
	// tracks active futures counts, Wait for idle capability
	workGroup workGroup

	// internal hook to report errors from futures
	onError func(error)

	// closes if/when firstError is thrown
	firstErrorChan chan interface{}

	// wait for firstErrorChan before querying
	firstError error
}

var _ Observer = &Waiter{}

func (w *Waiter) AwaitIdle() {
	w.workGroup.Wait()
}

func (w *Waiter) AwaitIdleOrFirstError() error {
	idle := make(chan interface{})

	go func() {
		w.AwaitIdle()
		close(idle)
	}()

	select {
	case <-w.firstErrorChan:
		return w.firstError
	case <-idle:
		return nil
	}

	return nil
}

func (w *Waiter) Created() {
	w.workGroup.Add(1)
}

func (w *Waiter) Rejected(err error) {
	w.workGroup.Done()
	if err != nil {
		w.onError(err)
	}
}

func (w *Waiter) Resolved() {
	w.workGroup.Done()
}

func NewWaiter() *Waiter {
	var once sync.Once
	var w *Waiter

	onError := func(err error) {
		if err == nil {
			return
		}
		once.Do(func() {
			w.firstError = err
			close(w.firstErrorChan)
		})
	}

	w = &Waiter{
		onError:        onError,
		firstErrorChan: make(chan interface{}),
	}
	return w
}

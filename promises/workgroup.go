package promises

import (
	"fmt"
	"sync"
)

// Mimicks the interface of `sync.WaitGroup` but does not panic in
// case of races between `Wait` and `Add` with a positive delta in the
// state with a zero counter. The reason `sync.WaitGroup` panics is to
// warn about a race condition. Using `workGroup` implicitly accept
// these race conditions instead. Use sparingly and document why it is
// used.
type workGroup struct {
	mutex   sync.Mutex
	cond    *sync.Cond
	counter int
}

func (wg *workGroup) Wait() {
	wg.mutex.Lock()
	defer wg.mutex.Unlock()

	if wg.cond == nil {
		wg.cond = sync.NewCond(&wg.mutex)
	}

	for wg.counter > 0 {
		wg.cond.Wait()
	}
}

func (wg *workGroup) Add(delta int) {
	wg.mutex.Lock()
	defer wg.mutex.Unlock()

	if wg.cond == nil {
		wg.cond = sync.NewCond(&wg.mutex)
	}

	c := wg.counter + delta

	if c < 0 {
		panic(fmt.Sprintf("Adding %d would make workGroup counter negative: %d + %d = %d",
			delta, wg.counter, delta, c))
	}

	wg.counter = c

	if c == 0 {
		wg.cond.Broadcast()
	}
}

func (wg *workGroup) Done() {
	wg.Add(-1)
}

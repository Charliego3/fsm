package fsm

import (
	"context"
	"fmt"
	"github.com/RussellLuo/timingwheel"
	"sync"
	"time"
)

var (
	timing   *timingwheel.TimingWheel
	timingMu sync.Mutex
)

func (f *FSM[T]) getTiming() *timingwheel.TimingWheel {
	timingMu.Lock()
	defer timingMu.Unlock()

	if timing == nil {
		timing = timingwheel.NewTimingWheel(f.tick, f.wheelSize)
		timing.Start()
	}
	return timing
}

type RetryErr[T comparable] struct {
	State T
}

func (r RetryErr[T]) Error() string {
	return fmt.Sprintf("fsm: will be retry with state: %v", r.State)
}

type CanceledErr[T comparable] struct {
	State T
}

func (r CanceledErr[T]) Error() string {
	return fmt.Sprintf("fsm: canceled with state: %v", r.State)
}

type Event[T comparable] struct {
	ctx  context.Context
	FSM  *FSM[T]
	Args []any
	Err  error
}

func (e *Event[T]) Retry(delay time.Duration) {
	if !e.FSM.async {
		time.Sleep(delay)
		e.FSM.transition(e)
		return
	}

	e.Err = RetryErr[T]{e.FSM.current}
	e.FSM.getTiming().AfterFunc(delay, func() {
		e.FSM.transition(e)
	})
}

func (e *Event[T]) Cancel() {
	e.Err = CanceledErr[T]{e.FSM.current}
}

package fsm

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Option[T comparable] func(*FSM[T])

// FSM is a Finite State Machine
type FSM[T comparable] struct {
	// current is the state of currently
	current T

	// transitions specify src state to dest state
	transitions map[T]T

	// actions called when state be transitioned
	actions map[actionKey[T]]Action[T]

	// async retry with async
	async bool

	// tick is retry delay duration, default is time.Second, only async
	tick time.Duration

	// wheelSize is a timingWheel delayQueue size, default is 1024, only async
	wheelSize int64

	handler func(T, error)

	mu, stateMu sync.Mutex
}

func WithAsync[T comparable]() Option[T] {
	return func(f *FSM[T]) {
		f.async = true
	}
}

func WithTransition[T comparable](src, dest T) Option[T] {
	return func(f *FSM[T]) {
		f.transitions[src] = dest
	}
}

func WithHandler[T comparable](handler func(state T, err error)) Option[T] {
	return func(f *FSM[T]) {
		f.handler = handler
	}
}

// NewFSM returnT a new finite state Machine with initial and transitions
func NewFSM[T comparable](initial T, opts ...Option[T]) *FSM[T] {
	fsm := &FSM[T]{
		current:     initial,
		transitions: make(map[T]T),
		actions:     make(map[actionKey[T]]Action[T]),
		tick:        time.Second,
		wheelSize:   1024,
	}
	for _, cb := range opts {
		cb(fsm)
	}
	return fsm
}

func (f *FSM[T]) State() T {
	f.stateMu.Lock()
	defer f.stateMu.Unlock()
	return f.current
}

func (f *FSM[T]) SetState(state T) *FSM[T] {
	f.stateMu.Lock()
	defer f.stateMu.Unlock()
	f.current = state
	return f
}

func (f *FSM[T]) Transition(ctx context.Context, args ...any) {
	f.mu.Lock()
	defer f.mu.Unlock()

	e := &Event[T]{ctx, f, args, nil}
	f.transition(e)
}

func (f *FSM[T]) transition(e *Event[T]) {
	if err := e.ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			err = errors.Join(CanceledErr[T]{f.State()}, err)
		}
		f.handlerErr(err)
		return
	}

	e.Err = nil
	ts := []ActionType{before, transition, after}
	for _, t := range ts {
		f.doAction(t, e)
		if e.Err != nil {
			f.handlerErr(e.Err)
			return
		}
	}

	if dest, ok := f.transitions[f.State()]; ok {
		f.SetState(dest)
		f.transition(e)
	}
}

func (f *FSM[T]) doAction(t ActionType, e *Event[T]) {
	if action, ok := f.actions[actionKey[T]{t, f.State()}]; ok {
		action(e.ctx, e)
	}
}

func (f *FSM[T]) handlerErr(err error) {
	if f.handler == nil {
		return
	}

	f.handler(f.State(), err)
}

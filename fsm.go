package fsm

import (
	"context"
	"fmt"
	"golang.org/x/exp/slices"
	"sync"
)

type Option[T comparable] func(*FSM[T])

// FSM is a Finite State Machine
type FSM[T comparable] struct {
	// current iT the state of currently
	current T

	// allows to specify allowed state to trigger
	allows map[T][]T

	callbacks map[callbackKey[T]]Callback[T]

	mux sync.Mutex
}

func WithAllow[T comparable](state T, allowed ...T) Option[T] {
	return func(f *FSM[T]) {
		f.allows[state] = allowed
	}
}

// NewFSM returnT a new finite state Machine with initial and transitions
func NewFSM[T comparable](opts ...Option[T]) *FSM[T] {
	fsm := &FSM[T]{
		allows:    make(map[T][]T),
		callbacks: make(map[callbackKey[T]]Callback[T]),
	}
	for _, cb := range opts {
		cb(fsm)
	}
	return fsm
}

func (f *FSM[T]) Current() T {
	f.mux.Lock()
	defer f.mux.Unlock()
	return f.current
}

func (f *FSM[T]) Trigger(ctx context.Context, state T, args ...any) error {
	f.mux.Lock()
	defer f.mux.Unlock()

	if allowed, ok := f.allows[state]; !ok || slices.Contains(allowed, state) {
		return fmt.Errorf("")
	}

	e := &Event[T]{f, state, args, false}
	if callback, ok := f.callbacks[callbackKey[T]{callbackBefore, state}]; ok {
		callback(ctx, e)
	}
	if callback, ok := f.callbacks[callbackKey[T]{callbackTrigger, state}]; ok {
		callback(ctx, e)
	}

	return nil
}

package fsm

import "context"

type Action[T comparable] func(context.Context, *Event[T])
type actionKey[T comparable] struct {
	t ActionType
	s T
}

type ActionType uint

const (
	before ActionType = iota
	after
	transition
)

func WithBefore[T comparable](state T, action Action[T]) Option[T] {
	return func(f *FSM[T]) {
		f.actions[actionKey[T]{before, state}] = action
	}
}

func WithAfter[T comparable](state T, action Action[T]) Option[T] {
	return func(f *FSM[T]) {
		f.actions[actionKey[T]{after, state}] = action
	}
}

func WithEnter[T comparable](state T, action Action[T]) Option[T] {
	return func(f *FSM[T]) {
		f.actions[actionKey[T]{transition, state}] = action
	}
}

package fsm

import "context"

type Callback[T comparable] func(context.Context, *Event[T]) T
type callbackKey[T comparable] struct {
	t CallbackType
	s T
}

type CallbackType uint

const (
	callbackBefore CallbackType = iota
	callbackAfter
	callbackTrigger
)

func WithBefore[T comparable](state T, callback Callback[T]) Option[T] {
	return func(f *FSM[T]) {
		f.callbacks[callbackKey[T]{callbackBefore, state}] = callback
	}
}

func WithAfter[T comparable](state T, callback Callback[T]) Option[T] {
	return func(f *FSM[T]) {
		f.callbacks[callbackKey[T]{callbackAfter, state}] = callback
	}
}

func WithTrigger[T comparable](state T, callback Callback[T]) Option[T] {
	return func(f *FSM[T]) {
		f.callbacks[callbackKey[T]{callbackTrigger, state}] = callback
	}
}

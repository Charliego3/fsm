package fsm

type Event[T comparable] struct {
	FSM      *FSM[T]
	State    T
	Args     []any
	Canceled bool
}

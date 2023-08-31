package fsm

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type Step uint

const (
	one Step = iota + 1
	two
	three
	four
	five
)

func TestGeneral(t *testing.T) {
	count := 0
	NewFSM(one,
		WithTransition(one, two),
		WithTransition(two, three),
		WithTransition(three, four),
		WithTransition(four, five),
		WithAsync[Step](),
		WithBefore(one, func(ctx context.Context, e *Event[Step]) {
			fmt.Println("one before")
		}),
		WithEnter(one, func(ctx context.Context, e *Event[Step]) {
			fmt.Println("one")
		}),
		WithAfter(one, func(ctx context.Context, e *Event[Step]) {
			fmt.Println("one after")
		}),
		WithEnter(two, func(ctx context.Context, e *Event[Step]) {
			fmt.Println("two")
			count++
			if count > 1 {
				return
			}
			e.Retry(time.Second * 5)
		}),
		WithEnter(three, func(ctx context.Context, e *Event[Step]) {
			fmt.Println("three")
			e.Cancel()
		}),
		WithEnter(four, func(ctx context.Context, e *Event[Step]) {
			fmt.Println("four")
		}),
		WithEnter(five, func(ctx context.Context, e *Event[Step]) {
			fmt.Println("five")
		}),
		WithHandler[Step](func(state Step, err error) {
			fmt.Printf("got err with state: %v, err: %v\n", state, err)
		}),
	).Transition(context.Background())

	time.Sleep(time.Second * 5)
}

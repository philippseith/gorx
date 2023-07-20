package rx

import "context"

type Subscribable[T any] interface {
	Subscribe(o Observer[T]) Subscription
}

func FromChan[T any](ctx context.Context, ch <-chan T) Subscribable[T] {
	s := NewSubject[T]()
	go func() {
		for {
			select {
			case t, ok := <-ch:
				if !ok {
					s.Complete()
					return
				}
				s.Next(t)
			case <-ctx.Done():
				s.Complete()
				return
			}
		}
	}()
	return s
}

func ToChan[T any](s Subscribable[T]) (<-chan Item[T], Subscription) {
	tc := &toChan[T]{
		ch: make(chan Item[T], 1),
	}
	return tc.ch, s.Subscribe(tc)
}

type Item[T any] struct {
	V T
	E error
}

type toChan[T any] struct {
	ch chan Item[T]
}

func (tc *toChan[T]) Next(value T) {
	tc.ch <- Item[T]{V: value}
}

func (tc *toChan[T]) Error(err error) {
	tc.ch <- Item[T]{E: err}
	close(tc.ch)
}

func (tc *toChan[T]) Complete() {
	close(tc.ch)
}

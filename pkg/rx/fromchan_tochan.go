package rx

import "context"

func FromChan[T any](ctx context.Context, ch <-chan T) Observable[T] {
	oo := &observableObserver[T]{}
	go func() {
		for {
			select {
			case t, ok := <-ch:
				if !ok {
					oo.Complete()
					return
				}
				oo.Next(t)
			case <-ctx.Done():
				oo.Complete()
				return
			}
		}
	}()
	return oo
}

func ToChan[T any](o Observable[T]) (<-chan Item[T], Subscription) {
	tc := &toChan[T]{
		ch: make(chan Item[T], 1),
	}
	return tc.ch, o.Subscribe(tc)
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

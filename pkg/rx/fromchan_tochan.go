package rx

import "context"

// FromChan creates an Observable[T] from a chan T. The channel will be read at
// once. All values sent before the Observable is subscribed to, will be ignored.
func FromChan[T any](ctx context.Context, ch <-chan T) Observable[T] {
	oo := &observableObserver[T, T]{
		t2u: func(t T) T {
			return t
		},
	}
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

// ToChan pushes the values from an Observable into a channel. It returns a
// channel and the Subscription to the Observable. the channel type Item[T]
// contains values and a possible error. Note that ToChan will block immediately
// with cold observables. You need to wrap the cold observable with
// ToConnectable and set up a receiving goroutine for the channel before you
// call Connectable.Connect
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

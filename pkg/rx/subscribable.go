package rx

import (
	"context"
)

type Subscribable[T any] interface {
	Subscribe(o Observer[T]) Subscription
}

func FromChan[T any](ctx context.Context, ch <-chan T) Subscribable[T] {
	p := &pipe[T, T]{
		n: func(t T) T {
			return t
		},
	}
	go func() {
		for {
			select {
			case t, ok := <-ch:
				if !ok {
					p.Complete()
					return
				}
				p.Next(t)
			case <-ctx.Done():
				p.Complete()
				return
			}
		}
	}()
	return p
}

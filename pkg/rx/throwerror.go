package rx

import "context"

func ThrowError[T any](ctx context.Context, err error) Observable[T] {
	return ToObservable[T](&throwError[T]{
		ctx: ctx,
		err: err,
	})
}

type throwError[T any] struct {
	ctx context.Context
	err error
}

func (te *throwError[T]) Subscribe(o Observer[T]) Subscription {
	o.Error(te.ctx, te.err)
	return NewSubscription(func() {})
}

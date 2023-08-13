package rx

func ThrowError[T any](err error) Observable[T] {
	return ToObservable[T](&throwError[T]{err: err})
}

type throwError[T any] struct {
	err error
}

func (te *throwError[T]) Subscribe(o Observer[T]) Subscription {
	o.Error(te.err)
	return NewSubscription(func() {})
}

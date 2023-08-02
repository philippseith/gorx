package rx

func ThrowError[T any](err error) Observable[T] {
	te := &throwError[T]{err: err}
	te.Subscribable = te
	return te
}

type throwError[T any] struct {
	observableObserver[T, T]
	err error
}

func (te *throwError[T]) Subscribe(o Observer[T]) Subscription {
	o.Error(te.err)
	return NewSubscription(func() {})
}

package rx

func Take[T any](s Subscribable[T], count int) Observable[T] {
	t := &take[T]{
		observableObserver: observableObserver[T, T]{
			t2u: func(t T) T { return t },
		},
		count: count,
	}
	t.sourceSub = func() Subscription {
		return s.Subscribe(t)
	}
	return t
}

type take[T any] struct {
	observableObserver[T, T]
	count int
}

func (t *take[T]) Next(value T) {
	t.observableObserver.Next(value)
	t.count--
	if t.count == 0 {
		t.Complete()
	}
}

package rx

func Take[T any](s Subscribable[T], count int) Observable[T] {
	t := &take[T]{
		Operator: Operator[T, T]{t2u: func(t T) T { return t }},
		count:    count,
	}
	t.prepareSubscribe(func() Subscription { return s.Subscribe(t) })
	return ToObservable[T](t)
}

type take[T any] struct {
	Operator[T, T]
	count int
}

func (t *take[T]) Next(value T) {
	if t.count != 0 {
		t.Operator.Next(value)
		t.count--

		if t.count == 0 {
			t.Complete()
		}
	}
}

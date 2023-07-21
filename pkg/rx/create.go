package rx

func Create[T any](subscribe func(o Observer[T])) Observable[T] {
	return &create[T]{s: subscribe}
}

type create[T any] struct {
	s func(Observer[T])
}

func (c *create[T]) Subscribe(o Observer[T]) Subscription {
	c.s(o)
	return &subscription{u: func() {}}
}

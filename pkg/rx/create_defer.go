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

func Defer[T any](factory func() Observable[T]) Observable[T] {
	return &deferImp[T]{f: factory}
}

type deferImp[T any] struct {
	f func() Observable[T]
}

func (d *deferImp[T]) Subscribe(o Observer[T]) Subscription {
	return d.f().Subscribe(o)
}

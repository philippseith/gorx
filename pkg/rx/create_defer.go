package rx

// Create creates an Observable from a subscribe function, which is called on every subscription
func Create[T any](subscribe func(o Observer[T])) Observable[T] {
	return &create[T]{s: subscribe}
}

type create[T any] struct {
	observable[T]
	s func(Observer[T])
}

func (c *create[T]) Subscribe(o Observer[T]) Subscription {
	c.s(o)
	return &subscription{}
}

// Defer creates an Observable for each subscription with the help of s factory function
func Defer[T any](factory func() Observable[T]) Observable[T] {
	return &deferImp[T]{f: factory}
}

type deferImp[T any] struct {
	observable[T]
	f func() Observable[T]
}

func (d *deferImp[T]) Subscribe(o Observer[T]) Subscription {
	return d.f().Subscribe(o)
}

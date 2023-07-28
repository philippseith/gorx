package rx

// Create creates an Observable from a subscribe function, which is called on every subscription
func Create[T any](subscribe func(o Observer[T]) Subscription) Observable[T] {
	c := &create[T]{s: subscribe}
	c.Subscribable = c // let our internal observable know that we are its subscribable
	return c
}

type create[T any] struct {
	observable[T]
	s func(Observer[T]) Subscription
}

func (c *create[T]) Subscribe(o Observer[T]) Subscription {
	return c.s(o)
}

// Defer creates an Observable for each subscription with the help of s factory function
func Defer[T any](factory func() Observable[T]) Observable[T] {
	d := &deferImp[T]{f: factory}
	d.Subscribable = d
	return d
}

type deferImp[T any] struct {
	observable[T]
	f func() Observable[T]
}

func (d *deferImp[T]) Subscribe(o Observer[T]) Subscription {
	return d.f().Subscribe(o)
}

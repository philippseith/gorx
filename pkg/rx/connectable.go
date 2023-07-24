package rx

type Connectable[T any] interface {
	Observable[T]

	Connect()
}

func ToConnectable[T any](o Observable[T]) Connectable[T] {
	return &connectable[T]{o: o}
}

type connectable[T any] struct {
	observableObserver[T]
	o Observable[T]
}

func (c *connectable[T]) Connect() {
	c.o.Subscribe(c)
}

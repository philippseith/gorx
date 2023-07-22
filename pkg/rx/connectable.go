package rx

type Connectable[T any] interface {
	Observable[T]

	Connect()
}

func ToConnectable[T any](o Observable[T]) Connectable[T] {
	return &connectable[T]{
		Subject: NewSubject[T](),
		o:       o,
	}
}

type connectable[T any] struct {
	Subject[T]
	o Observable[T]
}

func (c *connectable[T]) Connect() {
	c.o.Subscribe(c)
}

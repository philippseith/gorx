package rx

type Connectable[T any] interface {
	Observable[T]

	Connect()
}

// ToConnectable creates an observable that multicasts once connect() is called on it.
func ToConnectable[T any](s Subscribable[T]) Connectable[T] {
	return &connectable[T]{
		Subject: NewSubject[T](),
		s:       s}
}

type connectable[T any] struct {
	Subject[T]
	s Subscribable[T]
}

func (c *connectable[T]) Connect() {
	c.s.Subscribe(c)
}

package rx

type Connectable[T any] interface {
	Observable[T]

	Connect()
}

// ToConnectable creates an observable that multicasts once connect() is called on it.
func ToConnectable[T any](s Subscribable[T]) Connectable[T] {
	return &connectable[T]{
		observableObserver: observableObserver[T, T]{
			t2u: func(t T) T {
				return t
			},
		},
		s: s}
}

type connectable[T any] struct {
	observableObserver[T, T]
	s Subscribable[T]
}

func (c *connectable[T]) Connect() {
	c.s.Subscribe(c)
}

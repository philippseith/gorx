package rx

func Share[T any](o Observable[T]) Observable[T] {
	return &share[T]{
		subject: subject[T]{},
		o:       o,
	}
}

type share[T any] struct {
	subject[T]
	o Observable[T]
	s Subscription
}

func (s *share[T]) Subscribe(o Observer[T]) Subscription {
	su := s.subject.Subscribe(o)
	if len(s.subject.observers) == 1 {
		s.s = s.o.Subscribe(s)
	}
	return NewSubscription(func() {
		su.Unsubscribe()
		if len(s.subject.observers) == 0 {
			s.s.Unsubscribe()
		}
	})
}

package rx

func Share[T any](s Subscribable[T]) Observable[T] {
	return &share[T]{
		subject: subject[T]{},
		s:       s,
	}
}

type share[T any] struct {
	subject[T]
	s  Subscribable[T]
	sn Subscription
}

func (sh *share[T]) Subscribe(o Observer[T]) Subscription {
	su := sh.subject.Subscribe(o)
	if len(sh.subject.observers) == 1 {
		sh.sn = sh.s.Subscribe(sh)
	}
	su.AddTearDownLogic(func() {
		if len(sh.subject.observers) == 0 {
			sh.sn.Unsubscribe()
		}
	})
	return su
}

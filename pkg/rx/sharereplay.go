package rx

func ShareReplay[T any](s Subscribable[T], opts ...ReplayOption) Observable[T] {
	return &shareReplay[T]{
		ReplaySubject: *NewReplaySubject[T](opts...),
		s:             s,
	}
}

type shareReplay[T any] struct {
	ReplaySubject[T]
	s              Subscribable[T]
	sn             Subscription
	onceSubscribed bool
}

func (sh *shareReplay[T]) Subscribe(o Observer[T]) Subscription {
	su := sh.ReplaySubject.Subscribe(o)
	if (!sh.onceSubscribed || sh.opt.refCount) && len(sh.observers) == 1 {
		sh.sn = sh.s.Subscribe(sh)
		sh.onceSubscribed = true
	}
	su.AddTearDownLogic(func() {
		if sh.opt.refCount && len(sh.observers) == 0 {
			sh.sn.Unsubscribe()
		}
	})
	return su
}

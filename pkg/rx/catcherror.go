package rx

import "sync"

func CatchError[T any](s Subscribable[T], catch func(error) Subscribable[T]) Observable[T] {
	ce := &catchError[T]{
		observableObserver: observableObserver[T, T]{
			t2u: func(t T) T {
				return t
			},
		},
		catch: catch,
	}
	ce.mx.Lock()
	defer ce.mx.Unlock()

	ce.Subscribable = ce
	ce.sourceSub = func() Subscription { return s.Subscribe(ce) }
	return ce
}

type catchError[T any] struct {
	observableObserver[T, T]
	catch  func(error) Subscribable[T]
	errSub Subscription

	mx sync.RWMutex
}

func (ce *catchError[T]) Error(err error) {
	ce.mx.Lock()
	defer ce.mx.Unlock()

	ce.errSub = ce.catch(err).Subscribe(ce.o)
}

func (ce *catchError[T]) Subscribe(o Observer[T]) Subscription {
	return ce.observableObserver.Subscribe(o).
		AddTearDownLogic(
			func() {
				ce.mx.RLock()
				defer ce.mx.RUnlock()

				if ce.errSub != nil {
					ce.errSub.Unsubscribe()
				}
			})
}

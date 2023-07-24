package rx

import "sync"

type observableObserver[T any] struct {
	o  Observer[T]
	mx sync.RWMutex
}

func (oo *observableObserver[T]) Next(value T) {
	oo.mx.RLock()
	defer oo.mx.RUnlock()

	if oo.o != nil {
		oo.o.Next(value)
	}
}

func (oo *observableObserver[T]) Error(err error) {
	oo.mx.RLock()
	defer oo.mx.RUnlock()

	if oo.o != nil {
		oo.o.Error(err)
	}
}

func (oo *observableObserver[T]) Complete() {
	oo.mx.RLock()
	defer oo.mx.RUnlock()

	if oo.o != nil {
		oo.o.Complete()
	}
}

func (oo *observableObserver[T]) Subscribe(o Observer[T]) Subscription {
	oo.mx.Lock()
	defer oo.mx.Unlock()

	oo.o = o

	return NewSubscription(func() {
		oo.mx.Lock()
		defer oo.mx.Unlock()

		oo.o = nil
	})
}

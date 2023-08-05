package rx

import (
	"fmt"
	"runtime/debug"
	"sync"
)

type observableObserver[T any, U any] struct {
	observable[U]
	o         Observer[U]
	sourceSub func() Subscription
	mx        sync.RWMutex

	t2u func(T) U
}

func (oo *observableObserver[T, U]) Next(value T) {
	oo.mx.RLock()
	defer oo.mx.RUnlock()

	defer func() {
		if r := recover(); r != nil {
			if oo.o != nil {
				oo.o.Error(fmt.Errorf("panic in %T.Next(%v): %v.\n%s", oo.o, value, r, string(debug.Stack())))
			}
		}
	}()

	if oo.o != nil {
		oo.o.Next(oo.t2u(value))
	}
}

func (oo *observableObserver[T, U]) Error(err error) {
	func() {
		oo.mx.RLock()
		defer oo.mx.RUnlock()

		// no defer recover with sending to o.Error(), as this would build an endless loop

		if oo.o != nil {
			oo.o.Error(err)
		}
	}()
	oo.unsubscribe()
}

func (oo *observableObserver[T, U]) Complete() {
	func() {
		oo.mx.RLock()
		defer oo.mx.RUnlock()

		defer func() {
			if r := recover(); r != nil {
				if oo.o != nil {
					oo.o.Error(fmt.Errorf("panic in %T.Complete(): %v\n%s", oo.o, r, string(debug.Stack())))
				}
			}
		}()

		if oo.o != nil {
			oo.o.Complete()
		}
	}()
	oo.unsubscribe()
}

func (oo *observableObserver[T, U]) Subscribe(o Observer[U]) Subscription {
	var sourceSub Subscription
	if ss := func() func() Subscription {
		oo.mx.Lock()
		defer oo.mx.Unlock()

		oo.o = o
		return oo.sourceSub
	}(); ss != nil {
		sourceSub = ss()
	}
	sub := NewSubscription(oo.unsubscribe)
	if sourceSub != nil {
		sub.AddSubscription(sourceSub)
	}
	return sub
}

func (oo *observableObserver[T, U]) unsubscribe() {
	oo.mx.Lock()
	defer oo.mx.Unlock()

	oo.o = nil
}

package rx

import (
	"context"
	"sync"
)

// Catch catches errors on the Subscribable to be handled by returning a new
// Subscribable
func Catch[T any](s Subscribable[T], catchError func(context.Context, error) Subscribable[T]) Observable[T] {
	ce := &catch[T]{
		Operator: Operator[T, T]{t2u: func(t T) T { return t }},
		catch:    catchError,
	}
	ce.prepareSubscribe(func() Subscription {
		return s.Subscribe(ce).
			AddTearDownLogic(
				func() {
					ce.mx.RLock()
					defer ce.mx.RUnlock()

					if ce.errSub != nil {
						ce.errSub.Unsubscribe()
					}
				})
	})
	return ToObservable[T](ce)
}

type catch[T any] struct {
	Operator[T, T]
	catch  func(context.Context, error) Subscribable[T]
	errSub Subscription

	mx sync.RWMutex
}

func (ce *catch[T]) Error(ctx context.Context, err error) {
	ce.mx.Lock()
	defer ce.mx.Unlock()

	ce.errSub = ce.catch(ctx, err).Subscribe(ce.observer())
}

package rx

import (
	"sync"

	"golang.org/x/exp/slices"
)

// CombineLatest combines multiple Subscribables to create an Observable whose
// values are calculated from the latest values of each of its input Observables.
func CombineLatest[T any](combine func(next ...any) T, sources ...Subscribable[any]) Observable[T] {
	c := &combineLatest[T]{
		Operator:  Operator[any, T]{t2u: func(a any) T { return a.(T) }},
		lasts:     make([]any, len(sources)),
		subs:      make([]Subscription, len(sources)),
		completed: make([]bool, len(sources)),
	}
	c.prepareSubscribe(func() Subscription {
		for i, source := range sources {
			idx := i
			c.subs[idx] = source.Subscribe(NewObserver[any](
				// Next
				func(next any) {
					if func() bool {
						c.mx.Lock()
						defer c.mx.Unlock()

						c.lasts[idx] = next
						return slices.Contains(c.lasts, nil)
					}() {
						return
					}
					lasts := make([]any, len(c.lasts))
					copy(lasts, c.lasts)
					c.Operator.Next(combine(lasts...))
				},
				// Error
				c.Operator.Error,
				// Complete
				func() {
					if func() bool {
						c.mx.Lock()
						defer c.mx.Unlock()

						c.completed[idx] = true
						return slices.Contains(c.completed, false)
					}() {
						return
					}
					c.Operator.Complete()
				}))
		}
		return NewSubscription(func() {
			for _, sub := range c.subs {
				sub.Unsubscribe()
			}
		})
	})
	return ToObservable[T](c)
}

type combineLatest[T any] struct {
	Operator[any, T]
	lasts     []any
	subs      []Subscription
	completed []bool

	mx sync.RWMutex
}

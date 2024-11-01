package rx

import (
	"sync"
)

// CombineLatest combines multiple Subscribables to create an Observable whose
// values are calculated from the latest values of each of its input Observables.
func CombineLatest[T any](combine func(next ...any) T, sources ...Subscribable[any]) Observable[T] {
	c := &combineLatest[T]{
		Operator:  Operator[any, T]{t2u: func(a any) T { return a.(T) }},
		latest:    make([]any, len(sources)),
		subs:      make([]Subscription, len(sources)),
		completed: make([]bool, len(sources)),
	}
	c.prepareSubscribe(func() Subscription {
		for i, source := range sources {
			idx := i
			c.subs[idx] = source.Subscribe(NewObserver[any](
				// Next
				func(next any) { c.next(combine, idx, next) },
				// Error
				c.Operator.Error,
				// Complete
				func() { c.complete(idx) }))
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
	latest    []any
	subs      []Subscription
	completed []bool

	mx sync.RWMutex
}

func (c *combineLatest[T]) next(combine func(next ...any) T, idx int, next any) {
	if func() bool {
		c.mx.Lock()
		defer c.mx.Unlock()

		c.latest[idx] = next
		for _, last := range c.latest {
			if last == nil {
				return true
			}
		}
		return false
	}() {
		return
	}

	latest := func() []any {
		c.mx.Lock()
		defer c.mx.Unlock()

		latest := make([]any, len(c.latest))
		copy(latest, c.latest)
		return latest
	}()

	c.Operator.Next(combine(latest...))
}

func (c *combineLatest[T]) complete(idx int) {
	if func() bool {
		c.mx.Lock()
		defer c.mx.Unlock()

		c.completed[idx] = true
		for _, completed := range c.completed {
			if !completed {
				return true
			}
		}
		return false
	}() {
		return
	}
	c.Operator.Complete()
}

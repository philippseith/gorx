package rx

import "golang.org/x/exp/slices"

func CombineLatest[T any](combine func(next ...any) T, sources ...Observable[any]) Observable[T] {
	c := &combineLatest[T]{
		observableObserver: observableObserver[any, T]{
			t2u: func(a any) T { return a.(T) },
		},
		lasts: make([]any, len(sources)),
		subs:  make([]Subscription, len(sources)),
		compl: make([]bool, len(sources)),
	}
	c.sourceSub = func() {
		for i, source := range sources {
			idx := i
			c.subs[idx] = source.Subscribe(NewObserver[any](func(next any) {
				if func() bool {
					c.observableObserver.mx.Lock()
					defer c.observableObserver.mx.Unlock()

					c.lasts[idx] = next
					return slices.Contains(c.lasts, nil)
				}() {
					return
				}
				lasts := make([]any, len(c.lasts))
				copy(lasts, c.lasts)
				c.observableObserver.Next(combine(lasts...))
			}, c.observableObserver.Error, func() {
				if func() bool {
					c.observableObserver.mx.Lock()
					defer c.observableObserver.mx.Unlock()

					c.compl[idx] = true
					return slices.Contains(c.compl, false)
				}() {
					return
				}
				c.observableObserver.Complete()
			}))
		}
	}
	return c
}

type combineLatest[T any] struct {
	observableObserver[any, T]
	lasts []any
	subs  []Subscription
	compl []bool
}

func (c *combineLatest[T]) Subscribe(o Observer[T]) Subscription {
	subOo := c.observableObserver.Subscribe(o)
	return NewSubscription(func() {
		c.observableObserver.mx.RLock()
		defer c.observableObserver.mx.RUnlock()

		subOo.Unsubscribe()
		for _, sub := range c.subs {
			if sub != nil {
				sub.Unsubscribe()
			}
		}
	})
}

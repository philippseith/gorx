package rx

import "sync"

// Concat creates an output Observable which sequentially emits all values from
// the first given Observable and then moves on to the next.
func Concat[T any](sources ...Subscribable[T]) Observable[T] {
	c := &concat[T]{sources: sources}
	return ToObservable[T](c)
}

type concat[T any] struct {
	sub      Subscription
	sources  []Subscribable[T]
	o        Observer[T]
	mxState  sync.RWMutex
	mxEvents sync.Mutex
}

func (c *concat[T]) Next(t T) {
	if c.observer() != nil {
		func() {
			c.mxEvents.Lock()
			defer c.mxEvents.Unlock()

			c.observer().Next(t)
		}()
	}
}

func (c *concat[T]) Error(err error) {
	if c.observer() != nil {
		func() {
			c.mxEvents.Lock()
			defer c.mxEvents.Unlock()
			c.observer().Error(err)
		}()
	}
}

func (c *concat[T]) Complete() {
	if source := func() Subscribable[T] {
		c.mxState.Lock()
		defer c.mxState.Unlock()

		// Might be nil when the source completes immediately
		if c.sub != nil {
			c.sub.Unsubscribe()
		}
		c.sources = c.sources[1:]
		if len(c.sources) > 0 {
			return c.sources[0]
		} else if c.o != nil {
			func() {
				c.mxEvents.Lock()
				defer c.mxEvents.Unlock()

				c.o.Complete()
			}()
		}
		return nil
	}(); source != nil {
		sub := source.Subscribe(c)
		func() {
			c.mxState.Lock()
			defer c.mxState.Unlock()

			c.sub = sub
		}()
	}
}

func (c *concat[T]) Subscribe(o Observer[T]) Subscription {
	if source := func() Subscribable[T] {
		c.mxState.Lock()
		defer c.mxState.Unlock()

		c.o = o
		if len(c.sources) > 0 {
			return c.sources[0]
		}
		return nil
	}(); source != nil {
		sub := source.Subscribe(c)
		func() {
			c.mxState.Lock()
			defer c.mxState.Unlock()

			c.sub = sub
		}()
	}
	return NewSubscription(func() {
		c.o = nil
		if c.sub != nil {
			c.sub.Unsubscribe()
		}
	})
}

func (c *concat[T]) observer() Observer[T] {
	c.mxState.RLock()
	defer c.mxState.RUnlock()

	return c.o
}

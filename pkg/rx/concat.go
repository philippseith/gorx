package rx

import "sync"

// Concat creates an output Observable which sequentially emits all values from
// the first given Observable and then moves on to the next.
func Concat[T any](sources ...Subscribable[T]) Observable[T] {
	c := &concat[T]{sources: sources}
	if len(c.sources) > 0 {
		c.sub = c.sources[0].Subscribe(c)
	}
	return ToObservable[T](c)
}

type concat[T any] struct {
	sub     Subscription
	sources []Subscribable[T]
	o       Observer[T]
	mx      sync.RWMutex
}

func (c *concat[T]) Next(t T) {
	if c.observer() != nil {
		c.observer().Next(t)
	}
}

func (c *concat[T]) Error(err error) {
	if c.observer() != nil {
		c.observer().Error(err)
	}
}

func (c *concat[T]) Complete() {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.sub.Unsubscribe()
	c.sources = c.sources[1:]
	if len(c.sources) > 0 {
		c.sub = c.sources[0].Subscribe(c)
	} else if c.o != nil {
		c.o.Complete()
	}
}

func (c *concat[T]) Subscribe(o Observer[T]) Subscription {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.o = o
	return NewSubscription(func() {
		c.o = nil
		if c.sub != nil {
			c.sub.Unsubscribe()
		}
	})
}

func (c *concat[T]) observer() Observer[T] {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.o
}

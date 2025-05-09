package rx

import (
	"sync"
)

// Concat creates an output Observable which sequentially emits all values from
// the first given Observable and then moves on to the next.
func Concat[T any](sources ...Subscribable[T]) Observable[T] {
	c := &concat[T]{sources: sources}
	return ToObservable(c)
}

type concat[T any] struct {
	sub           Subscription
	subInComplete Subscription
	sources       []Subscribable[T]
	o             Observer[T]
	mxState       sync.RWMutex
	mxEvents      sync.Mutex
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
	c.unsubscribeAndNilSubscription()

	func() {
		c.mxState.Lock()
		defer c.mxState.Unlock()

		c.sources = c.sources[1:]
	}()

	if source := c.currentSource(); source == nil {
		c.sendComplete()
	} else {
		c.setSubscriptionInComplete(source.Subscribe(c))
	}
}

func (c *concat[T]) Subscribe(o Observer[T]) Subscription {
	source := c.currentSource()
	if source == nil {
		return NewSubscription(nil)
	}
	c.setObserver(o)
	sub := source.Subscribe(c)
	// Check if source has completed in Subscribe
	switch c.currentSource() {
	case source:
		// It has not completed
		c.setSourceSubscription(sub)

		return NewSubscription(func() {
			c.setObserver(nil)
			c.unsubscribeAndNilSubscription()
		})
	case nil:
		// It has completed and it was the last one
		return NewSubscription(nil)
	default:
		// It has completed but there has a next one been subscribed in Complete
		c.setSourceSubscriptionToSubscriptionInComplete()

		return NewSubscription(func() {
			c.setObserver(nil)
			c.unsubscribeAndNilSubscription()
		})
	}
}

func (c *concat[T]) currentSource() Subscribable[T] {
	c.mxState.RLock()
	defer c.mxState.RUnlock()
	if len(c.sources) > 0 {
		return c.sources[0]
	}
	return nil
}

func (c *concat[T]) observer() Observer[T] {
	c.mxState.RLock()
	defer c.mxState.RUnlock()

	return c.o
}

func (c *concat[T]) setObserver(o Observer[T]) {
	c.mxState.Lock()
	defer c.mxState.Unlock()

	c.o = o
}

func (c *concat[T]) setSourceSubscription(sub Subscription) {
	c.mxState.Lock()
	defer c.mxState.Unlock()

	c.sub = sub
}

func (c *concat[T]) setSubscriptionInComplete(sub Subscription) {
	c.mxState.Lock()
	defer c.mxState.Unlock()

	c.subInComplete = sub
}

func (c *concat[T]) setSourceSubscriptionToSubscriptionInComplete() {
	c.mxState.Lock()
	defer c.mxState.Unlock()

	c.sub = c.subInComplete
}

func (c *concat[T]) unsubscribeAndNilSubscription() {
	c.mxState.Lock()
	defer c.mxState.Unlock()

	if c.sub != nil {
		c.sub.Unsubscribe()
		c.sub = nil
	}
}

func (c *concat[T]) sendComplete() {
	c.mxEvents.Lock()
	defer c.mxEvents.Unlock()
	if o := c.observer(); o != nil {
		o.Complete()
	}
}

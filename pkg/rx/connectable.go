package rx

import (
	"sync"
	"time"
)

type Connectable[T any] interface {
	Subscribable[T]
	ObservableExtension[Connectable[T], Subscribable[T], T]

	Connect()
}

// ToConnectable creates an observable that multicasts once connect() is called on it.
func ToConnectable[T any](s Subscribable[T]) Connectable[T] {
	c := &connectable[T]{source: s}
	c.Subscribable = c
	return c
}

type connectable[T any] struct {
	observable[T]
	source             Subscribable[T]
	sourceSubscription Subscription
	outerObserver      Observer[T]
	connected          bool
	mx                 sync.RWMutex
}

func (c *connectable[T]) Connect() {
	if func() bool {
		c.mx.Lock()
		defer c.mx.Unlock()

		if !c.connected {
			c.connected = true
			return true
		}
		return false
	}() {
		sub := func() Subscription {
			c.mx.RLock()
			defer c.mx.RUnlock()

			return c.source.Subscribe(c.outerObserver)
		}()
		func() {
			c.mx.Lock()
			defer c.mx.Unlock()

			c.sourceSubscription = sub
		}()
	}
}

func (c *connectable[T]) Subscribe(o Observer[T]) Subscription {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.outerObserver = o
	return NewSubscription(func() {
		c.mx.RLock()
		defer c.mx.RUnlock()

		if c.sourceSubscription != nil {
			c.sourceSubscription.Unsubscribe()
		}
	})
}

func (c *connectable[T]) AddTearDownLogic(tld func()) Connectable[T] {
	c.tdls = append(c.tdls, tld)
	return c
}

func (c *connectable[T]) Catch(catch func(error) Connectable[T]) Connectable[T] {
	panic("")
}

func (c *connectable[T]) Concat(...Connectable[T]) Connectable[T] {
	// TODO the Connectables need to be connected when their predecessor has completed
	panic("")
}

func (c *connectable[T]) DebounceTime(duration time.Duration) Connectable[T] {
	// TODO this is wrong. c needs to be connected the moment the result of ToConnectable connects,
	// otherwise the Subscribe call chain of the operators will stop at c
	return ToConnectable[T](DebounceTime[T](c, duration))
}

func (c *connectable[T]) DistinctUntilChanged(equal func(T, T) bool) Connectable[T] {

}

func (c *connectable[T]) Log(id string) Connectable[T] {

}

func (c *connectable[T]) Merge(sources ...Connectable[T]) Connectable[T] {

}

func (c *connectable[T]) Share() Connectable[T] {

}

func (c *connectable[T]) ShareReplay(opts ...ReplayOption) Connectable[T] {

}

func (c *connectable[T]) Take(count int) Connectable[T] {

}

func (c *connectable[T]) Tap(subscribe func(Observer[T]), next func(T) T, err func(error) error, complete, unsubscribe func()) Connectable[T] {

}

func (c *connectable[T]) ToSlice() <-chan []T {

}

package rx

import "sync"

type Connectable[T any] interface {
	Observable[T]

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

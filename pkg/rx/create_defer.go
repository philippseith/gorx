package rx

import (
	"context"
	"fmt"
	"runtime/debug"
)

// Create creates an Observable from a subscribe function, which is called on every subscribe.
// The function should register the passed observer as subscriber and return the Subscription for this.
// The context passed to the observer can be enriched with context.WithValue(rx.ContextKeyCreate, ...)
func Create[T any](subscribe func(o Observer[T]) Subscription) Observable[T] {
	c := &create[T]{s: subscribe}
	c.Subscribable = c
	return c
}

type create[T any] struct {
	observable[T]
	s func(Observer[T]) Subscription
}

func (c *create[T]) Subscribe(o Observer[T]) Subscription {
	defer func() {
		if r := recover(); r != nil {
			o.Error(context.Background(), fmt.Errorf("panic in Create(): %#v\n%s", r, string(debug.Stack())))
		}
	}()

	return c.s(o)
}

// Defer creates an Observable for each call to Subscribe with the help of a factory function.
// The context passed to the observer can be enriched with context.WithValue(rx.ContextKeyDefer, ...)
func Defer[T any](factory func() Observable[T]) Observable[T] {
	d := &deferImp[T]{f: factory}
	d.Subscribable = d
	return d
}

type deferImp[T any] struct {
	observable[T]
	f func() Observable[T]
}

func (d *deferImp[T]) Subscribe(o Observer[T]) Subscription {
	defer func() {
		if r := recover(); r != nil {
			o.Error(context.Background(), fmt.Errorf("panic in Defer(): %#v\n%s", r, string(debug.Stack())))
		}
	}()

	return d.f().Subscribe(o)
}

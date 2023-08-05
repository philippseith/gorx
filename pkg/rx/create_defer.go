package rx

import (
	"fmt"
	"runtime/debug"
)

// Create creates an Observable from a subscribe function, which is called on every subscription
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
			o.Error(fmt.Errorf("panic in Create(): %v\n%s", r, string(debug.Stack())))
		}
	}()

	return c.s(o)
}

// Defer creates an Observable for each subscription with the help of s factory function
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
			o.Error(fmt.Errorf("panic in Defer(): %v\n%s", r, string(debug.Stack())))
		}
	}()

	return d.f().Subscribe(o)
}

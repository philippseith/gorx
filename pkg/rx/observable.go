package rx

import (
	"time"
)

type Subscribable[T any] interface {
	Subscribe(o Observer[T]) Subscription
}

type Observable[T any] interface {
	Subscribable[T]

	AddTearDownLogic(tld func()) Observable[T]
	CatchError(catch func(error) Subscribable[T]) Observable[T]
	DebounceTime(duration time.Duration) Observable[T]
	DistinctUntilChanged(equal func(T, T) bool) Observable[T]
	Share() Observable[T]
	ShareReplay(opts ...ReplayOption) Observable[T]
	Take(count int) Observable[T]
	Tap(next func(T), err func(error), complete func()) Observable[T]
	ToAny() Observable[any]
	ToConnectable() Connectable[T]
	ToSlice() <-chan []T
}

func ToObservable[T any](s Subscribable[T]) Observable[T] {
	return &observable[T]{Subscribable: s}
}

// observable is an internal type to enrich simple Subscribable objects with the
// Observable interface. Note that structs containing observable and implementing
// Subscribable need to initialize observable.Subscribable with themselves
type observable[T any] struct {
	Subscribable[T]
	tdls []func()
}

func (o *observable[T]) Subscribe(or Observer[T]) Subscription {
	sub := o.Subscribable.Subscribe(or)
	for _, tld := range o.tdls {
		sub.AddTearDownLogic(tld)
	}
	return sub
}

func (o *observable[T]) AddTearDownLogic(tld func()) Observable[T] {
	o.tdls = append(o.tdls, tld)
	return o
}

func (o *observable[T]) CatchError(catch func(error) Subscribable[T]) Observable[T] {
	return CatchError[T](o, catch)
}

func (o *observable[T]) DebounceTime(duration time.Duration) Observable[T] {
	return DebounceTime[T](o, duration)
}

func (o *observable[T]) DistinctUntilChanged(equal func(T, T) bool) Observable[T] {
	return DistinctUntilChanged[T](o, equal)
}

func (o *observable[T]) Share() Observable[T] {
	return Share[T](o)
}

func (o *observable[T]) ShareReplay(opts ...ReplayOption) Observable[T] {
	return ShareReplay[T](o, opts...)
}

func (o *observable[T]) Take(count int) Observable[T] {
	return Take[T](o, count)
}

func (o *observable[T]) Tap(next func(T), err func(error), complete func()) Observable[T] {
	return Tap[T](o, next, err, complete)
}

func (o *observable[T]) ToAny() Observable[any] {
	return ToAny[T](o)
}

func (o *observable[T]) ToConnectable() Connectable[T] {
	return ToConnectable[T](o)
}

func (o *observable[T]) ToSlice() <-chan []T {
	return ToSlice[T](o)
}

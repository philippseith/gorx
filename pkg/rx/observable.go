package rx

import (
	"time"
)

// Subscribable is an object which can be subscribed to.
//
//	Subscribe(o Observer[T]) Subscription
//
// Subscribe registers an Observer for the stream of events. All Observables
// returned by operators subscribe to their sources when their own Subscribe
// method is called
type Subscribable[T any] interface {
	Subscribe(o Observer[T]) Subscription
}

// Observable attaches all operators with only on generic parameter on a Subscribable
type Observable[T any] interface {
	Subscribable[T]

	AddTearDownLogic(tld func()) Observable[T]
	Catch(catch func(error) Subscribable[T]) Observable[T]
	Concat(...Subscribable[T]) Observable[T]
	DebounceTime(duration time.Duration) Observable[T]
	DistinctUntilChanged(equal func(T, T) bool) Observable[T]
	Share() Observable[T]
	ShareReplay(opts ...ReplayOption) Observable[T]
	Take(count int) Observable[T]
	Tap(next func(T) T, err func(error) error, complete func()) Observable[T]
	ToAny() Observable[any]
	ToConnectable() Connectable[T]
	ToSlice() <-chan []T
}

// ToObservable extends a Subscribable to an Observable
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

func (o *observable[T]) Catch(catch func(error) Subscribable[T]) Observable[T] {
	return Catch[T](o, catch)
}

func (o *observable[T]) Concat(sources ...Subscribable[T]) Observable[T] {

	return Concat[T](append([]Subscribable[T]{o}, sources...)...)
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

func (o *observable[T]) Tap(next func(T) T, err func(error) error, complete func()) Observable[T] {
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

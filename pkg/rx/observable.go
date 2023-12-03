package rx

import (
	"context"
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
	ObservableExtension[Observable[T], Subscribable[T], T]

	ToAny() Observable[any]
	ToConnectable() Connectable[T]
}

type ObservableExtension[Extended Subscribable[T], Extendable Subscribable[T], T any] interface {
	AddTearDownLogic(tld func()) Extended
	Catch(catch func(context.Context, error) Extendable) Extended
	Concat(...Extendable) Extended
	DebounceTime(duration time.Duration) Extended
	DistinctUntilChanged(equal func(T, T) bool) Extended
	Log(id string) Extended
	Merge(sources ...Extendable) Extended
	Share() Extended
	ShareReplay(opts ...ReplayOption) Extended
	Take(count int) Extended
	Tap(subscribe func(Observer[T]),
		next func(context.Context, T) (context.Context, T),
		err func(context.Context, error) (context.Context, error),
		complete func(context.Context) context.Context,
		unsubscribe func()) Extended
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

func (o *observable[T]) Catch(catch func(context.Context, error) Subscribable[T]) Observable[T] {
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

func (o *observable[T]) Log(id string) Observable[T] {
	return Log[T](o, id)
}

func (o *observable[T]) Merge(sources ...Subscribable[T]) Observable[T] {
	sources = append([]Subscribable[T]{o}, sources...)
	return Merge[T](sources...)
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

func (o *observable[T]) Tap(
	subscribe func(Observer[T]),
	next func(context.Context, T) (context.Context, T),
	err func(context.Context, error) (context.Context, error),
	complete func(context.Context) context.Context,
	unsubscribe func()) Observable[T] {
	return Tap[T](o, subscribe, next, err, complete, unsubscribe)
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

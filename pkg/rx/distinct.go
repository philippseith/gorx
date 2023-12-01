package rx

import (
	"context"
	"reflect"
)

// Distinct returns an Observable that emits all items emitted by the source
// Subscribable that are distinct by comparison from previous items.
func Distinct[T comparable](s Subscribable[T]) Observable[T] {
	values := map[T]struct{}{}
	d := &Operator[T, T]{t2u: func(t T) T { return t }}
	d.prepareSubscribe(func() Subscription {
		return s.Subscribe(NewObserverWithContext[T](func(ctx context.Context, value T) {
			if _, in := values[value]; !in {
				d.Next(ctx, value)
			}
			values[value] = struct{}{}
		}, d.Error, d.Complete))
	})
	return ToObservable[T](d)
}

// DistinctUntilChanged returns a Observable that emits all values pushed by the
// source Subscribable if they are distinct in comparison to the last value the
// result observable emitted.
func DistinctUntilChanged[T any](s Subscribable[T], equal func(T, T) bool) Observable[T] {
	var last *T
	d := &Operator[T, T]{t2u: func(t T) T { return t }}
	if equal == nil {
		equal = func(a, b T) bool { return reflect.DeepEqual(a, b) }
	}
	d.prepareSubscribe(func() Subscription {
		return s.Subscribe(NewObserverWithContext[T](func(ctx context.Context, value T) {
			if last == nil || !equal(*last, value) {
				d.Next(ctx, value)
			}
			last = &value
		}, d.Error, d.Complete))
	})
	return ToObservable[T](d)
}

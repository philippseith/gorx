package rx

import "reflect"

// Distinct returns an Observable that emits all items emitted by the source
// Subscribable that are distinct by comparison from previous items.
func Distinct[T comparable](s Subscribable[T]) Observable[T] {
	values := map[T]struct{}{}
	oo := &observableObserver[T, T]{
		t2u: func(t T) T {
			return t
		},
	}
	oo.Subscribable = oo
	oo.sourceSub = func() Subscription {
		return s.Subscribe(NewObserver[T](func(value T) {
			if _, in := values[value]; !in {
				oo.Next(value)
			}
			values[value] = struct{}{}
		}, oo.Error, oo.Complete))
	}
	return oo
}

// DistinctUntilChanged returns a Observable that emits all values pushed by the
// source Subscribable if they are distinct in comparison to the last value the
// result observable emitted.
func DistinctUntilChanged[T any](s Subscribable[T], equal func(T, T) bool) Observable[T] {
	var last *T
	oo := &observableObserver[T, T]{
		t2u: func(t T) T {
			return t
		},
	}
	oo.Subscribable = oo
	if equal == nil {
		equal = func(a, b T) bool { return reflect.DeepEqual(a, b) }
	}
	oo.sourceSub = func() Subscription {
		return s.Subscribe(NewObserver[T](func(value T) {
			if last == nil || !equal(*last, value) {
				last = &value
				oo.Next(value)
			}
		}, oo.Error, oo.Complete))
	}
	return oo
}

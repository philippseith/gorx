package rx

import "reflect"

func Distinct[T comparable](o Observable[T]) Observable[T] {
	values := map[T]struct{}{}
	oo := &observableObserver[T]{}
	o.Subscribe(NewObserver[T](func(value T) {
		if _, in := values[value]; !in {
			oo.Next(value)
		}
		values[value] = struct{}{}
	}, oo.Error, oo.Complete))
	return oo
}

func DistinctUntilChanged[T any](o Observable[T], equal func(T, T) bool) Observable[T] {
	var last *T
	oo := &observableObserver[T]{}
	if equal == nil {
		equal = func(a, b T) bool { return reflect.DeepEqual(a, b) }
	}
	o.Subscribe(NewObserver[T](func(value T) {
		if last == nil || !equal(*last, value) {
			oo.Next(value)
		}
		last = &value
	}, oo.Error, oo.Complete))
	return oo
}

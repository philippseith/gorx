package rx

import "reflect"

func Distinct[T comparable](o Observable[T]) Observable[T] {
	values := map[T]struct{}{}
	s := NewSubject[T]()
	o.Subscribe(NewObserver[T](func(value T) {
		if _, in := values[value]; !in {
			s.Next(value)
		}
		values[value] = struct{}{}
	}, s.Error, s.Complete))
	return s
}

func DistinctUntilChanged[T any](o Observable[T], equal func(T, T) bool) Observable[T] {
	var last *T
	s := NewSubject[T]()
	if equal == nil {
		equal = func(a, b T) bool { return reflect.DeepEqual(a, b) }
	}
	o.Subscribe(NewObserver[T](func(value T) {
		if last == nil || !equal(*last, value) {
			s.Next(value)
		}
		last = &value
	}, s.Error, s.Complete))
	return s
}

package rx

import (
	"context"
	"time"
)

type Subscribable[T any] interface {
	Subscribe(o Observer[T]) Subscription
}

type Observable[T any] interface {
	Subscribable[T]

	Debounce(duration time.Duration) Observable[T]
	DistinctUntilChanged(equal func(T, T) bool) Observable[T]
	Share() Observable[T]
	Take(count int) Observable[T]
	ToAny() Observable[any]
	ToConnectable() Connectable[T]
	ToSlice(ctx context.Context) []T
}

type observable[T any] struct {
	s Subscribable[T]
}

func (o *observable[T]) Subscribe(or Observer[T]) Subscription {
	return o.s.Subscribe(or)
}

func (o *observable[T]) Debounce(duration time.Duration) Observable[T] {
	return Debounce[T](o, duration)
}

func (o *observable[T]) DistinctUntilChanged(equal func(T, T) bool) Observable[T] {
	return DistinctUntilChanged[T](o, equal)
}

func (o *observable[T]) Share() Observable[T] {
	return Share[T](o)
}

func (o *observable[T]) Take(count int) Observable[T] {
	return Take[T](o, count)
}

func (o *observable[T]) ToAny() Observable[any] {
	return ToAny[T](o)
}

func (o *observable[T]) ToConnectable() Connectable[T] {
	return ToConnectable[T](o)
}

func (o *observable[T]) ToSlice(ctx context.Context) []T {
	return ToSlice[T](ctx, o)
}

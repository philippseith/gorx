package rx

import (
	"context"
	"fmt"
	"runtime/debug"
)

// From creates an Observable that emits all items and then completes
func From[T any](items ...T) Observable[T] {
	f := &from[T]{items: items}
	f.Subscribable = f
	return f
}

type from[T any] struct {
	observable[T]
	items []T
}

func (f *from[T]) Subscribe(o Observer[T]) Subscription {
	defer func() {
		if r := recover(); r != nil {
			o.Error(context.Background(), fmt.Errorf("panic in Subscribe(): %v\n%s", r, string(debug.Stack())))
		}
	}()

	for _, item := range f.items {
		o.Next(context.Background(), item)
	}
	o.Complete(context.Background())
	return &subscription{}
}

// Just creates an Observable that emits only the value and then completes
func Just[T any](value T) Observable[T] {
	j := &just[T]{value: value}
	j.Subscribable = j
	return j
}

type just[T any] struct {
	observable[T]
	value T
}

func (j *just[T]) Subscribe(o Observer[T]) Subscription {
	defer func() {
		if r := recover(); r != nil {
			o.Error(context.Background(), fmt.Errorf("panic in Subscribe(): %v\n%s", r, string(debug.Stack())))
		}
	}()

	o.Next(context.Background(), j.value)
	o.Complete(context.Background())
	return &subscription{}
}

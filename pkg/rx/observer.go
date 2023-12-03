package rx

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
)

type Observer[T any] interface {
	Next(ctx context.Context, value T)
	Error(ctx context.Context, err error)
	Complete(ctx context.Context)
}

func NewObserver[T any](next func(T), err func(error), complete func()) Observer[T] {
	return &observer[T]{
		next:     func(_ context.Context, value T) { next(value) },
		err:      func(_ context.Context, e error) { err(e) },
		complete: func(_ context.Context) { complete() },
	}
}

func NewObserverWithContext[T any](next func(context.Context, T), err func(context.Context, error), complete func(context.Context)) Observer[T] {
	return &observer[T]{
		next:     next,
		err:      err,
		complete: complete,
	}
}

func OnNext[T any](next func(T)) Observer[T] {
	return &observer[T]{next: func(_ context.Context, value T) { next(value) }, err: func(_ context.Context, err error) { log.Print(err) }}
}

func OnNextWithContext[T any](next func(context.Context, T)) Observer[T] {
	return &observer[T]{next: next, err: func(_ context.Context, err error) { log.Print(err) }}
}

type observer[T any] struct {
	next     func(context.Context, T)
	err      func(context.Context, error)
	complete func(context.Context)
}

func (o *observer[T]) Next(ctx context.Context, value T) {
	defer func() {
		if r := recover(); r != nil {
			o.Error(ctx, fmt.Errorf("panic in %T.Next(%v): %v.\n%s", o, value, r, string(debug.Stack())))
		}
	}()

	if o.next != nil {
		o.next(ctx, value)
	}
}

func (o *observer[T]) Error(ctx context.Context, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in %T.Error(%v): %v.\n%s", o, err, r, string(debug.Stack()))
		}
	}()

	if o.err != nil {
		o.err(ctx, err)
	}
}

func (o *observer[T]) Complete(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			o.Error(ctx, fmt.Errorf("panic in %T.Complete(): %v\n%s", o, r, string(debug.Stack())))
		}
	}()

	if o.complete != nil {
		o.complete(ctx)
	}
}

package rx

import (
	"context"
	"sync"
)

// BehaviorSubject is a variant of Subject that requires an initial value and
// emits its current value whenever it is subscribed to.
type BehaviorSubject[T any] struct {
	Subject[T]
	ctx   context.Context
	value T
	mx    sync.RWMutex
}

// NewBehaviorSubject creates a new BehaviorSubject
func NewBehaviorSubject[T any](ctx context.Context, value T) *BehaviorSubject[T] {
	return &BehaviorSubject[T]{
		Subject: NewSubject[T](),
		ctx:     ctx,
		value:   value}
}

func (bs *BehaviorSubject[T]) Value() T {
	bs.mx.RLock()
	defer bs.mx.RUnlock()

	return bs.value
}

func (bs *BehaviorSubject[T]) Subscribe(o Observer[T]) Subscription {
	s := bs.Subject.Subscribe(o)

	o.Next(func() (context.Context, T) {
		bs.mx.RLock()
		defer bs.mx.RUnlock()

		return bs.ctx, bs.value
	}())

	return s
}

func (bs *BehaviorSubject[T]) Next(ctx context.Context, value T) {
	func(ctx context.Context, value T) {
		bs.mx.Lock()
		defer bs.mx.Unlock()

		bs.ctx = ctx
		bs.value = value
	}(ctx, value)

	bs.Subject.Next(ctx, value)
}

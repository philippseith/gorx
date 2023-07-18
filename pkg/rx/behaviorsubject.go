package rx

import "sync"

type BehaviorSubject[T any] struct {
	Subject[T]
	value T
	mx    sync.RWMutex
}

func NewBehaviorSubject[T any](value T) *BehaviorSubject[T] {
	return &BehaviorSubject[T]{value: value}
}

func (bs *BehaviorSubject[T]) Value() T {
	bs.mx.RLock()
	defer bs.mx.RUnlock()

	return bs.value
}

func (bs *BehaviorSubject[T]) Subscribe(o Observer[T]) Subscription {
	s := bs.Subject.Subscribe(o)

	o.Next(func() T {
		bs.mx.RLock()
		defer bs.mx.RUnlock()

		return bs.value
	}())

	return s
}

func (bs *BehaviorSubject[T]) Next(value T) {
	func(value T) {
		bs.mx.Lock()
		defer bs.mx.Unlock()

		bs.value = value
	}(value)

	bs.Subject.Next(value)
}

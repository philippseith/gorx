package rx

import "sync"

type AsyncSubject[T any] struct {
	Subject[T]
	value T
	mx    sync.RWMutex
}

func (as *AsyncSubject[T]) Next(value T) {
	func(value T) {
		as.mx.Lock()
		defer as.mx.Unlock()

		as.value = value
	}(value)
}

func (as *AsyncSubject[T]) Complete() {
	as.Subject.Next(func() T {
		as.mx.RLock()
		defer as.mx.RUnlock()

		return as.value
	}())
	as.Subject.Complete()
}

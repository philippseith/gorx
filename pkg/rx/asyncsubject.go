package rx

import "sync"

// AsyncSubject is a variant of Subject that only emits a value when it
// completes. It will emit its latest value to all its observers on completion.
type AsyncSubject[T any] struct {
	Subject[T]
	value T
	mx    sync.RWMutex
}

func NewAsyncSubject[T any]() *AsyncSubject[T] {
	return &AsyncSubject[T]{}
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

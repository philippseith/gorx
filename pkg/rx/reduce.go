package rx

import "sync"

func Reduce[T any, U any](s Subscribable[T], acc func(U, T) U) Subscribable[U] {
	r := &reduce[T, U]{acc: acc}
	s.Subscribe(r)
	return r
}

type reduce[T any, U any] struct {
	Subject[U]
	acc   func(U, T) U
	value U
	mx    sync.RWMutex
}

func (r *reduce[T, U]) Next(value T) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.value = r.acc(r.value, value)
}

func (r *reduce[T, U]) Complete() {
	r.Subject.Next(func() U {
		r.mx.RLock()
		defer r.mx.RUnlock()

		return r.value
	}())
	r.Subject.Complete()
}

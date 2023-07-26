package rx

import (
	"sync"
)

type Subject[T any] interface {
	Observable[T]
	Observer[T]
}

func NewSubject[T any]() Subject[T] {
	return &subject[T]{}
}

type subject[T any] struct {
	observable[T]
	observers []Observer[T]
	mx        sync.RWMutex
}

func (s *subject[T]) Subscribe(o Observer[T]) Subscription {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.observers = append(s.observers, o)

	return &subscription{u: func() {
		s.mx.Lock()
		defer s.mx.Unlock()

		for i, v := range s.observers {
			if o == v {
				s.observers = append(s.observers[:i], s.observers[i+1:]...)
				return
			}
		}
	}}
}

func (s *subject[T]) Next(value T) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	for _, o := range s.observers {
		o.Next(value)
	}
}

func (s *subject[T]) Error(err error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	for _, o := range s.observers {
		o.Error(err)
	}
}

func (s *subject[T]) Complete() {
	s.mx.RLock()
	defer s.mx.RUnlock()

	for _, o := range s.observers {
		o.Complete()
	}
}

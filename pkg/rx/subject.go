package rx

import (
	"sync"
)

type Subject[T any] interface {
	Observable[T]
	Observer[T]
}

func NewSubject[T any]() Subject[T] {
	s := &subject[T]{}

	s.mx.Lock()
	defer s.mx.Unlock()

	s.observers = map[int]Observer[T]{}
	return s
}

type subject[T any] struct {
	observers map[int]Observer[T]
	nextId    int
	mx        sync.RWMutex
}

func (s *subject[T]) Subscribe(o Observer[T]) Subscription {
	s.mx.Lock()
	defer s.mx.Unlock()

	id := s.nextId
	s.nextId++

	s.observers[id] = o

	return &subscription{u: func() {
		s.mx.Lock()
		defer s.mx.Unlock()

		delete(s.observers, id)
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

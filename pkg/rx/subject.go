package rx

import (
	"sync"

	"golang.org/x/exp/slices"
)

type Subject[T any] struct {
	observers []Observer[T]
	mx        sync.RWMutex
}

func (s *Subject[T]) Subscribe(o Observer[T]) Subscription {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.observers = append(s.observers, o)
	return &subscription{u: func() {
		s.mx.Lock()
		defer s.mx.Unlock()

		if idx := slices.Index(s.observers, o); idx != -1 {
			s.observers = append(s.observers[:idx], s.observers[idx+1:]...)
		}
	}}
}

func (s *Subject[T]) Next(value T) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	for _, o := range s.observers {
		o.Next(value)
	}
}

func (s *Subject[T]) Error(err error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	for _, o := range s.observers {
		o.Error(err)
	}
}

func (s *Subject[T]) Complete() {
	s.mx.RLock()
	defer s.mx.RUnlock()

	for _, o := range s.observers {
		o.Complete()
	}
}

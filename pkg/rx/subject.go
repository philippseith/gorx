package rx

import (
	"fmt"
	"runtime/debug"
	"sync"
)

type Subject[T any] interface {
	Observable[T]
	Observer[T]
}

func NewSubject[T any]() Subject[T] {
	s := &subject[T]{}
	s.Subscribable = s
	return s
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

	return NewSubscription(func() {
		s.mx.Lock()
		defer s.mx.Unlock()

		for i, v := range s.observers {
			if o == v {
				s.observers = append(s.observers[:i], s.observers[i+1:]...)
				return
			}
		}
	})
}

func (s *subject[T]) Next(value T) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	for _, o := range s.observers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					o.Error(fmt.Errorf("panic in %T.Next(%v): %v\n%s", o, value, r, string(debug.Stack())))
				}
			}()
			o.Next(value)
		}()
	}
}

func (s *subject[T]) Error(err error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	// no defer recover with sending to o.Error(), as this would build an endless loop
	for _, o := range s.observers {
		o.Error(err)
	}
}

func (s *subject[T]) Complete() {
	s.mx.RLock()
	defer s.mx.RUnlock()

	for _, o := range s.observers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					o.Error(fmt.Errorf("panic in %T.Complete(): %v\n%s", o, r, string(debug.Stack())))
				}
			}()
			o.Complete()
		}()
	}
}

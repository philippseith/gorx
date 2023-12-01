package rx

import (
	"context"
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
	mxState   sync.RWMutex
	mxEvents  sync.Mutex
}

func (s *subject[T]) Subscribe(o Observer[T]) Subscription {
	s.mxState.Lock()
	defer s.mxState.Unlock()

	s.observers = append(s.observers, o)

	return NewSubscription(func() {
		s.mxState.Lock()
		defer s.mxState.Unlock()

		for i, v := range s.observers {
			if o == v {
				s.observers = append(s.observers[:i], s.observers[i+1:]...)
				return
			}
		}
	})
}

func (s *subject[T]) Next(ctx context.Context, value T) {
	oo := s.getObservers()

	s.mxEvents.Lock()
	defer s.mxEvents.Unlock()

	for _, o := range oo {
		func() {
			defer func() {
				if r := recover(); r != nil {
					o.Error(ctx, fmt.Errorf("panic in %T.Next(%v): %v\n%s", o, value, r, string(debug.Stack())))
				}
			}()
			o.Next(ctx, value)
		}()
	}
}

func (s *subject[T]) Error(ctx context.Context, err error) {
	oo := s.getObservers()

	s.mxEvents.Lock()
	defer s.mxEvents.Unlock()

	// no defer recover with sending to o.Error(), as this would build an endless loop
	for _, o := range oo {
		o.Error(ctx, err)
	}
}

func (s *subject[T]) Complete(ctx context.Context) {
	oo := s.getObservers()

	s.mxEvents.Lock()
	defer s.mxEvents.Unlock()

	for _, o := range oo {
		func() {
			defer func() {
				if r := recover(); r != nil {
					o.Error(ctx, fmt.Errorf("panic in %T.Complete(): %v\n%s", o, r, string(debug.Stack())))
				}
			}()
			o.Complete(ctx)
		}()
	}
}

func (s *subject[T]) getObservers() []Observer[T] {
	s.mxState.RLock()
	defer s.mxState.RUnlock()

	return s.observers
}

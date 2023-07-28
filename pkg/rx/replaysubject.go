package rx

import (
	"sync"
	"time"
)

type ReplaySubject[T any] struct {
	subject[T]
	buffer     []T
	timeStamps []time.Time
	maxBufSize int
	window     time.Duration
	refCount   bool
	complete   bool
	err        error

	mx sync.RWMutex
}

type ReplayOption[T any] func(r *ReplaySubject[T])

func MaxBufferSize[T any](maxBufferSize int) ReplayOption[T] {
	return func(r *ReplaySubject[T]) {
		r.maxBufSize = maxBufferSize
	}
}

func Window[T any](window time.Duration) ReplayOption[T] {
	return func(r *ReplaySubject[T]) {
		r.window = window
	}
}

func RefCount[T any](refCount bool) ReplayOption[T] {
	return func(r *ReplaySubject[T]) {
		r.refCount = refCount
	}
}

func NewReplaySubject[T any](options ...ReplayOption[T]) *ReplaySubject[T] {
	r := &ReplaySubject[T]{}
	for _, o := range options {
		o(r)
	}
	return r
}

func (r *ReplaySubject[T]) Subscribe(o Observer[T]) Subscription {
	s := r.subject.Subscribe(o)
	// TODO Replay (if window is set, remove too old entries)
	return s
}

func (r *ReplaySubject[T]) Next(value T) {
	func() {
		r.mx.RLock()
		defer r.mx.RUnlock()
		// TODO check maxBufferSize (0 means no check)
		r.buffer = append(r.buffer, value)
		r.timeStamps = append(r.timeStamps, time.Now())
	}()

	r.subject.Next(value)
}

func (r *ReplaySubject[T]) Error(err error) {
	func() {
		r.mx.RLock()
		defer r.mx.RUnlock()

		r.err = err
	}()

	r.subject.Error(err)
}

func (r *ReplaySubject[T]) Complete() {
	func() {
		r.mx.RLock()
		defer r.mx.RUnlock()

		r.complete = true
	}()

	r.subject.Complete()
}

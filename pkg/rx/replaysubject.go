package rx

import (
	"sync"
	"time"
)

type ReplaySubject[T any] struct {
	subject[T]
	buffer     []T
	timeStamps []time.Time
	complete   bool
	err        error
	opt        replaySubjectOptions

	mx sync.RWMutex
}

type replaySubjectOptions struct {
	maxBufSize int
	window     time.Duration
	refCount   bool // only used by ShareReplay
}

type ReplayOption func(r *replaySubjectOptions)

func MaxBufferSize(maxBufferSize int) ReplayOption {
	return func(opt *replaySubjectOptions) {
		opt.maxBufSize = maxBufferSize
	}
}

func Window(window time.Duration) ReplayOption {
	return func(opt *replaySubjectOptions) {
		opt.window = window
	}
}

func RefCount(refCount bool) ReplayOption {
	return func(opt *replaySubjectOptions) {
		opt.refCount = refCount
	}
}

func NewReplaySubject[T any](options ...ReplayOption) *ReplaySubject[T] {
	r := &ReplaySubject[T]{}
	for _, o := range options {
		o(&r.opt)
	}
	return r
}

func (r *ReplaySubject[T]) Subscribe(o Observer[T]) Subscription {
	s := r.subject.Subscribe(o)
	if r.opt.window == time.Duration(0) {
		for _, value := range r.buffer {
			r.subject.Next(value)
		}
	} else {
		now := time.Now()
		inWindow := 0
		for i, value := range r.buffer {
			if now.Sub(r.timeStamps[i]) <= r.opt.window {
				r.subject.Next(value)
			} else {
				inWindow = i + 1
			}
		}
		r.buffer = r.buffer[inWindow:]
		r.timeStamps = r.timeStamps[inWindow:]
	}
	return s
}

func (r *ReplaySubject[T]) Next(value T) {
	func() {
		r.mx.RLock()
		defer r.mx.RUnlock()
		if r.opt.maxBufSize == 0 || len(r.buffer) < r.opt.maxBufSize {
			r.buffer = append(r.buffer, value)
			if r.opt.window != time.Duration(0) {
				r.timeStamps = append(r.timeStamps, time.Now())
			}
		} else {
			r.buffer = append(r.buffer[1:], value)
			if r.opt.window != time.Duration(0) {
				r.timeStamps = append(r.timeStamps[1:], time.Now())
			}
		}
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

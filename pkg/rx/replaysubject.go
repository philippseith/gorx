package rx

import (
	"context"
	"sync"
	"time"
)

type ReplaySubject[T any] struct {
	subject[T]
	buffer      []T
	ctxBuffer   []context.Context
	timeStamps  []time.Time
	complete    bool
	ctxComplete context.Context
	err         error
	ctxErr      context.Context
	opt         replaySubjectOptions

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
	r.Subscribable = r
	for _, o := range options {
		o(&r.opt)
	}
	return r
}

func (r *ReplaySubject[T]) Subscribe(o Observer[T]) Subscription {
	s := r.subject.Subscribe(o)
	if r.opt.window == time.Duration(0) {
		func() {
			r.mx.RLock()
			defer r.mx.RUnlock()

			for i, value := range r.buffer {
				r.subject.Next(r.ctxBuffer[i], value)
			}
		}()
	} else {
		now := time.Now()
		inWindow := 0
		func() {
			r.mx.RLock()
			defer r.mx.RUnlock()

			for i, value := range r.buffer {
				if now.Sub(r.timeStamps[i]) <= r.opt.window {
					r.subject.Next(r.ctxBuffer[i], value)
				} else {
					inWindow = i + 1
				}
			}
		}()
		func() {
			r.mx.Lock()
			defer r.mx.Unlock()

			r.buffer = r.buffer[inWindow:]
			r.ctxBuffer = r.ctxBuffer[inWindow:]
			r.timeStamps = r.timeStamps[inWindow:]
		}()
	}
	func() {
		r.mx.RLock()
		defer r.mx.RUnlock()

		if r.err != nil {
			r.subject.Error(r.ctxErr, r.err)
		}
		if r.complete {
			r.subject.Complete(r.ctxComplete)
		}
	}()
	return s
}

func (r *ReplaySubject[T]) Next(ctx context.Context, value T) {
	func() {
		r.mx.Lock()
		defer r.mx.Unlock()

		if r.opt.maxBufSize == 0 || len(r.buffer) < r.opt.maxBufSize {
			r.buffer = append(r.buffer, value)
			r.ctxBuffer = append(r.ctxBuffer, ctx)
			if r.opt.window != time.Duration(0) {
				r.timeStamps = append(r.timeStamps, time.Now())
			}
		} else {
			r.buffer = append(r.buffer[1:], value)
			r.ctxBuffer = append(r.ctxBuffer, ctx)
			if r.opt.window != time.Duration(0) {
				r.timeStamps = append(r.timeStamps[1:], time.Now())
			}
		}
	}()

	r.subject.Next(ctx, value)
}

func (r *ReplaySubject[T]) Error(ctx context.Context, err error) {
	func() {
		r.mx.Lock()
		defer r.mx.Unlock()

		r.err = err
		r.ctxErr = ctx
	}()

	r.subject.Error(ctx, err)
}

func (r *ReplaySubject[T]) Complete(ctx context.Context) {
	func() {
		r.mx.Lock()
		defer r.mx.Unlock()

		r.complete = true
		r.ctxComplete = ctx
	}()

	r.subject.Complete(ctx)
}

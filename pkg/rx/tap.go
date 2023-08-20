package rx

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

func Tap[T any](s Subscribable[T], next func(T) T, err func(error) error, complete func()) Observable[T] {
	t := &tap[T]{
		next:     next,
		err:      err,
		complete: complete,
	}
	t.onSubscribe = func() Subscription { return s.Subscribe(t) }
	return ToObservable[T](t)
}

type tap[T any] struct {
	observer           Observer[T]
	onSubscribe        func() Subscription
	sourceSubscription Subscription
	next               func(T) T
	err                func(error) error
	complete           func()
	mxState            sync.RWMutex
	mxEvents           sync.Mutex
}

func (t *tap[T]) Next(value T) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in %T.Next(%v): %v.\n%s", t, value, r, string(debug.Stack()))
			if o := t.getObserver(); o != nil {
				o.Error(err)
			} else {
				log.Print(err)
			}
		}
	}()

	func() {
		t.mxEvents.Lock()
		defer t.mxEvents.Unlock()

		if t.next != nil {
			value = t.next(value)
		}
		if o := t.getObserver(); o != nil {
			o.Next(value)
		}
	}()

}

func (t *tap[T]) Error(err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in %T.Error(%v): %v.\n%s", t, err, r, string(debug.Stack()))
		}
	}()

	func() {
		t.mxEvents.Lock()
		defer t.mxEvents.Unlock()

		if t.err != nil {
			err = t.err(err)
		}
		if o := t.getObserver(); o != nil {
			o.Error(err)
		}
	}()
}

func (t *tap[T]) Complete() {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in %T.Complete(): %v\n%s", t, r, string(debug.Stack()))
			if o := t.getObserver(); o != nil {
				o.Error(err)
			} else {
				log.Print(err)
			}
		}
	}()

	func() {
		t.mxEvents.Lock()
		defer t.mxEvents.Unlock()

		if t.complete != nil {
			t.complete()
		}
		if o := t.getObserver(); o != nil {
			o.Complete()
		}
	}()
}

func (t *tap[T]) getObserver() Observer[T] {
	t.mxState.RLock()
	defer t.mxState.RUnlock()

	return t.observer
}

func (t *tap[T]) Subscribe(o Observer[T]) Subscription {
	t.mxState.Lock()
	defer t.mxState.Unlock()

	t.observer = o
	t.sourceSubscription = t.onSubscribe()

	return NewSubscription(func() {
		t.mxState.RLock()
		defer t.mxState.RUnlock()

		if t.sourceSubscription != nil {
			t.sourceSubscription.Unsubscribe()
		}
	})
}

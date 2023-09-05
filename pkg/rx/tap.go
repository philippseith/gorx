package rx

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

func Tap[T any](s Subscribable[T], subscribe func(Observer[T]), next func(T) T, err func(error) error, complete, unsubscribe func()) Observable[T] {
	t := &tap[T]{
		subscribe:   subscribe,
		next:        next,
		err:         err,
		complete:    complete,
		unsubscribe: unsubscribe,
	}
	t.onSubscribe = func() Subscription { return s.Subscribe(t) }
	return ToObservable[T](t)
}

type tap[T any] struct {
	observer           Observer[T]
	onSubscribe        func() Subscription
	sourceSubscription Subscription
	subscribe          func(Observer[T])
	next               func(T) T
	err                func(error) error
	complete           func()
	unsubscribe        func()
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
	// If anything is in the operator chain that directly calls Next
	// this would deadlock if we simply lock everything with mxState.
	// So we lock more fine grained
	func() {
		t.mxState.RLock()
		defer t.mxState.RUnlock()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic in %T.Subscribe(%v): %v.\n%s", t, o, r, string(debug.Stack()))
			}
		}()

		if t.subscribe != nil {
			t.subscribe(o)
		}
	}()

	sub := func() func() Subscription {
		t.mxState.Lock()
		defer t.mxState.Unlock()

		t.observer = o
		return t.onSubscribe
	}()()

	func() {
		t.mxState.Lock()
		defer t.mxState.Unlock()

		t.sourceSubscription = sub
	}()

	return NewSubscription(func() {
		t.mxState.RLock()
		defer t.mxState.RUnlock()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic in %T.Unsubscribe(): %v.\n%s", t, r, string(debug.Stack()))
			}
		}()

		if t.unsubscribe != nil {
			t.unsubscribe()
		}
		if t.sourceSubscription != nil {
			t.sourceSubscription.Unsubscribe()
		}
	})
}

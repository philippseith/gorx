package rx

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

// Tap allows to tap into all methods of the Subscriber/Subscribable interface.
// Where values are passed to the subscriber, Tap can alter them.
func Tap[T any](s Subscribable[T],
	subscribe func(Observer[T]),
	next func(context.Context, T) (context.Context, T),
	err func(context.Context, error) (context.Context, error),
	complete func(context.Context) context.Context,
	unsubscribe func()) Observable[T] {
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

// Log allows to log all method invocations of the Subscriber/Subscribable interface.
func Log[T any](s Subscribable[T], id string) Observable[T] {
	return Tap(s, func(o Observer[T]) {
		log.Printf("Subscribe %s: %T", id, o)
	}, func(ctx context.Context, t T) (context.Context, T) {
		log.Printf("Next %s: %+v", id, t)
		return ctx, t
	}, func(ctx context.Context, err error) (context.Context, error) {
		log.Printf("Error %s: %v", id, err)
		return ctx, err
	}, func(ctx context.Context) context.Context {
		log.Printf("Complete %s", id)
		return ctx
	}, func() {
		log.Printf("Unsubscribe %s", id)
	})
}

type tap[T any] struct {
	observer           Observer[T]
	onSubscribe        func() Subscription
	sourceSubscription Subscription
	subscribe          func(Observer[T])
	next               func(context.Context, T) (context.Context, T)
	err                func(context.Context, error) (context.Context, error)
	complete           func(context.Context) context.Context
	unsubscribe        func()
	mxState            sync.RWMutex
	mxEvents           sync.Mutex
}

func (t *tap[T]) Next(ctx context.Context, value T) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in %T.Next(%v): %v.\n%s", t, value, r, string(debug.Stack()))
			t.Error(ctx, err)
		}
	}()

	func() {
		t.mxEvents.Lock()
		defer t.mxEvents.Unlock()

		if t.next != nil {
			ctx, value = t.next(ctx, value)
		}
		if o := t.getObserver(); o != nil {
			o.Next(ctx, value)
		}
	}()

}

func (t *tap[T]) Error(ctx context.Context, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in %T.Error(%v): %v.\n%s", t, err, r, string(debug.Stack()))
		}
	}()

	func() {
		t.mxEvents.Lock()
		defer t.mxEvents.Unlock()

		if t.err != nil {
			ctx, err = t.err(ctx, err)
		}
		if o := t.getObserver(); o != nil {
			o.Error(ctx, err)
		}
	}()
}

func (t *tap[T]) Complete(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in %T.Complete(): %v\n%s", t, r, string(debug.Stack()))
			t.Error(ctx, err)
		}
	}()

	func() {
		t.mxEvents.Lock()
		defer t.mxEvents.Unlock()

		if t.complete != nil {
			ctx = t.complete(ctx)
		}
		if o := t.getObserver(); o != nil {
			o.Complete(ctx)
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

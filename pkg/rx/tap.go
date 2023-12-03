package rx

import (
	"context"
	"fmt"
	"log/slog"
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
		slog.LogAttrs(context.Background(), slog.LevelInfo, "Subscribe",
			slog.Attr{Key: "id", Value: slog.StringValue(id)},
			slog.Attr{Key: "observer", Value: slogTypeValue(o)})
	}, func(ctx context.Context, t T) (context.Context, T) {
		slog.LogAttrs(ctx, slog.LevelInfo, "Next",
			slog.Attr{Key: "id", Value: slog.StringValue(id)},
			slog.Attr{Key: "value", Value: slog.AnyValue(t)})
		return ctx, t
	}, func(ctx context.Context, err error) (context.Context, error) {
		slog.LogAttrs(ctx, slog.LevelInfo, "Error",
			slog.Attr{Key: "id", Value: slog.StringValue(id)},
			slog.Attr{Key: "error", Value: slog.AnyValue(err)})
		return ctx, err
	}, func(ctx context.Context) context.Context {
		slog.LogAttrs(ctx, slog.LevelInfo, "Complete", slog.Attr{Key: "id", Value: slog.StringValue(id)})
		return ctx
	}, func() {
		slog.LogAttrs(context.Background(), slog.LevelInfo, "Unsubscribe", slog.Attr{Key: "id", Value: slog.StringValue(id)})
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
			slog.LogAttrs(ctx, slog.LevelError, "panic in Next",
				slog.Attr{Key: "observer", Value: slogTypeValue(t.observer)},
				slog.Attr{Key: "value", Value: slog.AnyValue(value)},
				slog.Group("panic",
					slog.Attr{Key: "error", Value: slog.AnyValue(r)},
					slog.Attr{Key: "stack", Value: slog.StringValue(string(debug.Stack()))}))

			err := fmt.Errorf("panic in %T.Next(%v): %v.\n%s", t.observer, value, r, string(debug.Stack()))
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
			slog.LogAttrs(ctx, slog.LevelError, "panic in Error",
				slog.Attr{Key: "observer", Value: slogTypeValue(t.observer)},
				slog.Attr{Key: "error", Value: slog.AnyValue(err)},
				slog.Group("panic",
					slog.Attr{Key: "error", Value: slog.AnyValue(r)},
					slog.Attr{Key: "stack", Value: slog.StringValue(string(debug.Stack()))}))
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
			slog.LogAttrs(ctx, slog.LevelError, "panic in Complete",
				slog.Attr{Key: "observer", Value: slogTypeValue(t.observer)},
				slog.Group("panic",
					slog.Attr{Key: "error", Value: slog.AnyValue(r)},
					slog.Attr{Key: "stack", Value: slog.StringValue(string(debug.Stack()))}))

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
				slog.LogAttrs(context.Background(), slog.LevelError, "panic in Subscribe",
					slog.Attr{Key: "observer", Value: slogTypeValue(o)},
					slog.Group("panic",
						slog.Attr{Key: "error", Value: slog.AnyValue(r)},
						slog.Attr{Key: "stack", Value: slog.StringValue(string(debug.Stack()))}))

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
				slog.LogAttrs(context.Background(), slog.LevelError, "panic in Unsubscribe",
					slog.Attr{Key: "observer", Value: slogTypeValue(o)},
					slog.Group("panic",
						slog.Attr{Key: "error", Value: slog.AnyValue(r)},
						slog.Attr{Key: "stack", Value: slog.StringValue(string(debug.Stack()))}))
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

func slogTypeValue(a any) slog.Value {
	return slog.StringValue(fmt.Sprintf("%T", a))
}

package rx

import (
	"context"
	"sync"
	"time"
)

var contextKeyDebounceTrigger = ContextKey("Debounce Trigger")

// DebounceTime returns an Observable that delays the emissions of the source
// Subscribable by the specified duration and may drop some values if they occur
// too frequently.
func DebounceTime[T any](s Subscribable[T], duration time.Duration) Observable[T] {
	return Debounce[T, time.Time](s, NewTicker(context.Background(), 0, duration))
}

// Debounce emits a notification from the source Observable only after a
// particular time span has passed without another source emission
func Debounce[T any, U any](s Subscribable[T], trigger Subscribable[U]) Observable[T] {
	d := &debounce[T]{
		Operator: Operator[T, T]{t2u: func(t T) T { return t }},
	}

	d.prepareSubscribe(func() Subscription {
		triggerSub := trigger.Subscribe(OnNextWithContext(func(_ context.Context, u U) {
			if func() bool {
				d.mx.RLock()
				defer d.mx.RUnlock()

				return d.hasLast
			}() {
				func() {
					d.mx.Lock()
					defer d.mx.Unlock()

					d.hasLast = false
				}()
				d.Operator.Next(context.WithValue(d.ctx, contextKeyDebounceTrigger, u), d.last) // nolint:contextcheck
			}

		}))
		return s.Subscribe(d).AddSubscription(triggerSub)
	})
	ds := ToObservable[T](d)
	// Unsubscribe trigger
	return ds
}

type debounce[T any] struct {
	Operator[T, T]
	hasLast bool
	ctx     context.Context
	last    T
	mx      sync.RWMutex
}

func (d *debounce[T]) Next(ctx context.Context, next T) {
	d.mx.Lock()
	defer d.mx.Unlock()

	d.ctx = ctx
	d.last = next
	d.hasLast = true
}

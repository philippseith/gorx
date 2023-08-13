package rx

import (
	"sync"
	"time"
)

// DebounceTime returns an Observable that delays the emissions of the source
// Subscribable by the specified duration and may drop some values if they occur
// too frequently.
// TODO Behavior is not correct. It should emit after time, not when the first goes in after that time
func DebounceTime[T any](s Subscribable[T], duration time.Duration) Observable[T] {
	d := &debounceTime[T]{
		Operator: Operator[T, T]{t2u: func(t T) T { return t }},
		duration: duration,
		last:     time.Unix(0, 0),
	}
	d.SubscribeToSource(d, s)
	return ToObservable[T](d)
}

type debounceTime[T any] struct {
	Operator[T, T]
	duration time.Duration
	last     time.Time
}

func (d *debounceTime[T]) Next(value T) {
	if time.Since(d.last) > d.duration {
		d.last = time.Now()
		d.Operator.Next(value)
	}
}

func Debounce[T any, U any](s Subscribable[T], trigger Subscribable[U]) Observable[T] {
	d := &debounce[T]{
		Operator: Operator[T, T]{t2u: func(t T) T { return t }},
	}

	triggerSub := trigger.Subscribe(OnNext(func(U) {
		d.mx.RLock()
		defer d.mx.RUnlock()

		if d.hasLast {
			d.hasLast = false
			d.Next(d.last)
		}
	}))
	d.SubscribeToSource(d, s)
	ds := ToObservable[T](d)
	// Unsubscribe trigger
	ds.AddTearDownLogic(triggerSub.Unsubscribe)
	return ds
}

type debounce[T any] struct {
	Operator[T, T]
	hasLast bool
	last    T
	mx      sync.RWMutex
}

func (d *debounce[T]) Next(next T) {
	d.mx.Lock()
	defer d.mx.Unlock()

	d.last = next
	d.hasLast = true
}

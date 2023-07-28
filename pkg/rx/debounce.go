package rx

import "time"

// Debounce returns an Observable that delays the emissions of the source
// Subscribable by the specified duration and may drop some values if they occur
// too frequently.
func Debounce[T any](s Subscribable[T], duration time.Duration) Observable[T] {
	d := &debounce[T]{
		observableObserver: observableObserver[T, T]{
			t2u: func(t T) T {
				return t
			},
		},
		duration: duration,
		last:     time.Unix(0, 0),
	}
	d.Subscribable = d
	d.sourceSub = func() { s.Subscribe(d) }
	return d
}

type debounce[T any] struct {
	observableObserver[T, T]
	duration time.Duration
	last     time.Time
}

func (d *debounce[T]) Next(value T) {
	if time.Since(d.last) > d.duration {
		d.last = time.Now()
		d.observableObserver.Next(value)
	}
}

package rx

import "time"

func Debounce[T any](o Observable[T], duration time.Duration) Observable[T] {
	d := &debounce[T]{
		observableObserver: observableObserver[T, T]{
			t2u: func(t T) T {
				return t
			},
		},
		duration: duration,
		last:     time.Unix(0, 0),
	}
	d.sourceSub = func() { o.Subscribe(d) }
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

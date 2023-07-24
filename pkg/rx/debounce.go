package rx

import "time"

func Debounce[T any](o Observable[T], duration time.Duration) Observable[T] {
	d := &debounce[T]{
		duration: duration,
		last:     time.Unix(0, 0),
	}
	o.Subscribe(d)
	return d
}

type debounce[T any] struct {
	observableObserver[T]
	duration time.Duration
	last     time.Time
}

func (d *debounce[T]) Next(value T) {
	if time.Since(d.last) > d.duration {
		d.last = time.Now()
		d.observableObserver.Next(value)
	}
}

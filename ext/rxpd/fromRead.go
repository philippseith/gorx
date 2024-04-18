package rxpd

import (
	"context"
	"time"

	"github.com/philippseith/gorx/pkg/rx"
)

func FromRead[T any](ctx context.Context, readOption PropertyOption[T], options ...PropertyOption[T]) Property[T] {

	ch := make(chan T)
	ivCh := make(chan time.Duration, 1)

	p := &property[T]{Subscribable: rx.FromChan(ch)}
	for _, option := range options {
		option(&p.propertyOption)
	}

	readOption(&p.propertyOption)

	WithSetInterval[T](func(d time.Duration) {
		ivCh <- d
	})(&p.propertyOption)

	go func() {
		ticker := time.NewTicker(time.Second)
		ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				close(ch)
				close(ivCh)
				return
			case <-ticker.C:
				result := <-p.propertyOption.read()
				if result.Err != nil {
					continue
				}
				ch <- result.Ok
			case interval := <-ivCh:
				ticker.Stop()
				if interval > 0 {
					ticker = time.NewTicker(interval)
				}
			}
		}
	}()

	return p
}

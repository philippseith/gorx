package rxpd

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/philippseith/gorx/pkg/rx"
)

// FromRead creates a Property from a WithRead option. The reactive stream is created by a goroutine
// which calls the read function in the interval set by SetInterval.
func FromRead[T any](ctx context.Context, read func() rx.ResultChan[T], options ...PropertyOption[T]) Property[T] {

	ch := make(chan T)
	ivCh := make(chan time.Duration, 1)

	p := &property[T]{Subscribable: rx.FromChan(ch)}
	for _, option := range options {
		option(&p.propertyOption)
	}

	WithSetInterval[T](func(d time.Duration) {
		if ctx.Err() != nil {
			return
		}
		ivCh <- d
	})(&p.propertyOption)

	p.read = read

	go func() {
		ticker := time.NewTicker(time.Second)
		ticker.Stop()

		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				close(ch)
				close(ivCh)
				return
			case <-ticker.C:
				result := func() (result rx.Result[T]) {
					defer func() {
						if r := recover(); r != nil {
							result = rx.Result[T]{Err: fmt.Errorf("panic in FromRead using %T.Read(): %v.\n%s", p, r, string(debug.Stack()))}
						}
					}()
					result = <-p.propertyOption.read()
					return result
				}()
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

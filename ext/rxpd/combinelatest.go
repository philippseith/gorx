package rxpd

import (
	"fmt"
	"sync"

	"github.com/philippseith/gorx/pkg/rx"
)

func CombineLatest[T any](combine func(next ...any) T, sources ...Subscribable[any]) Property[T] {
	rxSources := make([]rx.Subscribable[any], 0, len(sources))
	for _, source := range sources {
		rxSources = append(rxSources, source)
	}
	q := &property[T]{
		Subscribable: rx.CombineLatest[T](combine, rxSources...),
		propertyOption: propertyOption[T]{
			setInterval: nil,
			read: func() <-chan rx.Result[T] {
				ch := make(chan rx.Result[T])
				go func() {
					var wg sync.WaitGroup
					wg.Add(len(sources))
					results := make([]rx.Result[any], len(sources))
					for i, s := range sources {
						idx := i
						source := s
						go func() {
							results[idx] = <-source.Read()
							wg.Done()
						}()
					}
					wg.Wait()
					for i, result := range results {
						if result.Err != nil {
							ch <- rx.Result[T]{Err: fmt.Errorf("CombineLatest.Read source %v: %w", i, result.Err)}
							close(ch)
							return
						}
					}
					next := make([]any, 0, len(sources))
					for _, result := range results {
						next = append(next, result)
					}
					ch <- rx.Result[T]{Ok: combine(next...)}
					close(ch)
				}()
				return ch
			},
		},
	}
	return q
}

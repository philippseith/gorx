package rx

import "context"

// ToSlice converts an observable sequence into a slice. ToSlice blocks the
// current goroutine until the observable completes or errors or the context is
// canceled. If the observable ended with error, the slice will be empty.
func ToSlice[T any](ctx context.Context, o Observable[T]) []T {
	var s []T
	ch := make(chan struct{})
	o.Subscribe(NewObserver(func(value T) {
		s = append(s, value)
	}, func(error) {
		s = []T{}
		close(ch)
	}, func() {
		close(ch)
	}))
	select {
	case <-ch:
	case <-ctx.Done():
	}
	return s
}

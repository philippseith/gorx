package rx

import "context"

func Go[T any](ctx context.Context, f func() (T, error)) <-chan Result[T] {
	ch := make(chan Result[T], 1)
	go func() {
		defer close(ch)

		select {
		case <-ctx.Done():
			ch <- Err[T](ctx.Err())
		case r := <-func() <-chan Result[T] {
			chf := make(chan Result[T])
			go func() {
				r := Result[T]{}
				r.Ok, r.Err = f()
				chf <- r
				close(chf)
			}()
			return chf
		}():
			ch <- r
		}
	}()
	return ch
}

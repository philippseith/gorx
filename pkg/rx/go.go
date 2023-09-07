package rx

import "context"

func Go[T any](ctx context.Context, f func() (T, error)) (T, error) {
	select {
	case <-ctx.Done():
		r := Result[T]{}
		return r.Ok, ctx.Err()
	case r := <-func() <-chan Result[T] {
		ch := make(chan Result[T])
		go func() {
			r := Result[T]{}
			r.Ok, r.Err = f()
			ch <- r
			close(ch)
		}()
		return ch
	}():
		return r.Ok, r.Err
	}
}

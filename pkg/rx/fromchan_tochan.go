package rx

import "context"

// FromChan creates an Observable[T] from a chan T. The channel will be read at
// once. All values sent before the Observable is subscribed to, will be ignored.
// FromChan should only be subscribed once simultaneously.
func FromChan[T any](ctx context.Context, ch <-chan T) Observable[T] {
	fc := &Operator[T, T]{t2u: func(t T) T { return t }}
	fc.prepareSubscribe(func() Subscription {
		go func() {
			for t := range ch {
				fc.Next(ctx, t)
			}
			fc.Complete(ctx)
		}()
		return NewSubscription(func() {})
	})
	// TODO prevent simultaneous subscriptions
	return ToObservable[T](fc)
}

// ResultChan is a sending channel for Result[T]
type ResultChan[T any] <-chan Result[T]

// ToObservable creates an Observable[T] from a ResultChan[T]. The channel will be read at
// once. All values sent before the Observable is subscribed to, will be ignored.
// FromResultChan should only be subscribed once simultaneously.
func (ch ResultChan[T]) ToObservable(ctx context.Context) Observable[T] {
	fc := &Operator[T, T]{t2u: func(t T) T { return t }}
	fc.prepareSubscribe(func() Subscription {
		go func() {
			for item := range ch {
				if item.Err != nil {
					fc.Error(ctx, item.Err)
				} else {
					fc.Next(ctx, item.Ok)
				}
			}
			fc.Complete(ctx)
		}()
		return NewSubscription(func() {})
	})
	return ToObservable[T](fc)
}

// OnNext adds a Next handler to a ResultChan
func (ch ResultChan[T]) OnNext(next func(T)) Subscription {
	return ch.ToObservable(context.Background()).Subscribe(OnNext[T](next))
}

func (ch ResultChan[T]) OnNextWithContext(ctx context.Context, next func(context.Context, T)) Subscription {
	return ch.ToObservable(ctx).Subscribe(OnNextWithContext[T](next))
}

// ToChan pushes the values from a Subscribable into a channel. It returns a
// channel and the Subscription to the Observable. the channel type Item[T]
// contains values and a possible error. Note that ToChan will block immediately
// with cold observables. You need to wrap the cold observable with
// ToConnectable and set up a receiving goroutine for the channel before you
// call Connectable.Connect
func ToChan[T any](s Subscribable[T]) (ResultChan[T], Subscription) {
	tc := &toChan[T]{
		ch: make(chan Result[T], 1),
	}
	return tc.ch, s.Subscribe(tc)
}

type toChan[T any] struct {
	ch chan Result[T]
}

func (tc *toChan[T]) Next(_ context.Context, value T) {
	tc.ch <- Result[T]{Ok: value}
}

func (tc *toChan[T]) Error(_ context.Context, err error) {
	tc.ch <- Result[T]{Err: err}
	close(tc.ch)
}

func (tc *toChan[T]) Complete(context.Context) {
	close(tc.ch)
}

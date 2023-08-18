package rx

// FromChan creates an Observable[T] from a chan T. The channel will be read at
// once. All values sent before the Observable is subscribed to, will be ignored.
// FromChan should only be subscribed once simultaneously.
func FromChan[T any](ch <-chan T) Observable[T] {
	fc := &Operator[T, T]{t2u: func(t T) T { return t }}
	fc.prepareSubscribe(func() Subscription {
		go func() {
			for t := range ch {
				fc.Next(t)
			}
			fc.Complete()
		}()
		return NewSubscription(func() {})
	})
	// TODO prevent simultaneous subscriptions
	return ToObservable[T](fc)
}

// FromItemChan creates an Observable[T] from a chan Item[T]. The channel will be read at
// once. All values sent before the Observable is subscribed to, will be ignored.
// FromItemChan should only be subscribed once simultaneously.
func FromItemChan[T any](ch <-chan Item[T]) Observable[T] {
	fc := &Operator[T, T]{t2u: func(t T) T { return t }}
	fc.prepareSubscribe(func() Subscription {
		go func() {
			for item := range ch {
				if item.E != nil {
					fc.Error(item.E)
				} else {
					fc.Next(item.V)
				}
			}
			fc.Complete()
		}()
		return NewSubscription(func() {})
	})
	return ToObservable[T](fc)
}

// ToChan pushes the values from a Subscribable into a channel. It returns a
// channel and the Subscription to the Observable. the channel type Item[T]
// contains values and a possible error. Note that ToChan will block immediately
// with cold observables. You need to wrap the cold observable with
// ToConnectable and set up a receiving goroutine for the channel before you
// call Connectable.Connect
func ToChan[T any](s Subscribable[T]) (<-chan Item[T], Subscription) {
	tc := &toChan[T]{
		ch: make(chan Item[T], 1),
	}
	return tc.ch, s.Subscribe(tc)
}

type Item[T any] struct {
	V T
	E error
}

type toChan[T any] struct {
	ch chan Item[T]
}

func (tc *toChan[T]) Next(value T) {
	tc.ch <- Item[T]{V: value}
}

func (tc *toChan[T]) Error(err error) {
	tc.ch <- Item[T]{E: err}
	close(tc.ch)
}

func (tc *toChan[T]) Complete() {
	close(tc.ch)
}

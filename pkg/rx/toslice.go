package rx

// ToSlice converts a Subscribable sequence into a slice. If the Subscribable
// ends with error, the channel will be closed without result. If the
// Subscribable produces its values in the same goroutine, waiting on the channel
// will deadlock.
func ToSlice[T any](s Subscribable[T]) <-chan []T {
	var sl []T
	ch := make(chan []T, 1)
	s.Subscribe(NewObserver(func(value T) {
		sl = append(sl, value)
	}, func(error) {
		close(ch)
	}, func() {
		ch <- sl
		close(ch)
	}))
	return ch
}

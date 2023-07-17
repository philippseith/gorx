package rx

type Subscriber[T any] struct {
}

func (s *Subscriber[T]) Next(value T) {
}

func (s *Subscriber[T]) Error(err error) {
}

func (s *Subscriber[T]) Complete() {
}

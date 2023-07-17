package rx

type DefaultObserver[T any] struct {
}

func (d DefaultObserver[T]) Next(T) {
}

func (d DefaultObserver[T]) Error(error) {
}

func (d DefaultObserver[T]) Complete() {
}

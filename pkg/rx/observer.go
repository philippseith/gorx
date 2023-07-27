package rx

type Observer[T any] interface {
	Next(value T)
	Error(err error)
	Complete()
}

func NewObserver[T any](next func(T), err func(error), complete func()) Observer[T] {
	return &observer[T]{
		next:     next,
		err:      err,
		complete: complete,
	}
}

func OnNext[T any](next func(T)) Observer[T] {
	return &observer[T]{next: next}
}

type observer[T any] struct {
	next     func(T)
	err      func(error)
	complete func()
}

func (o *observer[T]) Next(value T) {
	if o.next != nil {
		o.next(value)
	}
}

func (o *observer[T]) Error(err error) {
	if o.err != nil {
		o.err(err)
	}
}

func (o *observer[T]) Complete() {
	if o.complete != nil {
		o.complete()
	}
}

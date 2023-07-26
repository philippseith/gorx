package rx

func From[T any](items ...T) Observable[T] {
	return &from[T]{
		items: items,
	}
}

type from[T any] struct {
	observable[T]
	items []T
}

func (f *from[T]) Subscribe(o Observer[T]) Subscription {
	for _, item := range f.items {
		o.Next(item)
	}
	o.Complete()
	return &subscription{u: func() {}}
}

func Just[T any](value T) Observable[T] {
	return &just[T]{value: value}
}

type just[T any] struct {
	observable[T]
	value T
}

func (j *just[T]) Subscribe(o Observer[T]) Subscription {
	o.Next(j.value)
	o.Complete()
	return &subscription{u: func() {}}
}

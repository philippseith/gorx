package rx

func Of[T any](items ...T) Subscribable[T] {
	return &of[T]{
		items: items,
	}
}

type of[T any] struct {
	items []T
}

func (of *of[T]) Subscribe(o Observer[T]) Subscription {
	for _, item := range of.items {
		o.Next(item)
		o.Complete()
	}
	return &subscription{u: func() {}}
}

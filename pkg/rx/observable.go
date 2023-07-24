package rx

type Observable[T any] interface {
	Subscribe(o Observer[T]) Subscription
}

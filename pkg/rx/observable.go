package rx

type Observable[T any] interface {
	Subscribe(o Observer[T]) Subscription
}

// ToAny converts an Observable[T] to an Observable[any]
func ToAny[T any](o Observable[T]) Observable[any] {
	oa := &observableObserver[T, any]{
		t2u: func(t T) any {
			return t
		},
	}
	oa.sourceSub = func() {
		o.Subscribe(oa)
	}
	return oa
}

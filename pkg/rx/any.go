package rx

// ToAny converts an Subscribable[T] to an Observable[any]
func ToAny[T any](s Subscribable[T]) Observable[any] {
	oa := &observableObserver[T, any]{
		t2u: func(t T) any {
			return t
		},
	}
	oa.sourceSub = func() Subscription {
		return s.Subscribe(oa)
	}
	return oa
}

package rx

// ToAny converts an Observable[T] to an Observable[any]
func ToAny[T any](s Subscribable[T]) Observable[any] {
	oa := &observableObserver[T, any]{
		t2u: func(t T) any {
			return t
		},
	}
	oa.sourceSub = func() {
		s.Subscribe(oa)
	}
	return oa
}

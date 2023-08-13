package rx

// ToAny converts a Subscribable[T] to an Observable[any]
func ToAny[T any](s Subscribable[T]) Observable[any] {
	oa := &Operator[T, any]{t2u: func(t T) any { return t }}
	oa.SubscribeToSource(oa, s)
	return ToObservable[any](oa)
}

package rx

// ToAny converts a Subscribable[T] to an Observable[any]
func ToAny[T any](s Subscribable[T]) Observable[any] {
	oa := &Operator[T, any]{t2u: func(t T) any { return t }}
	oa.prepareSubscribe(func() Subscription { return s.Subscribe(oa) })
	return ToObservable[any](oa)
}

func ToTyped[T any](s Subscribable[any]) Observable[T] {
	oa := &Operator[any, T]{t2u: func(a any) T { return a.(T) }}
	oa.prepareSubscribe(func() Subscription { return s.Subscribe(oa) })
	return ToObservable[T](oa)
}

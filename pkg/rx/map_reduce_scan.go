package rx

import "context"

func Map[T any, U any](s Subscribable[T], mapper func(T) U) Observable[U] {
	m := &Operator[T, U]{t2u: mapper}
	m.prepareSubscribe(func() Subscription { return s.Subscribe(m) })
	return ToObservable[U](m)
}

func Reduce[T any, U any](s Subscribable[T], acc func(U, T) U, seed U) Observable[U] {
	result := &seed
	r := &Operator[U, U]{t2u: func(u U) U { return u }}
	r.prepareSubscribe(func() Subscription {
		return s.Subscribe(NewObserverWithContext[T](func(ctx context.Context, value T) {
			*result = acc(*result, value)
		}, r.Error, func(ctx context.Context) {
			r.Next(ctx, *result)
			r.Complete(ctx)
		}))
	})
	return ToObservable[U](r)
}

func Scan[T any, U any](s Subscribable[T], acc func(U, T) U, seed U) Observable[U] {
	sc := &struct {
		Operator[T, U]
		acc U
	}{acc: seed}
	sc.Operator = Operator[T, U]{t2u: func(t T) U {
		sc.acc = acc(sc.acc, t)
		return sc.acc
	}}
	sc.prepareSubscribe(func() Subscription { return s.Subscribe(sc) })
	return ToObservable[U](sc)
}

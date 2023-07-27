package rx

func Map[T any, U any](s Subscribable[T], mapper func(T) U) Observable[U] {
	oou := &observableObserver[U, U]{
		t2u: func(t U) U {
			return t
		},
	}
	oou.sourceSub = func() {
		s.Subscribe(NewObserver[T](func(value T) {
			oou.Next(mapper(value))
		}, oou.Error, oou.Complete))
	}
	return oou
}

func Reduce[T any, U any](s Subscribable[T], acc func(U, T) U, seed U) Observable[U] {
	result := []U{seed}
	oou := &observableObserver[U, U]{
		t2u: func(t U) U {
			return t
		},
	}
	oou.sourceSub = func() {
		s.Subscribe(NewObserver[T](func(value T) {
			result[0] = acc(result[0], value)
		}, oou.Error, func() {
			oou.Next(result[0])
			oou.Complete()
		}))
	}
	return oou
}

func Scan[T any, U any](s Subscribable[T], acc func(U, T) U, seed U) Observable[U] {
	sc := &struct {
		observableObserver[T, U]
		acc U
	}{
		observableObserver: observableObserver[T, U]{},
		acc:                seed,
	}
	sc.t2u = func(t T) U {
		sc.acc = acc(sc.acc, t)
		return sc.acc
	}
	sc.sourceSub = func() {
		s.Subscribe(sc)
	}
	return sc
}

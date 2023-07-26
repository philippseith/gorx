package rx

func Map[T any, U any](o Observable[T], mapper func(T) U) Observable[U] {
	oou := &observableObserver[U, U]{
		t2u: func(t U) U {
			return t
		},
	}
	oou.sourceSub = func() {
		o.Subscribe(NewObserver[T](func(value T) {
			oou.Next(mapper(value))
		}, oou.Error, oou.Complete))
	}
	return oou
}

func Reduce[T any, U any](o Observable[T], acc func(U, T) U, seed U) Observable[U] {
	result := []U{seed}
	oou := &observableObserver[U, U]{
		t2u: func(t U) U {
			return t
		},
	}
	oou.sourceSub = func() {
		o.Subscribe(NewObserver[T](func(value T) {
			result[0] = acc(result[0], value)
		}, oou.Error, func() {
			oou.Next(result[0])
			oou.Complete()
		}))
	}
	return oou
}

func Scan[T any, U any](o Observable[T], acc func(U, T) U, seed U) Observable[U] {
	s := &struct {
		observableObserver[T, U]
		acc U
	}{
		observableObserver: observableObserver[T, U]{},
		acc:                seed,
	}
	s.t2u = func(t T) U {
		s.acc = acc(s.acc, t)
		return s.acc
	}
	s.sourceSub = func() {
		o.Subscribe(s)
	}
	return s
}

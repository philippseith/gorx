package rx

func Map[T any, U any](s Subscribable[T], mapper func(T) U) Subscribable[U] {
	su := NewSubject[U]()
	s.Subscribe(NewObserver[T](func(value T) {
		su.Next(mapper(value))
	}, su.Error, su.Complete))
	return su
}

func Reduce[T any, U any](s Subscribable[T], acc func(U, T) U, seed U) Subscribable[U] {
	result := []U{seed}
	su := NewSubject[U]()
	s.Subscribe(NewObserver[T](func(value T) {
		result[0] = acc(result[0], value)
	}, su.Error, func() {
		su.Next(result[0])
		su.Complete()
	}))
	return su
}

func Scan[T any, U any](s Subscribable[T], acc func(U, T) U, seed U) Subscribable[U] {
	result := []U{seed}
	su := NewSubject[U]()
	s.Subscribe(NewObserver[T](func(value T) {
		result[0] = acc(result[0], value)
		su.Next(result[0])
	}, su.Error, su.Complete))
	return su
}

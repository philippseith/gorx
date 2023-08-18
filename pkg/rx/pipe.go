package rx

import "time"

func DebounceX[T any, U any](trigger Subscribable[U]) func(Subscribable[T]) Observable[T] {
	return func(s Subscribable[T]) Observable[T] {
		return Debounce(s, trigger)
	}
}

func DebounceTimeX[T any](duration time.Duration) func(Subscribable[T]) Observable[T] {
	return func(s Subscribable[T]) Observable[T] {
		return Debounce[T, time.Time](s, NewTicker(0, duration))
	}
}

func MapX[T any, U any](mapper func(T) U) func(Subscribable[T]) Observable[U] {
	return func(s Subscribable[T]) Observable[U] {
		return Map(s, mapper)
	}
}

func test1() {
	s := NewSubject[int]()
	t := NewTicker(0, time.Millisecond)
	DebounceX[float32, time.Time](t)(MapX[int, float32](func(i int) float32 {
		return 0
	})(s))
}

func Pipe2[T1 any, T2 any, T3 any](f1 func(Subscribable[T1]) Observable[T2], f2 func(Subscribable[T2]) Observable[T3]) func(Subscribable[T1]) Subscribable[T3] {
	return func(a Subscribable[T1]) Subscribable[T3] {
		return f2(f1(a))
	}
}

func test2() {
	Pipe2(MapX(func(int) float32 { return 0 }), DebounceTimeX[float32](time.Millisecond))(NewSubject[int]()).Subscribe(OnNext[float32](nil))
}

type P1 interface {
	Any(subscribable Subscribable[any]) Observable[any]
}

type P2[T1 any, T2 any] interface {
	P1
	Typed(Subscribable[T1]) Observable[T2]
}

type p2[T1 any, T2 any] struct {
	typed func(Subscribable[T1]) Observable[T2]
}

func (p *p2[T1, T2]) Any(s Subscribable[any]) Observable[any] {
	return p.Typed(ToTyped[T1](s)).ToAny()
}

func (p *p2[T1, T2]) Typed(s Subscribable[T1]) Observable[T2] {
	return p.typed(s)
}

func Pipe3(ff ...P1) func(Subscribable[any]) Subscribable[any] {
	return func(s Subscribable[any]) Subscribable[any] {
		for _, f := range ff {
			s = f.Any(s)
		}
		return s
	}
}

func MapY[T any, U any](mapper func(T) U) P2[T, U] {
	return nil
}

func test3() {
	Pipe3(MapY[int, float32](nil), MapY[uint, int](nil))
}

package rxpd

import (
	"fmt"

	"github.com/philippseith/gorx/pkg/rx"
)

func Map[T any, U any](s Subscribable[T], mapToU func(T) U, mapToT func(U) (T, error)) Property[U] {
	return &property[U]{
		Subscribable: rx.Map[T, U](s, mapToU),
		propertyOption: propertyOption[U]{
			setInterval: s.SetInterval,
			read: func() <-chan rx.Result[U] {
				ch := make(chan rx.Result[U], 1)
				go func() {
					result := <-s.Read()
					if result.Err != nil {
						ch <- rx.Result[U]{Err: result.Err}
					} else {
						ch <- rx.Result[U]{Ok: mapToU(result.Ok)}
					}
					close(ch)
				}()
				return ch
			},
			write: func(u U) <-chan error {
				ch := make(chan error, 1)
				go func() {
					if t, err := mapToT(u); err != nil {
						ch <- err
					} else {
						err := <-s.Write(t)
						ch <- err
					}
					close(ch)
				}()
				return ch
			},
		},
	}
}

func ToAny[T any](s Subscribable[T]) Property[any] {
	return Map[T, any](s, func(t T) any { return t }, func(a any) (T, error) {
		if valueT, ok := a.(T); ok {
			return valueT, nil
		}
		return *new(T), fmt.Errorf("can not convert %v to %T", a, *new(T))
	})
}

package rx

import (
	"fmt"
	"log"
	"runtime/debug"
)

type Observer[T any] interface {
	Next(value T)
	Error(err error)
	Complete()
}

func NewObserver[T any](next func(T), err func(error), complete func()) Observer[T] {
	return &observer[T]{
		next:     next,
		err:      err,
		complete: complete,
	}
}

func OnNext[T any](next func(T)) Observer[T] {
	return &observer[T]{next: next, err: func(err error) { log.Print(err) }}
}

type observer[T any] struct {
	next     func(T)
	err      func(error)
	complete func()
}

func (o *observer[T]) Next(value T) {
	defer func() {
		if r := recover(); r != nil {
			o.Error(fmt.Errorf("panic in %T.Next(%v): %v.\n%s", o, value, r, string(debug.Stack())))
		}
	}()

	if o.next != nil {
		o.next(value)
	}
}

func (o *observer[T]) Error(err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in %T.Error(%v): %v.\n%s", o, err, r, string(debug.Stack()))
		}
	}()

	if o.err != nil {
		o.err(err)
	}
}

func (o *observer[T]) Complete() {
	defer func() {
		if r := recover(); r != nil {
			o.Error(fmt.Errorf("panic in %T.Complete(): %v\n%s", o, r, string(debug.Stack())))
		}
	}()

	if o.complete != nil {
		o.complete()
	}
}

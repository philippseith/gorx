package rxpd

import (
	"log"
	"time"

	"github.com/philippseith/gorx/pkg/rx"
)

type PropertyOption[T any] func(*propertyOption[T])

type propertyOption[T any] struct {
	setInterval   func(time.Duration)
	read          func() rx.ResultChan[T]
	write         func(T) <-chan error
	readCallsNext bool
	readTimeout   time.Duration
	onReadTimeout func() rx.Result[T]
}

func WithSetInterval[T any](setInterval func(time.Duration)) func(*propertyOption[T]) {
	return func(po *propertyOption[T]) {
		po.setInterval = setInterval
	}
}

func WithRead[T any](read func() rx.ResultChan[T]) func(*propertyOption[T]) {
	return func(po *propertyOption[T]) {
		po.read = read
	}
}

func WithTapRead[T any](tapRead func(func() rx.ResultChan[T]) rx.ResultChan[T]) func(*propertyOption[T]) {
	return func(po *propertyOption[T]) {
		poRead := po.read
		po.read = func() rx.ResultChan[T] {
			return tapRead(poRead)
		}
	}
}

func WithWrite[T any](write func(T) <-chan error) func(*propertyOption[T]) {
	return func(po *propertyOption[T]) {
		po.write = write
	}
}

func WithTapWrite[T any](tapWrite func(T, func(T) <-chan error) <-chan error) func(*propertyOption[T]) {
	return func(po *propertyOption[T]) {
		poWrite := po.write
		po.write = func(t T) <-chan error {
			return tapWrite(t, poWrite)
		}
	}
}

func WithLogReadWrite[T any](id string) func(*propertyOption[T]) {
	return func(po *propertyOption[T]) {
		poRead := po.read
		po.read = func() rx.ResultChan[T] {
			log.Printf("Before Read %s", id)
			ch := make(chan rx.Result[T])
			go func() {
				result := <-poRead()
				log.Printf("After Read %s: %v", id, result)
				ch <- result
				close(ch)
			}()
			return ch
		}
		po.write = func(t T) <-chan error {
			poWrite := po.write
			log.Printf("Before Write %s: %v", id, t)
			ch := make(chan error)
			go func() {
				err := <-poWrite(t)
				log.Printf("After Write %s: %v", id, err)
				ch <- err
				close(ch)
			}()
			return ch
		}
	}
}

// WithReadCallsNext configures that the result from the Read rx.ResultChan is passed to the Next function.
// When used together with WithReadTimeout and the timeout elapsed, first the timeout value is passed to the Next function.
// When Read finally returns, this value is also passed to the Next function.
func WithReadCallsNext[T any]() func(*propertyOption[T]) {
	return func(po *propertyOption[T]) {
		po.readCallsNext = true
	}
}

// WithReadTimeout configures a timeout for the Read function. If the timeout is reached,
// the onTimeout function is called and its result is returned.
func WithReadTimeout[T any](timeout time.Duration, onTimeout func() rx.Result[T]) func(*propertyOption[T]) {
	return func(po *propertyOption[T]) {
		po.readTimeout = timeout
		po.onReadTimeout = onTimeout
	}
}

// With allows changing PropertyOption after creation.
func (p *property[T]) With(options ...PropertyOption[T]) Property[T] {
	for _, option := range options {
		option(&p.propertyOption)
	}
	return p
}

package rxpd

import (
	"fmt"
	"time"

	"github.com/philippseith/gorx/pkg/rx"
)

type Subscribable[T any] interface {
	rx.Subscribable[T]

	SetInterval(time.Duration)
	Interval() time.Duration

	Read() <-chan rx.Result[T]
	Write(T) <-chan error
}

type Property[T any] interface {
	Subscribable[T]
	rx.ObservableExtension[Property[T], Subscribable[T], T]

	ToAny() Property[any]
}

// ToProperty extends a Subscribable to a Property
func ToProperty[T any](s rx.Subscribable[T], options ...PropertyOption[T]) Property[T] {
	p := &property[T]{Subscribable: s}
	for _, option := range options {
		option(&p.propertyOption)
	}
	return p
}

type PropertyOption[T any] func(*propertyOption[T])

type propertyOption[T any] struct {
	setInterval func(time.Duration)
	read        func() <-chan rx.Result[T]
	write       func(T) <-chan error
}

func WithSetInterval[T any](setInterval func(time.Duration)) func(*propertyOption[T]) {
	return func(po *propertyOption[T]) {
		po.setInterval = setInterval
	}
}

func WithRead[T any](read func() <-chan rx.Result[T]) func(*propertyOption[T]) {
	return func(po *propertyOption[T]) {
		po.read = read
	}
}

func WithWrite[T any](write func(T) <-chan error) func(*propertyOption[T]) {
	return func(po *propertyOption[T]) {
		po.write = write
	}
}

type property[T any] struct {
	rx.Subscribable[T]
	propertyOption[T]

	tearDownLogics []func()
	interval       time.Duration
}

func (p *property[T]) Subscribe(o rx.Observer[T]) rx.Subscription {
	sub := p.Subscribable.Subscribe(o)
	for _, tld := range p.tearDownLogics {
		sub.AddTearDownLogic(tld)
	}
	return sub
}

func (p *property[T]) AddTearDownLogic(logic func()) Property[T] {
	p.tearDownLogics = append(p.tearDownLogics, logic)
	return p
}

func (p *property[T]) Catch(catch func(error) Subscribable[T]) Property[T] {
	// TODO the returned Property has undefined Read/Write behavior
	return ToProperty[T](rx.Catch[T](p, func(err error) rx.Subscribable[T] {
		return catch(err)
	}))
}

func (p *property[T]) Concat(sources ...Subscribable[T]) Property[T] {
	// TODO the returned Property has undefined Read/Write behavior
	// TODO Which of the sources should be read/written -> The one that will send next, but how to find out which is this?
	rxSources := make([]rx.Subscribable[T], 0, len(sources))
	for _, source := range sources {
		rxSources = append(rxSources, source)
	}
	return ToProperty[T](rx.Concat[T](append([]rx.Subscribable[T]{p}, rxSources...)...))
}

func (p *property[T]) DebounceTime(duration time.Duration) Property[T] {
	return p.toProperty(rx.DebounceTime[T](p, duration))
}

func (p *property[T]) DistinctUntilChanged(equal func(T, T) bool) Property[T] {
	return p.toProperty(rx.DistinctUntilChanged[T](p, equal))
}

func (p *property[T]) Log(id string) Property[T] {
	return p.toProperty(rx.Log[T](p, id))
}

func (p *property[T]) Merge(sources ...Subscribable[T]) Property[T] {
	//TODO implement me
	panic("implement me")
}

func (p *property[T]) Share() Property[T] {
	return p.toProperty(rx.Share[T](p))
}

func (p *property[T]) ShareReplay(opts ...rx.ReplayOption) Property[T] {
	return p.toProperty(rx.ShareReplay[T](p, opts...))
}

func (p *property[T]) Take(count int) Property[T] {
	return p.toProperty(rx.Take[T](p, count))
}

func (p *property[T]) Tap(subscribe func(rx.Observer[T]), next func(T) T, err func(error) error, complete, unsubscribe func()) Property[T] {
	return p.toProperty(rx.Tap[T](p, subscribe, next, err, complete, unsubscribe))
}

func (p *property[T]) ToAny() Property[any] {
	var read func() <-chan rx.Result[any]
	if p.read != nil {
		read = func() <-chan rx.Result[any] {
			ch := make(chan rx.Result[any], 1)
			if rT, ok := <-p.read(); ok {
				ch <- rx.Result[any]{Ok: rT.Ok, Err: rT.Err}
			}
			close(ch)
			return ch
		}
	}
	var write func(value any) <-chan error
	if p.write != nil {
		write = func(value any) <-chan error {
			ch := make(chan error, 1)
			if valueT, ok := value.(T); ok {
				if err, ok := <-p.write(valueT); ok {
					ch <- err
				}
			} else {
				ch <- fmt.Errorf("can not convert %v to %T", value, valueT)
			}
			close(ch)
			return ch
		}
	}

	return &property[any]{
		Subscribable: rx.ToAny[T](p.Subscribable),
		interval:     p.interval,
		propertyOption: propertyOption[any]{
			setInterval: p.setInterval,
			read:        read,
			write:       write,
		},
	}
}

func (p *property[T]) ToSlice() <-chan []T {
	return rx.ToSlice[T](p)
}

func (p *property[T]) SetInterval(interval time.Duration) {
	p.interval = interval
	if p.setInterval != nil {
		p.setInterval(interval)
	}
}

func (p *property[T]) Interval() time.Duration {
	return p.interval
}

func (p *property[T]) Read() <-chan rx.Result[T] {
	if p.read != nil {
		return p.Read()
	}
	ch := make(chan rx.Result[T])
	close(ch)
	return ch
}

func (p *property[T]) Write(value T) <-chan error {
	if p.write != nil {
		return p.write(value)
	}
	ch := make(chan error)
	close(ch)
	return ch
}

func (p *property[T]) toProperty(s rx.Subscribable[T]) Property[T] {
	return &property[T]{
		Subscribable: s,
		interval:     p.interval,
		propertyOption: propertyOption[T]{
			setInterval: p.setInterval,
			read:        p.read,
			write:       p.write,
		},
	}
}

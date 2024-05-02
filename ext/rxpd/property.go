package rxpd

import (
	"fmt"
	"runtime/debug"
	"slices"
	"sync"
	"time"

	"github.com/philippseith/gorx/pkg/rx"
)

type Subscribable[T any] interface {
	rx.Subscribable[T]

	SetInterval(time.Duration)
	Interval() time.Duration

	Read() rx.ResultChan[T]
	Write(T) <-chan error
}

type Property[T any] interface {
	Subscribable[T]
	rx.ObservableExtension[Property[T], Subscribable[T], T]

	ToAny() Property[any]
	With(options ...PropertyOption[T]) Property[T]
}

// ToProperty extends a Subscribable to a Property
func ToProperty[T any](s rx.Subscribable[T], options ...PropertyOption[T]) Property[T] {
	p := &property[T]{Subscribable: s}
	for _, option := range options {
		option(&p.propertyOption)
	}
	return p
}

type property[T any] struct {
	rx.Subscribable[T]
	propertyOption[T]

	tearDownLogics []func()
	interval       time.Duration
	readNextCh     []chan T
	mxRNC          sync.Mutex
}

func (p *property[T]) Subscribe(o rx.Observer[T]) rx.Subscription {
	sub := p.Subscribable.Subscribe(o)
	for _, tld := range p.tearDownLogics {
		sub.AddTearDownLogic(tld)
	}
	if p.readCallsNext {

		readNextCh := make(chan T, 1)
		p.readNextCh = append(p.readNextCh, readNextCh)
		sub.AddTearDownLogic(func() {
			p.mxRNC.Lock()
			defer p.mxRNC.Unlock()

			if i := slices.Index(p.readNextCh, readNextCh); i != -1 {
				p.readNextCh[i] = p.readNextCh[len(p.readNextCh)-1]
				p.readNextCh = p.readNextCh[:len(p.readNextCh)-1]
			}
			close(readNextCh)
		})

		go func() {
			for value := range readNextCh {
				o.Next(value)
			}
		}()
	}
	return sub
}

func (p *property[T]) AddTearDownLogic(logic func()) Property[T] {
	p.tearDownLogics = append(p.tearDownLogics, logic)
	return p
}

func (p *property[T]) Catch(catch func(error) Subscribable[T]) Property[T] {
	q := &property[T]{
		interval: p.interval,
		propertyOption: propertyOption[T]{
			setInterval: p.setInterval,
			read:        p.read,
			write:       p.write,
		},
	}
	q.Subscribable = rx.Catch[T](p, func(err error) rx.Subscribable[T] {
		c := catch(err)
		q.interval = c.Interval()
		q.setInterval = c.SetInterval
		q.read = c.Read
		q.write = c.Write

		return c
	})

	return q
}

// Concat creates an output Property which sequentially emits all values from
// the first given Subscribable and then moves on to the next.
// Interval, SetInterval, Read, Write of the output Property do not invoke any method of the sources.
// If they should have an effect, they have to be set with PropertyOptions in the ToProperty or With methods.
func (p *property[T]) Concat(sources ...Subscribable[T]) Property[T] {
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

// Merge subscribes to each given input Subscribable (as arguments), and simply
// forwards (without doing any transformation) all the values from all the input
// Subscribables to the output Property. The output Property only completes
// once all input Subscribables have completed. Any error delivered by an input
// Subscribable will be immediately emitted on the output Property.
// Interval, SetInterval, Read, Write of the output Property do not invoke any method of the sources.
// If they should have an effect, they have to be set with PropertyOptions in the ToProperty or With methods.
func (p *property[T]) Merge(sources ...Subscribable[T]) Property[T] {
	rxSources := make([]rx.Subscribable[T], 0, len(sources))
	for _, source := range sources {
		rxSources = append(rxSources, source)
	}
	return ToProperty[T](rx.Merge[T](rxSources...))
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
	return ToAny[T](p)
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

func (p *property[T]) Read() rx.ResultChan[T] {
	ch := make(chan rx.Result[T])
	defer close(ch)

	if p.read != nil {
		return func() (rCh rx.ResultChan[T]) {
			defer func() {
				if r := recover(); r != nil {
					ch <- rx.Result[T]{Err: fmt.Errorf("panic in %T.Read(): %v.\n%s", p, r, string(debug.Stack()))}
					rCh = ch
				}
			}()

			switch {
			case p.readTimeout > 0:
				rCh = p.onReadWithTimeout(p.read())
			case p.readCallsNext && p.readTimeout == 0:
				rCh = p.onReadCallsNext(p.read())
			default:
				// Trick to make sending an error to the result chan in defer possible
				rCh = p.read()
			}
			return rCh
		}()
	}
	return ch
}

func (p *property[T]) onReadCallsNext(rCh rx.ResultChan[T]) rx.ResultChan[T] {
	ch := make(chan rx.Result[T], 1)
	go func() {
		defer close(ch)

		result := <-rCh
		ch <- result
		if result.Err == nil {
			p.callNextWithRead(result.Ok)
		}
	}()
	return ch
}

func (p *property[T]) callNextWithRead(value T) {
	p.mxRNC.Lock()
	defer p.mxRNC.Unlock()

	for _, readNextCh := range p.readNextCh {
		go func(ch chan T) {
			ch <- value
		}(readNextCh)
	}
}

func (p *property[T]) onReadWithTimeout(rCh rx.ResultChan[T]) rx.ResultChan[T] {
	ch := make(chan rx.Result[T], 1)
	go func() {
		defer close(ch)

		if p.readCallsNext {
			rCh = p.onReadCallsNext(rCh)
		}

		select {
		case result := <-rCh:
			ch <- result
		case <-time.After(p.readTimeout):
			timeoutResult := p.onReadTimeout()
			ch <- timeoutResult
			if timeoutResult.Err == nil {
				p.callNextWithRead(timeoutResult.Ok)
			}

		}
	}()
	return ch
}

func (p *property[T]) Write(value T) <-chan error {
	ch := make(chan error)
	defer close(ch)

	if p.write != nil {
		return func() (wCh <-chan error) {
			defer func() {
				if r := recover(); r != nil {
					ch <- fmt.Errorf("panic in %T.Write(%v): %v.\n%s", p, value, r, string(debug.Stack()))
					wCh = ch
				}
			}()
			wCh = p.write(value)
			return wCh
		}()
	}
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

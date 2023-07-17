package rx

import (
	"context"
	"sync"

	"golang.org/x/exp/slices"
)

type Subscribable[T any] interface {
	Subscribe(o Observer[T]) Subscription
}

func FromChan[T any](ctx context.Context, ch <-chan T) Subscribable[T] {
	p := &pipe[T, T]{
		n: func(t T) T {
			return t
		},
	}
	go func() {
		for {
			select {
			case t, ok := <-ch:
				if !ok {
					p.Complete()
					return
				}
				p.Next(t)
			case <-ctx.Done():
				p.Complete()
				return
			}
		}
	}()
	return p
}

func ToChan[T any](s Subscribable[T]) <-chan T {
	panic("imp!")
}

type toChan[T any] struct {
	ch chan T
}

func (tc *toChan[T]) Next(value T) {

}

func Pipe[T any, U any](source Subscribable[T], next func(value T) U) Subscribable[U] {
	p := &pipe[T, U]{
		n: next,
	}
	source.Subscribe(p)
	return p
}

type pipe[T any, U any] struct {
	obs   []Observer[U]
	mxObs sync.RWMutex
	n     func(value T) U
}

func (p *pipe[T, U]) Subscribe(o Observer[U]) Subscription {
	p.obs = append(p.obs, o)
	return &subscription{func() {
		p.mxObs.Lock()
		defer p.mxObs.Unlock()

		idx := slices.Index(p.obs, o)
		p.obs = append(p.obs[:idx], p.obs[idx+1:]...)
	}}
}

func (p *pipe[T, U]) Next(t T) {
	u := p.n(t)
	p.mxObs.RLock()
	defer p.mxObs.RUnlock()

	for _, ob := range p.obs {
		ob.Next(u)
	}
}

func (p *pipe[T, U]) Error(err error) {
	p.mxObs.RLock()
	defer p.mxObs.RUnlock()

	for _, ob := range p.obs {
		ob.Error(err)
	}
}

func (p *pipe[T, U]) Complete() {
	p.mxObs.RLock()
	defer p.mxObs.RUnlock()

	for _, ob := range p.obs {
		ob.Complete()
	}
}

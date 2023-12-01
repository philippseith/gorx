package rx_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestTapRef(t *testing.T) {
	x := 0
	s := rx.NewSubject[*int]()
	ta := s.Log("before").Tap(nil, func(ctx context.Context, i *int) (context.Context, *int) {
		*i = 1
		return ctx, i
	}, nil, nil, nil).Log("after")

	var y int
	ta.Subscribe(rx.OnNext(func(i *int) {
		y = *i
	}))

	s.Next(context.Background(), &x)
	assert.Equal(t, 1, y)
}

func TestTap(t *testing.T) {
	s := rx.NewSubject[int]()
	ta := s.Log("before").Tap(nil, func(ctx context.Context, i int) (context.Context, int) {
		return ctx, i + 1
	}, nil, nil, nil).Log("after")

	var y int
	ta.Subscribe(rx.OnNext(func(i int) {
		y = i
	}))

	s.Next(context.Background(), 0)
	assert.Equal(t, 1, y)
}

func TestPanic(t *testing.T) {
	s := rx.NewSubject[int]()
	tn := s.Log("before").Tap(nil, func(context.Context, int) (context.Context, int) {
		panic("Next")
	}, nil, nil, nil).Log("after")

	var panicNext error
	tn.Subscribe(rx.NewObserver[int](nil, func(err error) {
		panicNext = err
	}, nil))
	s.Next(context.Background(), 0)
	assert.NotNil(t, panicNext)
}

package rx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestTapRef(t *testing.T) {
	x := 0
	s := rx.NewSubject[*int]()
	ta := s.Log("before").Tap(nil, func(i *int) *int {
		*i = 1
		return i
	}, nil, nil, nil).Log("after")

	var y int
	ta.Subscribe(rx.OnNext(func(i *int) {
		y = *i
	}))

	s.Next(&x)
	assert.Equal(t, 1, y)
}

func TestTap(t *testing.T) {
	s := rx.NewSubject[int]()
	ta := s.Log("before").Tap(nil, func(i int) int {
		return i + 1
	}, nil, nil, nil).Log("after")

	var y int
	ta.Subscribe(rx.OnNext(func(i int) {
		y = i
	}))

	s.Next(0)
	assert.Equal(t, 1, y)
}

func TestPanic(t *testing.T) {
	s := rx.NewSubject[int]()
	tn := s.Log("before").Tap(nil, func(i int) int {
		panic("Next")
	}, nil, nil, nil).Log("after")

	var panicNext error
	tn.Subscribe(rx.NewObserver[int](nil, func(err error) {
		panicNext = err
	}, nil))
	s.Next(0)
	assert.NotNil(t, panicNext)
}

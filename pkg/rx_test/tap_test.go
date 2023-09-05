package rx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestTapRef(t *testing.T) {
	x := 0
	s := rx.NewSubject[*int]()
	ta := s.Tap(nil, func(i *int) *int {
		*i = 1
		return i
	}, nil, nil, nil)

	var y int
	ta.Subscribe(rx.OnNext(func(i *int) {
		y = *i
	}))

	s.Next(&x)
	assert.Equal(t, 1, y)
}

func TestTap(t *testing.T) {
	s := rx.NewSubject[int]()
	ta := s.Tap(nil, func(i int) int {
		return i + 1
	}, nil, nil, nil)

	var y int
	ta.Subscribe(rx.OnNext(func(i int) {
		y = i
	}))

	s.Next(0)
	assert.Equal(t, 1, y)
}

package rx_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestDistinctUntilChanged(t *testing.T) {
	s := rx.NewSubject[int]()
	var actual []int
	rx.DistinctUntilChanged[int](s, func(a, b int) bool { return a == b }).
		Subscribe(rx.NewObserver(func(value int) {
			actual = append(actual, value)
		}, nil, nil))
	for _, i := range []int{1, 1, 2, 2, 3, 3} {
		s.Next(i)
	}
	assert.Equal(t, []int{1, 2, 3}, actual)
}

func TestDistinct(t *testing.T) {
	s := rx.NewSubject[int]()
	var actual []int
	rx.Distinct[int](s).
		Subscribe(rx.NewObserver(func(value int) {
			actual = append(actual, value)
		}, nil, nil))
	for _, i := range []int{1, 1, 2, 2, 3, 3, 1, 2, 4} {
		s.Next(i)
	}
	assert.Equal(t, []int{1, 2, 3, 4}, actual)
}

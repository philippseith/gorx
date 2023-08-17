package rx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestSubjectNext(t *testing.T) {
	s := rx.NewSubject[int]()

	var ov []int
	o := rx.NewObserver[int](func(v int) {
		ov = append(ov, v)
	}, nil, nil)
	var pv []int
	p := rx.NewObserver[int](func(v int) {
		pv = append(pv, v)
	}, nil, nil)

	ou := s.Subscribe(o)
	pu := s.Subscribe(p)

	n1 := []int{1, 2, 3}
	for _, v := range n1 {
		s.Next(v)
	}
	assert.Equal(t, n1, ov)
	assert.Equal(t, n1, pv)

	pu.Unsubscribe()

	s.Next(4)

	assert.Equal(t, []int{1, 2, 3, 4}, ov)
	assert.Equal(t, n1, pv)

	ou.Unsubscribe()

	s.Next(5)

	assert.Equal(t, []int{1, 2, 3, 4}, ov)
	assert.Equal(t, n1, pv)
}

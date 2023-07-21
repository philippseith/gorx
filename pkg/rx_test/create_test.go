package rx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestCreate(t *testing.T) {
	c := rx.Create(func(o rx.Observer[int]) {
		o.Next(1)
		o.Next(2)
	})

	var ov []int
	o := rx.NewObserver[int](func(v int) {
		ov = append(ov, v)
	}, nil, nil)
	var pv []int
	p := rx.NewObserver[int](func(v int) {
		pv = append(pv, v)
	}, nil, nil)

	c.Subscribe(o)
	assert.Equal(t, []int{1, 2}, ov)
	c.Subscribe(p)
	assert.Equal(t, []int{1, 2}, pv)
}

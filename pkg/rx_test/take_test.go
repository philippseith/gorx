package rx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestTake(t *testing.T) {
	c := rx.From(1, 2, 3, 4).ToConnectable()
	ta := rx.Take[int](c, 2)
	var actual []int
	ta.Subscribe(rx.OnNext(func(t int) {
		actual = append(actual, t)
	}))
	c.Connect()
	assert.Equal(t, []int{1, 2}, actual)
}

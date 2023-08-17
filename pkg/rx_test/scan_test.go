package rx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestScan(t *testing.T) {
	c := rx.ToConnectable[int](rx.From(1, 2, 3))
	s := rx.Scan[int, int](c, func(acc int, val int) int { return acc + val }, 0)
	var actual []int
	s.Subscribe(rx.OnNext(func(t int) {
		actual = append(actual, t)
	}))
	c.Connect()
	assert.Equal(t, []int{1, 3, 6}, actual)
}

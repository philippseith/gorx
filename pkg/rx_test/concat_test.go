package rx_test

import (
	"testing"

	"github.com/philippseith/gorx/pkg/rx"
	"github.com/stretchr/testify/assert"
)

func TestConcat(t *testing.T) {
	c1 := rx.From[int](1, 2, 3)
	c2 := rx.From[int](4, 5, 6)
	cc := c1.Concat(c2)

	var ii []int
	cc.Subscribe(rx.OnNext[int](func(i int) {
		ii = append(ii, i)
	}))
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, ii)
}

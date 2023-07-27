package rx_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestTake(t *testing.T) {
	ta := rx.Take[int](rx.From(1, 2, 3, 4), 2)
	actual := <-rx.ToSlice[int](ta)
	assert.Equal(t, []int{1, 2}, actual)
}

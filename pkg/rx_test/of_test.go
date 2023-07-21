package rx_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestOf(t *testing.T) {
	expected := []int{1, 2, 3}
	o := rx.Of(expected...)
	var actual []int
	o.Subscribe(rx.NewObserver[int](func(i int) {
		actual = append(actual, i)
	}, nil, nil))

	assert.Equal(t, expected, actual)
}

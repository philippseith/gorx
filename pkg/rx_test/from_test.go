package rx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestJust(t *testing.T) {
	expected := []int{1, 2, 3}
	f := rx.From(expected...)
	var actual []int
	f.Subscribe(rx.NewObserver[int](func(i int) {
		actual = append(actual, i)
	}, nil, nil))

	assert.Equal(t, expected, actual)
}

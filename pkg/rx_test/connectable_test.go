package rx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestConnectable(t *testing.T) {
	expected := []int{1, 2, 3}
	fromUnsubscribed := false
	c := rx.From(expected...).AddTearDownLogic(func() {
		fromUnsubscribed = true
	}).ToConnectable()

	var actual []int
	sub := c.Subscribe(rx.NewObserver[int](func(i int) {
		actual = append(actual, i)
	}, nil, nil))

	assert.Empty(t, actual)

	c.Connect()

	assert.Equal(t, expected, actual)

	sub.Unsubscribe()
	assert.True(t, fromUnsubscribed)
}

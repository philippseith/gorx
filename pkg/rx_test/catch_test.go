package rx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestCatch(t *testing.T) {
	sourceUnsubscribed := false
	catchUnsubscribed := false
	c := rx.Create(func(o rx.Observer[int]) rx.Subscription {
		o.Next(1)
		o.Next(2)
		o.Error(nil)
		return rx.NewSubscription(func() {
			sourceUnsubscribed = true
		})
	}).ToConnectable()

	cc := c.Catch(func(err error) rx.Subscribable[int] {
		return rx.Create(func(o rx.Observer[int]) rx.Subscription {
			o.Next(3)
			o.Next(4)
			return rx.NewSubscription(func() {
				catchUnsubscribed = true
			})
		})
	})

	var actual []int
	sub := cc.Subscribe(rx.OnNext(func(next int) {
		actual = append(actual, next)
	}))

	c.Connect()

	sub.Unsubscribe()
	assert.True(t, sourceUnsubscribed)
	assert.True(t, catchUnsubscribed)
	assert.Equal(t, []int{1, 2, 3, 4}, actual)
}

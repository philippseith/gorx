package rx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestShareReplayWithRefCount(t *testing.T) {
	subscribedOnce := false
	source := rx.Create[int](func(o rx.Observer[int]) rx.Subscription {
		if !subscribedOnce {
			o.Next(0)
			o.Next(1)
			o.Next(2)
			subscribedOnce = true
		} else {
			o.Next(314)
		}
		return rx.NewSubscription(func() {})
	})

	sr := source.ShareReplay(rx.MaxBufferSize(2), rx.RefCount(true))
	var act1 []int
	sub1 := sr.Subscribe(rx.OnNext(func(t int) {
		act1 = append(act1, t)
	}))
	assert.Equal(t, []int{0, 1, 2}, act1)
	var act2 []int
	sub2 := sr.Subscribe(rx.OnNext(func(t int) {
		act2 = append(act2, t)
	}))
	assert.Equal(t, []int{1, 2}, act2)

	sub1.Unsubscribe()
	sub2.Unsubscribe()

	var act3 []int
	sr.Subscribe(rx.OnNext(func(t int) {
		act3 = append(act3, t)
	}))
	assert.Equal(t, []int{1, 2, 314}, act3)
}

func TestShareReplayWithoutRefCount(t *testing.T) {
	subscribedOnce := false
	source := rx.Create[int](func(o rx.Observer[int]) rx.Subscription {
		if !subscribedOnce {
			o.Next(0)
			o.Next(1)
			o.Next(2)
			subscribedOnce = true
		} else {
			o.Next(314)
		}
		return rx.NewSubscription(func() {})
	})

	sr := source.ShareReplay(rx.MaxBufferSize(2), rx.RefCount(false))
	var act1 []int
	sub1 := sr.Subscribe(rx.OnNext(func(t int) {
		act1 = append(act1, t)
	}))
	assert.Equal(t, []int{0, 1, 2}, act1)

	sub1.Unsubscribe()

	var act2 []int
	sr.Subscribe(rx.OnNext(func(t int) {
		act2 = append(act2, t)
	}))
	assert.Equal(t, []int{1, 2}, act2)
}

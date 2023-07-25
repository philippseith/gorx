package rx_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestShare(t *testing.T) {
	ch := make(chan int, 1)
	sh := rx.Share[int](rx.FromChan[int](context.Background(), ch))

	var wg sync.WaitGroup
	wg.Add(2)
	var sl1, sl2 []int
	sh.Subscribe(rx.NewObserver(func(value int) {
		sl1 = append(sl1, value)
	}, nil, func() {
		wg.Done()
	}))
	sh.Subscribe(rx.NewObserver(func(value int) {
		sl2 = append(sl2, value)
	}, nil, func() {
		wg.Done()
	}))

	for i := 1; i < 4; i++ {
		ch <- i
	}
	close(ch)
	wg.Wait()
	expected := []int{1, 2, 3}
	assert.Equal(t, expected, sl1)
	assert.Equal(t, expected, sl2)
}

func TestShareResetOnNewSubscriber0(t *testing.T) {
	sh := rx.Share(rx.From(1, 2, 3))
	expected := []int{1, 2, 3}
	var sl1 []int
	sh.Subscribe(rx.NewObserver(func(value int) {
		sl1 = append(sl1, value)
	}, nil, nil)).Unsubscribe()
	assert.Equal(t, expected, sl1)

	var sl2 []int
	sh.Subscribe(rx.NewObserver(func(value int) {
		sl2 = append(sl2, value)
	}, nil, nil))
	assert.Equal(t, expected, sl2)

	// Without Unsubscribe, rx.From will not be restarted
	var sl3 []int
	sh.Subscribe(rx.NewObserver(func(value int) {
		sl3 = append(sl3, value)
	}, nil, nil))
	assert.Equal(t, 0, len(sl3))
}

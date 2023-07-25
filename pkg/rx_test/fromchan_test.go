package rx_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestFromChan(t *testing.T) {
	ch := make(chan int, 1)
	fc := rx.FromChan[int](context.Background(), ch)

	var sl []int
	done := make(chan struct{}, 1)
	fc.Subscribe(rx.NewObserver[int](func(value int) {
		sl = append(sl, value)
	}, nil, func() {
		close(done)
	}))

	for i := 0; i < 3; i++ {
		ch <- i
	}
	close(ch)
	<-done
	assert.Equal(t, []int{0, 1, 2}, sl)
}

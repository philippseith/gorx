package rxpd_test

import (
	"context"
	"testing"
	"time"

	"github.com/philippseith/gorx/ext/rxpd"
	"github.com/philippseith/gorx/pkg/rx"
	"github.com/stretchr/testify/assert"
)

func TestReadWithTimeout(t *testing.T) {
	value := 0

	p := rxpd.FromRead[int](context.Background(), func() rx.ResultChan[int] {
		ch := make(chan rx.Result[int], 1)
		go func() {
			<-time.After(50 * time.Millisecond)
			ch <- rx.Ok(value)
			close(ch)
		}()
		return ch
	}, rxpd.WithReadTimeout[int](20*time.Millisecond, func() rx.Result[int] { return rx.Ok(-1) }))

	result := <-p.Read()
	assert.Equal(t, -1, result.Ok)
}
func TestReadWithTimeoutAndCallsNext(t *testing.T) {
	value := 0

	p := rxpd.FromRead[int](context.Background(), func() rx.ResultChan[int] {
		ch := make(chan rx.Result[int], 1)
		go func() {
			<-time.After(50 * time.Millisecond)
			ch <- rx.Ok(value)
			close(ch)
		}()
		return ch
	},
		rxpd.WithReadTimeout[int](20*time.Millisecond, func() rx.Result[int] { return rx.Ok(-1) }),
		rxpd.WithReadCallsNext[int]())

	ch := make(chan struct{})
	nextValue := []int{0}
	p.Subscribe(rx.OnNext(func(t int) {
		nextValue[0] = t
		ch <- struct{}{}
	}))

	result := <-p.Read()
	assert.Equal(t, -1, result.Ok)
	<-ch
	assert.Equal(t, -1, nextValue[0])
	<-ch
	assert.Equal(t, 0, nextValue[0])
}

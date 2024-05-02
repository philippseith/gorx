package rxpd_test

import (
	"context"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/philippseith/gorx/ext/rxpd"
	"github.com/philippseith/gorx/pkg/rx"
)

func TestReadCallsNext(t *testing.T) {
	value := 0
	p := rxpd.FromRead[int](context.Background(), func() rx.ResultChan[int] {
		ch := make(chan rx.Result[int], 1)
		ch <- rx.Ok(value)
		close(ch)
		return ch
	}, rxpd.WithReadCallsNext[int]())

	ch := make(chan struct{})
	nextValue := []int{0, 0}
	sub1 := p.Subscribe(rx.OnNext(func(t int) {
		nextValue[0] = t
		ch <- struct{}{}
	}))
	p.Subscribe(rx.OnNext(func(t int) {
		nextValue[1] = t
		ch <- struct{}{}
	}))

	value = 1
	<-p.Read()
	<-ch
	<-ch

	assert.Equal(t, value, nextValue[0])
	assert.Equal(t, value, nextValue[1])

	sub1.Unsubscribe()
	value = 2
	<-p.Read()
	<-ch

	assert.Equal(t, 1, nextValue[0])
	assert.Equal(t, value, nextValue[1])

}

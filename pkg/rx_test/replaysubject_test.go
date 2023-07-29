package rx_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestReplaySubjectBuffer(t *testing.T) {
	r := rx.NewReplaySubject[int](rx.MaxBufferSize(2))

	r.Next(1)
	r.Next(2)
	r.Next(3)

	var a1 []int
	r.Subscribe(rx.OnNext(func(value int) {
		a1 = append(a1, value)
	}))

	assert.Equal(t, []int{2, 3}, a1)

	var a2 []int
	r.Subscribe(rx.OnNext(func(value int) {
		a2 = append(a2, value)
	}))

	assert.Equal(t, []int{2, 3}, a2)

	r.Next(4)

	var a3 []int
	r.Subscribe(rx.OnNext(func(value int) {
		a3 = append(a3, value)
	}))

	assert.Equal(t, []int{3, 4}, a3)
}

func TestReplaySubjectEndlessBuffer(t *testing.T) {
	r := rx.NewReplaySubject[int]()

	r.Next(1)

	var a1 []int
	r.Subscribe(rx.OnNext(func(value int) {
		a1 = append(a1, value)
	}))

	assert.Equal(t, []int{1}, a1)

	r.Next(2)

	var a2 []int
	r.Subscribe(rx.OnNext(func(value int) {
		a2 = append(a2, value)
	}))

	assert.Equal(t, []int{1, 2}, a2)

	r.Next(3)

	var a3 []int
	r.Subscribe(rx.OnNext(func(value int) {
		a3 = append(a3, value)
	}))

	assert.Equal(t, []int{1, 2, 3}, a3)
}

func TestReplaySubjectWindow(t *testing.T) {
	r := rx.NewReplaySubject[int](rx.Window(5 * time.Millisecond))

	var a1 []int
	r.Subscribe(rx.OnNext(func(value int) {
		a1 = append(a1, value)
	}))
	r.Next(1)

	assert.Equal(t, []int{1}, a1)

	<-time.After(5 * time.Millisecond)
	var a2 []int
	r.Subscribe(rx.OnNext(func(value int) {
		a2 = append(a2, value)
	}))

	assert.Equal(t, []int(nil), a2)

	r.Next(2)
	<-time.After(1 * time.Millisecond)
	r.Next(3)
	<-time.After(1 * time.Millisecond)

	assert.Equal(t, []int{2, 3}, a2)

	<-time.After(3 * time.Millisecond)

	var a3 []int
	r.Subscribe(rx.OnNext(func(value int) {
		a3 = append(a3, value)
	}))

	assert.Equal(t, []int{3}, a3)
}

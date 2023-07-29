package rx_test

import (
	"errors"
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

func TestReplaySubjectComplete(t *testing.T) {
	r := rx.NewReplaySubject[int](rx.MaxBufferSize(2))

	r.Next(1)
	r.Next(2)
	r.Next(3)
	r.Complete()

	var a1 []int
	var c1 bool
	r.Subscribe(rx.NewObserver(func(value int) {
		a1 = append(a1, value)
	}, nil, func() {
		c1 = true
	}))

	assert.Equal(t, []int{2, 3}, a1)
	assert.True(t, c1)

	var a2 []int
	var c2 bool
	r.Subscribe(rx.NewObserver(func(value int) {
		a2 = append(a2, value)
	}, nil, func() {
		c2 = true
	}))

	assert.Equal(t, []int{2, 3}, a2)
	assert.True(t, c2)
}

func TestReplaySubjectError(t *testing.T) {
	r := rx.NewReplaySubject[int](rx.MaxBufferSize(2))

	r.Next(1)
	r.Next(2)
	r.Next(3)
	err := errors.New("shit happens")
	r.Error(err)

	var a1 []int
	var err1 error
	r.Subscribe(rx.NewObserver(func(value int) {
		a1 = append(a1, value)
	}, func(err error) {
		err1 = err
	}, nil))

	assert.Equal(t, []int{2, 3}, a1)
	assert.Equal(t, err, err1)

	var a2 []int
	var err2 error
	r.Subscribe(rx.NewObserver(func(value int) {
		a2 = append(a2, value)
	}, func(err error) {
		err2 = err
	}, nil))

	assert.Equal(t, []int{2, 3}, a2)
	assert.Equal(t, err, err2)
}

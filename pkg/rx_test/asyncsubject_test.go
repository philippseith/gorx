package rx_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestAsyncSubject_EmitsValueOnComplete(t *testing.T) {
	flags := []bool{false, false, false}
	r := 0
	rr := &r
	as := rx.NewAsyncSubject[int]()

	as.Subscribe(rx.NewObserver[int](func(v int) {
		assert.True(t, flags[0])
		*rr = v
	}, func(error) {
		flags[1] = true
	}, func() {
		assert.True(t, flags[0])
		flags[2] = true
	}))

	as.Next(1)
	assert.Equal(t, 0, r)
	as.Next(99)
	assert.Equal(t, 0, r)
	assert.False(t, flags[1])
	assert.False(t, flags[2])
	flags[0] = true
	as.Complete()
	assert.Equal(t, 99, r)
	assert.False(t, flags[1])
	assert.True(t, flags[2])
}

func TestAsyncSubject_DoesEmitError(t *testing.T) {
	flags := []bool{false}
	err := errors.New("error")

	as := rx.NewAsyncSubject[struct{}]()
	as.Subscribe(rx.NewObserver[struct{}](nil, func(e error) {
		assert.Equal(t, err, e)
		flags[0] = true
	}, nil))
	as.Error(err)

	assert.True(t, flags[0])
}

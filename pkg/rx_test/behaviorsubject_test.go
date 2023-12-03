package rx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestBehaviorSubject_DoesEmitInitialValueOnSubscribe(t *testing.T) {
	r := 0
	rr := &r
	bs := rx.NewBehaviorSubject[int](context.Background(), 99)

	bs.Subscribe(rx.NewObserver(func(v int) {
		*rr = v
	}, nil, nil))

	assert.Equal(t, r, 99)

	r = 0
	bs.Subscribe(rx.NewObserver(func(v int) {
		*rr = v
	}, nil, nil))

	assert.Equal(t, r, 99)
}

func TestBehaviorSubject_DoesEmitError(t *testing.T) {
	flags := []bool{false}
	err := errors.New("error")

	bs := rx.NewBehaviorSubject[int](context.Background(), 0)
	bs.Subscribe(rx.NewObserver[int](nil, func(e error) {
		assert.Equal(t, err, e)
		flags[0] = true
	}, nil))
	bs.Error(context.Background(), err)

	assert.True(t, flags[0])
}

func TestBehaviorSubject_DoesComplete(t *testing.T) {
	flags := []bool{false}

	bs := rx.NewBehaviorSubject[int](context.Background(), 0)
	bs.Subscribe(rx.NewObserver[int](nil, nil, func() {
		flags[0] = true
	}))
	bs.Complete(context.Background())

	assert.True(t, flags[0])
}

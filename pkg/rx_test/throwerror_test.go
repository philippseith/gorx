package rx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestThrowError(t *testing.T) {
	err := errors.New("test")
	te := rx.ThrowError[any](context.Background(), err)

	var actual error

	te.Subscribe(rx.NewObserver[any](nil, func(err error) {
		actual = err
	}, nil))

	assert.Equal(t, err, actual)
}

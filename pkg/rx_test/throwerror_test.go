package rx_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestThrowError(t *testing.T) {
	err := errors.New("test")
	te := rx.ThrowError[any](err)

	var actual error

	te.Subscribe(rx.NewObserver[any](nil, func(err error) {
		actual = err
	}, nil))

	assert.Equal(t, err, actual)
}

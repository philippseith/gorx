package rx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/philippseith/gorx/pkg/rx"
	"github.com/stretchr/testify/assert"
)

func TestGo(t *testing.T) {
	x, err := rx.Go(context.Background(), func() (int, error) {
		return 1, nil
	})

	assert.Equal(t, 1, x)
	assert.NoError(t, err)

	y, err := rx.Go(context.Background(), func() (float64, error) {
		return 1.0, errors.New("fail")
	})

	assert.Equal(t, 1.0, y)
	assert.Error(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	z, err := rx.Go(ctx, func() (string, error) {
		return "z", nil
	})

	assert.Equal(t, "", z)
	assert.Equal(t, ctx.Err(), err)
}

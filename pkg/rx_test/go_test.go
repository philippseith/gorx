package rx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/philippseith/gorx/pkg/rx"
	"github.com/stretchr/testify/assert"
)

func TestGo(t *testing.T) {
	x := <-rx.Go(context.Background(), func() (int, error) {
		return 1, nil
	})

	assert.Equal(t, 1, x.Ok)
	assert.NoError(t, x.Err)

	y := <-rx.Go(context.Background(), func() (float64, error) {
		return 1.0, errors.New("fail")
	})

	assert.Equal(t, 1.0, y.Ok)
	assert.Error(t, y.Err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	z := <-rx.Go(ctx, func() (string, error) {
		return "z", nil
	})

	assert.Equal(t, "", z.Ok)
	assert.Equal(t, ctx.Err(), z.Err)
}

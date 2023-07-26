package rx_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestTake(t *testing.T) {
	ta := rx.Take(rx.From(1, 2, 3, 4), 2)
	actual := rx.ToSlice(context.Background(), ta)
	assert.Equal(t, []int{1, 2}, actual)
}

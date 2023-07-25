package rx_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestToSlice(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, rx.ToSlice[int](context.Background(), rx.From[int](1, 2, 3)))
}

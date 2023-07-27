package rx_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestScan(t *testing.T) {
	s := rx.Scan[int, int](rx.From(1, 2, 3), func(acc int, val int) int { return acc + val }, 0)
	assert.Equal(t, []int{1, 3, 6}, <-rx.ToSlice[int](s))
}

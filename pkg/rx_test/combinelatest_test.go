package rx_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestCombineLatest(t *testing.T) {
	var actual []any
	ab := rx.ToConnectable(rx.From[any]('A', 'B'))
	oneTwoThree := rx.ToConnectable(rx.From[any](1, 2, 3))

	rx.CombineLatest[[]any](func(n ...any) []any {
		var nn []any
		for _, ni := range n {
			nn = append(nn, ni)
		}
		return nn
	}, ab, oneTwoThree).Subscribe(rx.NewObserver[[]any](func(next []any) {
		actual = next
	}, nil, nil))

	ab.Connect()
	oneTwoThree.Connect()

	assert.Equal(t, []any{'B', 3}, actual)
}

func TestCombineLatestTicker(t *testing.T) {
	t1 := rx.NewTicker(0, 1*time.Millisecond)
	t2 := rx.NewTicker(0, 2*time.Millisecond)
	tt1 := rx.Scan[time.Time, int](rx.Take[time.Time](t1, 3), func(i int, t time.Time) int { return i + 1 }, 0)
	tt2 := rx.Scan[time.Time, int](rx.Take[time.Time](t2, 3), func(i int, t time.Time) int { return i + 1 }, 0)

	sl := rx.ToSlice(context.Background(), rx.CombineLatest[[]any](func(ts ...any) []any { return ts }, rx.ToAny(tt1), rx.ToAny(tt2)))
	assert.Equal(t, [][]any{{1, 1}, {2, 1}, {2, 2}, {3, 2}, {3, 3}}, sl)
}

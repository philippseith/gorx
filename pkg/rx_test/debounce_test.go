package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestDebounce(t *testing.T) {
	ticker := rx.NewTicker(context.Background(), 10*time.Millisecond, 10*time.Millisecond)
	tc := ticker.ToConnectable()
	counts := []int{0, 0}
	done := make(chan struct{})
	tc.Tap(nil, func(ctx context.Context, t time.Time) (context.Context, time.Time) {
		counts[0] = counts[0] + 1
		if counts[0] == 10 {
			ticker.Stop()
			done <- struct{}{}
		}
		return ctx, t
	}, nil, nil, nil).
		DebounceTime(20 * time.Millisecond).Subscribe(rx.OnNext(func(t time.Time) {
		counts[1]++
	}))
	tc.Connect()
	<-done
	assert.Equal(t, 10, counts[0])
	assert.InDelta(t, 5, counts[1], 1.1) // Jitter
}

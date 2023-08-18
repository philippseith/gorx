package rx_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestDebounce(t *testing.T) {
	ticker := rx.NewTicker(1*time.Millisecond, 1*time.Millisecond)
	tc := ticker.ToConnectable()
	counts := []int{0, 0}
	done := make(chan struct{})
	tc.Tap(func(t time.Time) time.Time {
		counts[0] = counts[0] + 1
		if counts[0] == 10 {
			ticker.Stop()
			done <- struct{}{}
		}
		return t
	}, nil, nil).
		DebounceTime(2 * time.Millisecond).Subscribe(rx.OnNext(func(t time.Time) {
		counts[1]++
	}))
	tc.Connect()
	<-done
	assert.Equal(t, 10, counts[0])
	assert.InDelta(t, 5, counts[1], 1.1) // Jitter
}

package rx_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestDebounce(t *testing.T) {
	ticker := rx.NewTicker(1*time.Millisecond, 1*time.Millisecond)
	s := rx.NewSubject[time.Time]()
	ticker.Subscribe(s)
	counts := []int{0, 0}
	done := make(chan struct{})
	s.Subscribe(rx.NewObserver[time.Time](func(t time.Time) {
		counts[0] = counts[0] + 1
		if counts[0] == 10 {
			ticker.Stop()
			done <- struct{}{}
		}
	}, nil, nil))

	rx.Debounce[time.Time](s, 2*time.Millisecond).Subscribe(rx.NewObserver(func(t time.Time) {
		counts[1]++
	}, nil, nil))
	<-done
	assert.Equal(t, 10, counts[0])
	assert.Equal(t, 5, counts[1])
}

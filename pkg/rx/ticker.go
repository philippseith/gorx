package rx

import (
	"sync"
	"time"
)

type Ticker interface {
	Observable[time.Time]

	Reset(interval time.Duration)
	Stop()
}

func NewTicker(due, interval time.Duration) Ticker {
	t := &ticker{
		Subject:  NewSubject[time.Time](),
		interval: interval,
	}
	go func() {
		t.Next(<-time.After(due))

		t.mx.Lock()
		defer t.mx.Unlock()

		if !t.stopped {
			t.ticker = time.NewTicker(t.interval)
			go func() {
				for now := range t.ticker.C {
					t.Next(now)
				}
				t.Complete()
			}()
		}
	}()
	return t
}

type ticker struct {
	Subject[time.Time]
	ticker   *time.Ticker
	interval time.Duration
	stopped  bool
	mx       sync.Mutex
}

func (t *ticker) Reset(interval time.Duration) {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.interval = interval
	if t.ticker != nil {
		t.ticker.Reset(interval)
	}
}

func (t *ticker) Stop() {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.stopped = true
	if t.ticker != nil {
		t.ticker.Stop()
	}
}

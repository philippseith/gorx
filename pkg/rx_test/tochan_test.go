package rx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestToChan(t *testing.T) {
	cn := rx.ToConnectable[int](rx.From(1, 2, 3))
	ch, _ := rx.ToChan[int](cn)
	var sl []int
	done := make(chan struct{})
	go func() {
		for value := range ch {
			sl = append(sl, value.V)
		}
		close(done)
	}()
	cn.Connect()
	<-done
	assert.Equal(t, []int{1, 2, 3}, sl)
}

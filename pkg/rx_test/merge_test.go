package rx_test

import (
	"github.com/philippseith/gorx/pkg/rx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMerge(t *testing.T) {
	t1 := rx.Map[time.Time, int](rx.NewTicker(0, 100*time.Millisecond), func(t time.Time) int {
		return 100
	}).ToConnectable()
	t2 := rx.Map[time.Time, int](rx.NewTicker(50*time.Millisecond, 100*time.Millisecond), func(t time.Time) int {
		return 150
	}).ToConnectable()
	m := rx.Merge[int](
		t1.Take(5),
		t2.Take(5))
	t1.Connect()
	t2.Connect()
	actual := <-m.ToSlice()

	assert.Len(t, actual, 10)

	// We don't know in which phase connect is called
	var expected []int
	if actual[0] == 150 {
		expected = []int{150, 100, 150, 100, 150, 100, 150, 100, 150, 100}
	} else {
		expected = []int{100, 150, 100, 150, 100, 150, 100, 150, 100, 150}
	}
	assert.Equal(t, expected, actual)

}

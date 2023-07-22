package rx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/pkg/rx"
)

func TestMap(t *testing.T) {
	s := rx.NewSubject[int]()

	var actual []string
	rx.Map[int, string](s, func(i int) string {
		return fmt.Sprint(i)
	}).Subscribe(rx.NewObserver(func(i string) {
		actual = append(actual, i)
	}, nil, nil))

	for _, i := range []int{3, 1, 4, 1, 5} {
		s.Next(i)
	}
	assert.Equal(t, []string{"3", "1", "4", "1", "5"}, actual)
}

func TestReduce(t *testing.T) {
	result := []int{0}
	fc := rx.ToConnectable(rx.From[int](1, 2, 3))

	rx.Reduce[int, int](fc, func(u, t int) int {
		return u + t
	}, 0).Subscribe(rx.NewObserver(func(v int) {
		result[0] = v
	}, nil, nil))

	fc.Connect()

	assert.Equal(t, 6, result[0])

}

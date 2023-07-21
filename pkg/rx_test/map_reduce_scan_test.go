package rx_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

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

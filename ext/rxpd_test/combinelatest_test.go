package rxpd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/ext/rxpd"
	"github.com/philippseith/gorx/pkg/rx"
)

func TestCombineLatest(t *testing.T) {
	var actual []any
	ab := rx.ToConnectable[any](rx.From[any]('A', 'B'))
	oneTwoThree := rx.ToConnectable[any](rx.From[any](1, 2, 3))

	abProp := rxpd.ToProperty[any](ab)
	oneTwoThreeProp := rxpd.ToProperty[any](oneTwoThree)

	rxpd.CombineLatest[[]any](func(n ...any) []any {
		var nn []any
		nn = append(nn, n...)
		return nn
	}, abProp, oneTwoThreeProp).Subscribe(rx.NewObserver[[]any](func(next []any) {
		actual = next
	}, nil, nil))

	ab.Connect()
	oneTwoThree.Connect()

	assert.Equal(t, []any{'B', 3}, actual)
}

//func TestCombineLatestTicker(t *testing.T) {
//	for i := 0; i < 1; i++ {
//		// 149 and 199 are prime
//		t1 := rx.NewTicker(100*time.Millisecond, 1*time.Microsecond)
//		ch1 := make(chan int, 1)
//		ch2 := make(chan int, 1)
//		j := 0
//		t1.Subscribe(rx.OnNext[time.Time](func(time.Time) {
//			j++
//			//fmt.Println(j)
//			if j == 6 {
//				close(ch1)
//				close(ch2)
//			}
//			if j >= 6 {
//				return
//			}
//			ch1 <- j
//			if j%2 == 0 {
//				ch2 <- j
//			}
//		}))
//		tt1 := rx.Scan[int, int](rx.FromChan[int](ch1), func(i int, k int) int {
//			//fmt.Println("A", i, " ", k)
//			return i + 1
//		}, 0)
//		tt2 := rx.Scan[int, int](rx.FromChan[int](ch2), func(i int, k int) int {
//			//fmt.Println("B", i, " ", k)
//			return i + 1
//		}, 0)
//
//		sl := <-rx.ToSlice[[]any](rx.CombineLatest[[]any](func(ts ...any) []any { return ts }, rx.ToAny[int](tt1), rx.ToAny[int](tt2)))
//
//		for _, s := range sl {
//			fmt.Println(s)
//		}
//		assert.Equal(t, [][]any{{2, 1}, {2, 2}, {3, 2}, {4, 2}, {5, 2}}, sl)
//	}
//}

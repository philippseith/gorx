package rxlive_test

import (
	"time"

	"github.com/philippseith/gorx/pkg/rx"
)

var rxModel rx.Observable[[]Row]

func BackEnd() rx.Subscribable[[]Row] {
	if rxModel == nil {
		rxModel = rx.Map[time.Time, []Row](rx.NewTicker(0, time.Second), func(t time.Time) []Row {
			if t.Second()%2 == 0 {
				return []Row{
					{Id: "BBB", States: []string{"456"}},
					{Id: "CCC"},
				}
			} else {
				return []Row{
					{Id: "AAA", States: []string{"123"}},
					{Id: "BBB", States: []string{"456"}},
					{Id: "CCC"},
				}
			}
		}).ShareReplay(rx.MaxBufferSize(1))
	}
	return rxModel
}

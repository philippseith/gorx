package rxpd

import (
	"time"

	"github.com/philippseith/gorx/pkg/rx"
)

type Property[T any] interface {
	rx.Observable[T]

	SetInterval(time.Duration)
	Interval() time.Duration

	Read() <-chan rx.Result[T]
	Write(T) <-chan error
}

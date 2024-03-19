package rxpd

import "github.com/philippseith/gorx/pkg/rx"

// Concat creates an output Property which sequentially emits all values from
// the first given Subscribable and then moves on to the next.
func Concat[T any](sources ...Subscribable[T]) Property[T] {
	rxSources := make([]rx.Subscribable[T], 0, len(sources))
	for _, source := range sources {
		rxSources = append(rxSources, source)
	}
	c := rx.Concat[T](rxSources...)
	return ToProperty[T](c)
}

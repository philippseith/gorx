package rxpd

import (
	"sync"

	"github.com/philippseith/gorx/pkg/rx"
)

// PropertySource creates and stores Property objects by IDs that can be reused.
type PropertySource[ID comparable, T any] struct {
	mx             sync.RWMutex
	properties     map[ID]Property[T]
	createProperty func(ID) Property[T]
}

// NewPropertySource creates a PropertySource which creates Property objects for IDs with the createProperty function.
func NewPropertySource[ID comparable, T any](createProperty func(ID) Property[T]) *PropertySource[ID, T] {
	return &PropertySource[ID, T]{
		properties:     map[ID]Property[T]{},
		createProperty: createProperty,
	}
}

// GetProperty returns the Property object for an ID. The Property has ShareReplay behavior with MaxBufferSize 1.
// This means each new subscriber will get the last value of the Property, if there is any.
func (ps *PropertySource[ID, T]) GetProperty(id ID) Property[T] {

	if property := func() Property[T] {
		ps.mx.RLock()
		defer ps.mx.RUnlock()

		return ps.properties[id]
	}(); property != nil {
		return property
	}

	property := ps.createProperty(id).ShareReplay(rx.MaxBufferSize(1), rx.RefCount(true))

	func() {
		ps.mx.Lock()
		defer ps.mx.Unlock()

		ps.properties[id] = property
	}()

	return property
}

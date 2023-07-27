package rx

import "sync"

type Subscription interface {
	Unsubscribe()
	AddSubscription(Subscription)
	AddTearDownLogic(func())
}

func NewSubscription(unsubscribe func()) Subscription {
	s := &subscription{}

	s.mx.Lock()
	defer s.mx.Unlock()

	s.tdls = append(s.tdls, unsubscribe)
	return s
}

type subscription struct {
	tdls []func()
	mx   sync.Mutex
}

func (s *subscription) AddSubscription(su Subscription) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.tdls = append(s.tdls, su.Unsubscribe)
}

func (s *subscription) AddTearDownLogic(tld func()) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.tdls = append(s.tdls, tld)
}

func (s *subscription) Unsubscribe() {
	s.mx.Lock()
	defer s.mx.Unlock()

	for _, tld := range s.tdls {
		tld()
	}
	s.tdls = nil
}

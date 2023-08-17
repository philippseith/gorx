package rx

import "sync"

type Subscription interface {
	Unsubscribe()
	AddSubscription(Subscription) Subscription
	AddTearDownLogic(func()) Subscription
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

func (s *subscription) AddSubscription(su Subscription) Subscription {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.tdls = append(s.tdls, su.Unsubscribe)
	return s
}

func (s *subscription) AddTearDownLogic(tld func()) Subscription {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.tdls = append(s.tdls, tld)
	return s
}

func (s *subscription) Unsubscribe() {
	s.mx.Lock()
	defer s.mx.Unlock()

	for _, tld := range s.tdls {
		tld()
	}
	s.tdls = nil
}

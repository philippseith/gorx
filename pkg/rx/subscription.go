package rx

type Subscription interface {
	Unsubscribe()
}

func NewSubscription(unsubscribe func()) Subscription {
	return &subscription{u: unsubscribe}
}

type subscription struct {
	u func()
}

func (s *subscription) Unsubscribe() {
	s.u()
}

package rx

type Subscription interface {
	Unsubscribe()
}

type subscription struct {
	u func()
}

func (s *subscription) Unsubscribe() {
	s.u()
}

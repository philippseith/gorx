package rx

func Merge[T any](sources ...Subscribable[T]) Observable[T] {
	m := &merge[T]{
		Operator: Operator[T, T]{t2u: func(t T) T { return t }},
	}
	m.completed = make([]bool, len(sources))
	for i, source := range sources {
		ci := i
		m.subs = append(m.subs, source.Subscribe(NewObserver[T](m.Next, m.Error, func() {
			func() {
				m.mxState.Lock()
				defer m.mxState.Unlock()

				m.completed[ci] = true
			}()
			if func() bool {
				m.mxState.RLock()
				defer m.mxState.RUnlock()

				for _, completed := range m.completed {
					if !completed {
						return false
					}
				}
				return true
			}() {
				m.Complete()
			}
		})))
	}
	return ToObservable[T](m)
}

type merge[T any] struct {
	Operator[T, T]
	subs      []Subscription
	completed []bool
}

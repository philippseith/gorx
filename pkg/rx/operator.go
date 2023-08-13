package rx

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

type Operator[T any, U any] struct {
	outObserver        Observer[U]
	sourceSubscription Subscription
	mx                 sync.RWMutex
	t2u                func(T) U
}

func (op *Operator[T, U]) Next(t T) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in %T.Next(%v): %v.\n%s", op, t, r, string(debug.Stack()))
			if o := op.getObserver(); o != nil {
				o.Error(err)
			} else {
				log.Print(err)
			}
		}
	}()

	if o := op.getObserver(); o != nil {
		o.Next(op.t2u(t))
	}
}

func (op *Operator[T, U]) Error(err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in %T.Error(%v): %v.\n%s", op, err, r, string(debug.Stack()))
		}
	}()

	if o := op.getObserver(); o != nil {
		o.Error(err)
	}
}

func (op *Operator[T, U]) Complete() {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in %T.Complete(): %v\n%s", op, r, string(debug.Stack()))
			if o := op.getObserver(); o != nil {
				o.Error(err)
			} else {
				log.Print(err)
			}
		}
	}()

	if o := op.getObserver(); o != nil {
		o.Complete()
	}
}

func (op *Operator[T, U]) getObserver() Observer[U] {
	op.mx.RLock()
	defer op.mx.RUnlock()

	return op.outObserver
}

func (op *Operator[T, U]) Subscribe(o Observer[U]) Subscription {
	op.mx.Lock()
	defer op.mx.Unlock()
	op.outObserver = o

	return NewSubscription(func() {
		op.mx.RLock()
		defer op.mx.RUnlock()

		if op.sourceSubscription != nil {
			op.sourceSubscription.Unsubscribe()
		}
	})
}

func (op *Operator[T, U]) SubscribeToSource(ob Observer[T], s Subscribable[T]) {
	op.sourceSubscription = s.Subscribe(ob)
}

package rx

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

type Operator[T any, U any] struct {
	onSubscribe        func() Subscription
	sourceSubscription Subscription
	outObserver        Observer[U]
	mxState            sync.RWMutex
	mxEvents           sync.Mutex
	t2u                func(T) U
}

func (op *Operator[T, U]) Next(t T) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in %T.Next(%v): %v.\n%s", op, t, r, string(debug.Stack()))
			if o := op.observer(); o != nil {
				o.Error(err)
			} else {
				log.Print(err)
			}
		}
	}()

	if o := op.observer(); o != nil {
		func() {
			op.mxEvents.Lock()
			defer op.mxEvents.Unlock()

			o.Next(op.t2u(t))
		}()
	}
}

func (op *Operator[T, U]) Error(err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in %T.Error(%v): %v.\n%s", op, err, r, string(debug.Stack()))
		}
	}()

	if o := op.observer(); o != nil {
		func() {
			op.mxEvents.Lock()
			defer op.mxEvents.Unlock()

			o.Error(err)
		}()
	}
}

func (op *Operator[T, U]) Complete() {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in %T.Complete(): %v\n%s", op, r, string(debug.Stack()))
			if o := op.observer(); o != nil {
				o.Error(err)
			} else {
				log.Print(err)
			}
		}
	}()

	if o := op.observer(); o != nil {
		func() {
			op.mxEvents.Lock()
			defer op.mxEvents.Unlock()

			o.Complete()
		}()
	}
}

func (op *Operator[T, U]) observer() Observer[U] {
	op.mxState.RLock()
	defer op.mxState.RUnlock()

	return op.outObserver
}

// Subscribe connects an Observer.
// The Operator does not call Subscribe of its source until its own Subscribe method is called.
func (op *Operator[T, U]) Subscribe(o Observer[U]) Subscription {
	// If anything is in the operator chain that directly calls Next
	// this would deadlock if we simply lock everything with mxState.
	// So we lock more fine-grained
	if onSubscribe := func() func() Subscription {
		op.mxState.Lock()
		defer op.mxState.Unlock()

		op.outObserver = o

		return op.onSubscribe

	}(); onSubscribe != nil {

		sub := op.onSubscribe()

		func() {
			op.mxState.Lock()
			defer op.mxState.Unlock()

			op.sourceSubscription = sub
		}()
	}

	return NewSubscription(func() {
		op.mxState.RLock()
		defer op.mxState.RUnlock()

		if op.sourceSubscription != nil {
			op.sourceSubscription.Unsubscribe()
		}
	})
}

// prepareSubscribe sets the onSubscribe method which is called in Subscribe.
// Operators should subscribe their source only in the onSubscribe method.
func (op *Operator[T, U]) prepareSubscribe(onSubscribe func() Subscription) {
	op.onSubscribe = onSubscribe
}

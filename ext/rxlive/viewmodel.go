package rxlive

import (
	"context"
	"log"
	"sync"

	"github.com/jfyne/live"

	"github.com/philippseith/gorx/pkg/rx"
)

type viewModel[T any] struct {
	data T
	sub  rx.Subscription
	mx   sync.Mutex
}

func (v *viewModel[T]) Data() T {
	return v.data
}

func (v *viewModel[T]) subscribe(sub rx.Subscription) {
	v.mx.Lock()
	defer v.mx.Unlock()

	v.sub = sub
}

func (v *viewModel[T]) unsubscribe() {
	v.mx.Lock()
	defer v.mx.Unlock()

	if v.sub != nil {
		v.sub.Unsubscribe()
		v.sub = nil
	}
}

func newViewModel[T any](ctx context.Context, socket live.Socket, model rx.Subscribable[T]) (any, error) {
	vm, ok := socket.Assigns().(*viewModel[T])
	if !ok {
		vm = &viewModel[T]{}
		if socket.Connected() {
			// We come here when the page is already open or in the second Mount with the Websocket
			if err := socket.Send("reload", "nil"); err != nil {
				log.Print(err)
			}
			var mx sync.Mutex
			vm.subscribe(model.Subscribe(rx.OnNext(func(data T) {
				// The order of the update events needs to be preserved,
				// as socket.Self seems not to be reentrancy-safe
				mx.Lock()
				defer mx.Unlock()

				if err := socket.Self(ctx, "vmChanged", data); err != nil {
					log.Print(err)
				}
			})))
			go func() {
				<-ctx.Done()
				vm.unsubscribe()
			}()
		}
	}
	return vm, nil
}

func tearDownViewModel[T any](socket live.Socket) error {
	if vm, ok := socket.Assigns().(*viewModel[T]); ok {
		vm.unsubscribe()
	}
	return nil
}

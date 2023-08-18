package rxlive

import (
	"context"
	"log"
	"sync"

	"github.com/jfyne/live"

	"github.com/philippseith/gorx/pkg/rx"
)

type viewModel[T any] struct {
	data   T
	sub    rx.Subscription
	ctx    context.Context
	socket live.Socket
}

func (v *viewModel[T]) Data() T {
	return v.data
}

func newViewModel[T any](ctx context.Context, socket live.Socket, model rx.Subscribable[T]) (any, error) {
	vm, ok := socket.Assigns().(*viewModel[T])
	if !ok {
		vm = &viewModel[T]{
			ctx:    ctx,
			socket: socket,
		}
		if socket.Connected() {
			// We come here when the page is already open or in the second Mount with the Websocket
			if err := socket.Send("reload", "nil"); err != nil {
				log.Print(err)
			}
			var mx sync.Mutex
			vm.sub = model.Subscribe(rx.OnNext(func(data T) {
				// The order of the update events needs to be preserved,
				// as socket.Self seems not to be reentrancy-safe
				mx.Lock()
				defer mx.Unlock()

				if err := vm.socket.Self(vm.ctx, "vmChanged", data); err != nil {
					log.Print(err)
				}
			}))
			// TODO When is ctx.Done?
			go func() {
				<-ctx.Done()
				vm.sub.Unsubscribe()
			}()
		}
	}
	return vm, nil
}

func tearDownViewModel[T any](socket live.Socket) error {
	vm, ok := socket.Assigns().(*viewModel[T])
	if ok && vm.sub != nil {
		vm.sub.Unsubscribe()
	}
	return nil
}

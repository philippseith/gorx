package rxlive

import (
	"context"
	"log"

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
			vm.sub = model.Subscribe(rx.OnNext(func(data T) {
				if err := vm.socket.Self(vm.ctx, "vmChanged", data); err != nil {
					log.Print(err)
				}
			}))
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

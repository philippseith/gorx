package rxlive

import (
	"context"
	_ "embed"
	"html/template"
	"log"
	"net/http"

	"github.com/jfyne/live"
	"github.com/philippseith/gorx/pkg/rx"
)

// NewHandler creates an http.Handler wrapping a live.Handler which fills the
// template with data of type T. To handle temporary connection losses on the
// client side, the live-hook-reload template can be used.
func NewHandler[T any](tmpl *template.Template, model rx.Subscribable[T]) http.Handler {

	var err error
	tmpl, err = tmpl.Parse(liveHookReload)
	if err != nil {
		log.Fatal(err)
	}

	h := live.NewHandler(live.WithTemplateRenderer(tmpl))

	h.HandleMount(func(ctx context.Context, socket live.Socket) (any, error) {
		return newViewModel[T](ctx, socket, model)
	})

	h.HandleUnmount(func(socket live.Socket) error {
		return tearDownViewModel[T](socket)
	})

	h.HandleSelf("vmChanged", func(ctx context.Context, socket live.Socket, d any) (any, error) {
		vm, okVm := socket.Assigns().(*viewModel[T])
		data, okData := d.(T)
		if okVm && okData {
			vm.data = data
		}
		return vm, nil
	})

	return live.NewHttpHandler(live.NewCookieStore("session-name", []byte("weak-secret")), h)
}

//go:embed reload.html
var liveHookReload string

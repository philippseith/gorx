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
// template with data of type T from the model parameter. The sessionAuthKey is
// used to build session cookies, by which the client of a live session can be
// identified. To handle temporary connection losses on the client side, the
// live-hook-reload template can be included with {{ template live-hook-reload
// }} in your template.
func NewHandler[T any](tmpl *template.Template, model rx.Subscribable[T], sessionAuthKey []byte) http.Handler {

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

	return live.NewHttpHandler(live.NewCookieStore("session-name", sessionAuthKey), h)
}

//go:embed reload.html
var liveHookReload string

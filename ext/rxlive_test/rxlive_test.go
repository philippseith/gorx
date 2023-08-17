package rxlive

import (
	"embed"
	"html/template"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jfyne/live"
	"github.com/stretchr/testify/assert"

	"github.com/philippseith/gorx/ext/rxlive"
	"github.com/philippseith/gorx/pkg/rx"
)

func TestRxLive(t *testing.T) {
	tmpl, err := template.ParseFS(FS, "root.html", "overview.html")
	assert.NoError(t, err)

	http.Handle("/", rxlive.NewHandler[[]Row](tmpl, BackEnd()))
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	assert.NoError(t, http.ListenAndServe(":8086", nil))
}

func TestRxLiveGin(t *testing.T) {
	router := gin.Default()

	tmpl, err := template.ParseFS(FS, "root.html", "overview.html")
	assert.NoError(t, err)

	router.GET("/", gin.WrapH(rxlive.NewHandler[[]Row](tmpl, BackEnd())))
	router.GET("/live.js", gin.WrapH(live.Javascript{}))
	router.GET("/auto.js.map", gin.WrapH(live.JavascriptMap{}))
	assert.NoError(t, router.Run(":8087"))
}

type Row struct {
	Id     string
	States []string
}

var rxModel rx.Observable[[]Row]

func BackEnd() rx.Subscribable[[]Row] {
	if rxModel == nil {
		rxModel = rx.Map[time.Time, []Row](rx.NewTicker(0, time.Second), func(t time.Time) []Row {
			if t.Second()%2 == 0 {
				return []Row{
					{Id: "BBB", States: []string{"456"}},
					{Id: "CCC"},
				}
			} else {
				return []Row{
					{Id: "AAA", States: []string{"123"}},
					{Id: "BBB", States: []string{"456"}},
					{Id: "CCC"},
				}
			}
		}).ShareReplay(rx.MaxBufferSize(1))
	}
	return rxModel
}

//go:embed *.html
var FS embed.FS

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

	http.Handle("/", rxlive.NewHandler[[]Row](tmpl, Model(), []byte("____os.Getenv('SESSION_KEY')____")))
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	assert.NoError(t, http.ListenAndServe(":8086", nil))
}

func TestRxLiveGin(t *testing.T) {
	router := gin.Default()

	tmpl, err := template.ParseFS(FS, "root.html", "overview.html")
	assert.NoError(t, err)

	router.GET("/", gin.WrapH(
		rxlive.NewHandler[[]Row](tmpl, Model(), []byte("____os.Getenv('SESSION_KEY')____"))))
	router.GET("/live.js", gin.WrapH(live.Javascript{}))
	router.GET("/auto.js.map", gin.WrapH(live.JavascriptMap{}))
	assert.NoError(t, router.Run(":8087"))
}

type Row struct {
	Id   string
	Time string
}

func Model() rx.Subscribable[[]Row] {
	t1 := rx.NewTicker(0, 1000*time.Millisecond)
	t2 := rx.NewTicker(300*time.Millisecond, 10*time.Millisecond)
	return rx.CombineLatest2(func(t1, t2 time.Time) []Row {
		if t1.Second()%2 == 0 {
			return []Row{
				{Id: "BBB", Time: t2.Format(time.StampMilli)},
			}
		} else {
			return []Row{
				{Id: "AAA",
					Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.StampMilli)},
				{Id: "BBB", Time: t2.Format(time.StampMilli)},
			}
		}
	}, t1, t2).DistinctUntilChanged(nil).ShareReplay(rx.MaxBufferSize(1))
}

//go:embed *.html
var FS embed.FS

package rxlive_test

import (
	"html/template"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jfyne/live"
	"github.com/philippseith/gorx/ext/rxlive"
	"github.com/stretchr/testify/assert"
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

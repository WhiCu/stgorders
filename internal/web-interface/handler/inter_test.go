package handler

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WhiCu/stgorders/pkg/logger"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFS(t *testing.T) {
	Convey("FS returns index.html", t, func() {
		html, err := fs.Sub(templates, "templates")
		So(err, ShouldBeNil)
		http_html := http.FS(html)
		f, err := http_html.Open("index.html")

		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
	})
}

func TestInter(t *testing.T) {
	Convey("Inter serves index.html", t, func() {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		r := gin.New()
		h := NewHandler(&mockService{}, logger.NewNOPSlog())
		r.GET("/", h.Inter)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, http.StatusOK)
		So(w.Body.String(), ShouldContainSubstring, "Поиск заказа")
	})
}

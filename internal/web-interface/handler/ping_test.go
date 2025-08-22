package handler

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPing(t *testing.T) {
	Convey("Ping returns pong", t, func() {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		r := gin.New()
		h := &Handler{log: slog.Default()}
		r.GET("/ping", h.Ping)

		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		r.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusOK)
		So(w.Body.String(), ShouldEqual, "pong")
	})
}

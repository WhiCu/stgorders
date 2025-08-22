package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/WhiCu/stgorders/db/model"
	"github.com/WhiCu/stgorders/pkg/logger"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

type mockService struct{}

func (m *mockService) GetJsonOrderByUID(ctx context.Context, orderUID string) (*model.JsonOrder, error) {
	if orderUID == "ok" {
		return &model.JsonOrder{Order: model.Order{OrderUID: "ok", DateCreated: time.Now()}}, nil
	}
	return nil, errors.New("not found")
}

func TestOrder_OK(t *testing.T) {
	Convey("Order returns 200 with JSON", t, func() {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		r := gin.New()
		h := NewHandler(&mockService{}, logger.NewNOPSlog())
		r.GET("/:orderUID", h.Order)

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		r.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

func TestOrder_BadURI(t *testing.T) {
	Convey("Router returns 404 when route param missing", t, func() {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		r := gin.New()
		h := NewHandler(&mockService{}, logger.NewNOPSlog())
		r.GET("/:orderUID", h.Order)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

func TestOrder_NotFound(t *testing.T) {
	Convey("Order returns 404 when service returns error", t, func() {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		r := gin.New()
		h := NewHandler(&mockService{}, logger.NewNOPSlog())
		r.GET("/:orderUID", h.Order)

		req := httptest.NewRequest(http.MethodGet, "/missing", nil)
		r.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

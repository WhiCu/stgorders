package service

import (
	"context"
	"errors"
	"testing"

	"github.com/WhiCu/stgorders/db/model"
	"github.com/WhiCu/stgorders/pkg/logger"
	. "github.com/smartystreets/goconvey/convey"
)

type mockStorage struct {
	json model.JsonOrder
	err  error
}

func (m *mockStorage) GetJsonOrderByUID(ctx context.Context, orderUID string) (*model.JsonOrder, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &m.json, nil
}

func TestService_GetJsonOrderByUID(t *testing.T) {
	Convey("Service returns json order and logs on success", t, func() {
		m := &mockStorage{json: model.JsonOrder{Order: model.Order{OrderUID: "u"}}}
		s := NewService(m, logger.NewNOPSlog())
		jo, err := s.GetJsonOrderByUID(context.Background(), "u")
		So(err, ShouldBeNil)
		So(jo, ShouldNotBeNil)
		So(jo.OrderUID, ShouldEqual, "u")
	})

	Convey("Service propagates error", t, func() {
		m := &mockStorage{err: errors.New("err")}
		s := NewService(m, logger.NewNOPSlog())
		jo, err := s.GetJsonOrderByUID(context.Background(), "u")
		So(jo, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}

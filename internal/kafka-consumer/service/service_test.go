package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/WhiCu/stgorders/db/model"
	"github.com/WhiCu/stgorders/pkg/logger"
	. "github.com/smartystreets/goconvey/convey"
)

type mockSaver struct {
	err error
}

func (m *mockSaver) Save(ctx context.Context, order model.JsonOrder) error { return m.err }
func (m *mockSaver) Close() error                                          { return nil }

func TestService_Serve_Success(t *testing.T) {
	Convey("Serve unmarshals, validates and saves order", t, func() {
		s := NewService(&mockSaver{}, logger.NewNOPSlog())
		payload, _ := json.Marshal(model.JsonOrder{
			Order:    model.Order{OrderUID: "1", TrackNumber: "t", Entry: "e", Locale: "ru", CustomerID: "c", DeliveryService: "d", Shardkey: "s", SmID: 1, DateCreated: time.Now(), OofShard: "o"},
			Delivery: model.Delivery{Name: "n", Phone: "+1234567890", Zip: "z", City: "c", Address: "a", Region: "r", Email: "e@e.e"},
			Payment:  model.Payment{Transaction: "1", Currency: "RUB", Provider: "p", Amount: 1, PaymentDt: 1, Bank: "b", DeliveryCost: 0, GoodsTotal: 1},
			Items:    []model.Item{{ChrtID: 1, TrackNumber: "t", Price: 1, Rid: "r", Name: "n", Sale: 0, Size: "s", TotalPrice: 1, NmID: 1, Brand: "b", Status: 1}},
		})

		err := s.Serve(context.Background(), payload)
		So(err, ShouldBeNil)
	})
}

func TestService_Serve_InvalidJSON(t *testing.T) {
	Convey("Serve returns error on invalid JSON", t, func() {
		s := NewService(&mockSaver{}, logger.NewNOPSlog())
		err := s.Serve(context.Background(), []byte("{"))
		So(err, ShouldNotBeNil)
	})
}

func TestService_Serve_ValidationError(t *testing.T) {
	Convey("Serve returns error on validation failure", t, func() {
		s := NewService(&mockSaver{}, logger.NewNOPSlog())
		payload, _ := json.Marshal(model.JsonOrder{})
		err := s.Serve(context.Background(), payload)
		So(err, ShouldNotBeNil)
	})
}

func TestService_Serve_SaveError(t *testing.T) {
	Convey("Serve returns error when storage fails", t, func() {
		s := NewService(&mockSaver{err: errors.New("err")}, logger.NewNOPSlog())
		payload, _ := json.Marshal(model.JsonOrder{
			Order:    model.Order{OrderUID: "1", TrackNumber: "t", Entry: "e", Locale: "ru", CustomerID: "c", DeliveryService: "d", Shardkey: "s", SmID: 1, DateCreated: time.Now(), OofShard: "o"},
			Delivery: model.Delivery{Name: "n", Phone: "+1234567890", Zip: "z", City: "c", Address: "a", Region: "r", Email: "e@e.e"},
			Payment:  model.Payment{Transaction: "1", Currency: "RUB", Provider: "p", Amount: 1, PaymentDt: 1, Bank: "b", DeliveryCost: 0, GoodsTotal: 1},
			Items:    []model.Item{{ChrtID: 1, TrackNumber: "t", Price: 1, Rid: "r", Name: "n", Sale: 0, Size: "s", TotalPrice: 1, NmID: 1, Brand: "b", Status: 1}},
		})
		err := s.Serve(context.Background(), payload)
		So(err, ShouldNotBeNil)
	})
}

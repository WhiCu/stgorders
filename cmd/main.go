package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/WhiCu/stgorders/internal/config"
	"github.com/WhiCu/stgorders/internal/kafka-consumer/client"
	"github.com/WhiCu/stgorders/internal/kafka-consumer/service"
	"github.com/WhiCu/stgorders/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type jsonOrder struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int32     `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int32  `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int32  `json:"delivery_cost"`
	GoodsTotal   int32  `json:"goods_total"`
	CustomFee    int32  `json:"custom_fee"`
}

type Item struct {
	ChrtID      int64  `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int32  `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int32  `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int32  `json:"total_price"`
	NmID        int64  `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int32  `json:"status"`
}

func main() {
	ctx := context.Background()
	cfg := config.MustLoadWithDefault("./config/config.yaml")
	log := slog.New(logger.MustInitLogger("debug"))
	log.Info("config loaded", slog.String("DSN", cfg.Storage.DSN()))
	p, err := pgxpool.New(ctx, cfg.Storage.DSN())
	if err != nil {
		panic(err)
	}
	defer p.Close()
	s := client.NewStorage(p)
	svc := service.NewService(s, log)
	err = svc.Serve(ctx, []byte(`{
    "order_uid": "b563feb7b2b84b6test1",
    "track_number": "WBILMTESTTRACK1",
    "entry": "WBIL",
    "delivery": {
      "name": "Alice Smith",
      "phone": "+9720000001",
      "zip": "2639801",
      "city": "Kiryat Mozkin",
      "address": "Ploshad Mira 1",
      "region": "Kraiot",
      "email": "alice1@gmail.com"
    },
    "payment": {
      "transaction": "b563feb7b2b84b6test1",
      "request_id": "",
      "currency": "USD",
      "provider": "wbpay",
      "amount": 1200,
      "payment_dt": 1637907727,
      "bank": "alpha",
      "delivery_cost": 800,
      "goods_total": 400,
      "custom_fee": 0
    },
    "items": [
      {
        "chrt_id": 9934931,
        "track_number": "WBILMTESTTRACK1",
        "price": 400,
        "rid": "ab4219087a764ae01",
        "name": "Mascaras",
        "sale": 20,
        "size": "0",
        "total_price": 320,
        "nm_id": 2389213,
        "brand": "Vivienne Sabo",
        "status": 202
      }
    ],
    "locale": "en",
    "internal_signature": "",
    "customer_id": "test1",
    "delivery_service": "meest",
    "shardkey": "1",
    "sm_id": 1,
    "date_created": "2021-11-26T06:22:19Z",
    "oof_shard": "1"
  }`))

	if err != nil {
		panic(err)
	}

	orders, err := s.ListOrders(ctx)
	if err != nil {
		panic(err)
	}

	for i, order := range orders {
		fmt.Println(i, order)
	}

}

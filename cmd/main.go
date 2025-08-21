package main

import (
	"fmt"

	"github.com/WhiCu/stgorders/internal/config"
)

var data = `{
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
  }`

func main() {
	cfg := config.MustLoadWithDefault("./config/config.yaml")
	fmt.Println("cfg: ", cfg.Cache)
	// // var jo model.JsonOrder
	// // if err := json.Unmarshal([]byte(data), jo); err != nil {
	// // 	log.Fatal(err)
	// // }

	// log := slog.New(logger.MustInitLogger("debug"))
	// c := cache.NewLRUCache[string, model.JsonOrder](cfg.Cache.Size, log.WithGroup("cache"))
	// p, err := pgxpool.New(context.Background(), cfg.Storage.DSN())
	// if err != nil {
	// 	panic(err)
	// }
	// stg := storage.NewStorage(p, log.WithGroup("storage"))
	// err = stg.InitCache(context.Background(), c)
	// if err != nil {
	// 	log.Error("could not init cache", slog.String("ERR", err.Error()))
	// }

	// // jsonOrder, err := client.NewStorage(stg, c, log).GetJsonOrderByUID(context.Background(), "b563feb7b2b84b6test9")
	// // if err != nil {
	// // 	log.Error("could not get order", slog.String("ERR", err.Error()))
	// // }
	// // fmt.Println("jo", jsonOrder)
	// s := client.NewStorage(stg, c, log.WithGroup("storageAdapter").With("orderUID", "b563feb7b2b84b6test9"))
	// s.GetJsonOrderByUID(context.Background(), "b563feb7b2b84b6test9")
	// s.Close()
}

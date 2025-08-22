package model

import (
	"time"

	"github.com/WhiCu/stgorders/db/pg"
)

type JsonOrder struct {
	Order
	Delivery Delivery `json:"delivery" validate:"required"`
	Payment  Payment  `json:"payment" validate:"required"`
	Items    []Item   `json:"items" validate:"required,min=1,dive"`
}

type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required"`
	TrackNumber       string    `json:"track_number" validate:"required"`
	Entry             string    `json:"entry" validate:"required"`
	Locale            string    `json:"locale" validate:"required"`
	InternalSignature string    `json:"internal_signature" validate:"omitempty"`
	CustomerID        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	Shardkey          string    `json:"shardkey" validate:"required"`
	SmID              int32     `json:"sm_id" validate:"required"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required,e164"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`
	RequestID    string `json:"request_id" validate:"omitempty"`
	Currency     string `json:"currency" validate:"required"`
	Provider     string `json:"provider" validate:"required"`
	Amount       int32  `json:"amount" validate:"required"`
	PaymentDt    int64  `json:"payment_dt" validate:"required"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost int32  `json:"delivery_cost" validate:"omitempty,gte=0"`
	GoodsTotal   int32  `json:"goods_total" validate:"required"`
	CustomFee    int32  `json:"custom_fee" validate:"omitempty"`
}

type Item struct {
	ChrtID      int64  `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int32  `json:"price" validate:"required"`
	Rid         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        int32  `json:"sale" validate:"omitempty,gte=0,lte=100"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  int32  `json:"total_price" validate:"required,gte=0"`
	NmID        int64  `json:"nm_id" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	Status      int32  `json:"status" validate:"required"`
}

func ModelJSONOrder(o pg.Order, d pg.Delivery, p pg.Payment, items []pg.Item) JsonOrder {
	// преобразуем items
	mItems := make([]Item, len(items))
	for i, it := range items {
		mItems[i] = Item{
			ChrtID:      it.ChrtID,
			TrackNumber: it.TrackNumber,
			Price:       it.Price,
			Rid:         it.Rid,
			Name:        it.Name,
			Sale:        it.Sale,
			Size:        it.Size,
			TotalPrice:  it.TotalPrice,
			NmID:        it.NmID,
			Brand:       it.Brand,
			Status:      it.Status,
		}
	}

	return JsonOrder{
		Order: Order{
			OrderUID:          o.OrderUid,
			TrackNumber:       o.TrackNumber,
			Entry:             o.Entry,
			Locale:            o.Locale,
			InternalSignature: o.InternalSignature,
			CustomerID:        o.CustomerID,
			DeliveryService:   o.DeliveryService,
			Shardkey:          o.Shardkey,
			SmID:              o.SmID,
			DateCreated:       o.DateCreated,
			OofShard:          o.OofShard,
		},
		Delivery: Delivery{
			Name:    d.Name,
			Phone:   d.Phone,
			Zip:     d.Zip,
			City:    d.City,
			Address: d.Address,
			Region:  d.Region,
			Email:   d.Email,
		},
		Payment: Payment{
			Transaction:  p.Transaction,
			RequestID:    p.RequestID,
			Currency:     p.Currency,
			Provider:     p.Provider,
			Amount:       p.Amount,
			PaymentDt:    p.PaymentDt,
			Bank:         p.Bank,
			DeliveryCost: p.DeliveryCost,
			GoodsTotal:   p.GoodsTotal,
			CustomFee:    p.CustomFee,
		},
		Items: mItems,
	}
}

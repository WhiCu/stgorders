package service

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/WhiCu/stgorders/db/pg"
	"github.com/WhiCu/stgorders/internal/kafka-consumer/client"
	"github.com/jackc/pgx/v5/pgtype"
)

type Storage interface {
}

type Service struct {
	storage *client.Storage

	log *slog.Logger
}

func NewService(storage *client.Storage, log *slog.Logger) *Service {
	return &Service{
		storage: storage,
		log:     log,
	}
}

func (s *Service) Serve(ctx context.Context, data []byte) (err error) {
	db, rollback, commit, err := s.storage.WithTx(ctx)
	defer rollback(ctx)

	var jo jsonOrder
	if err = json.Unmarshal(data, &jo); err != nil {
		panic(err)
	}

	orderID, err := db.CreateOrder(ctx, pg.CreateOrderParams{
		OrderUid:          jo.OrderUID,
		TrackNumber:       jo.TrackNumber,
		Entry:             jo.Entry,
		Locale:            jo.Locale,
		InternalSignature: jo.InternalSignature,
		CustomerID:        jo.CustomerID,
		DeliveryService:   jo.DeliveryService,
		Shardkey:          jo.Shardkey,
		SmID:              jo.SmID,
		DateCreated: pgtype.Timestamp{
			Time:  jo.DateCreated,
			Valid: true,
		},
		OofShard: jo.OofShard,
	})
	if err != nil {
		panic(err)
	}

	_, err = db.CreateDelivery(ctx, pg.CreateDeliveryParams{
		OrderID: orderID,
		Name:    jo.Delivery.Name,
		Phone:   jo.Delivery.Phone,
		Zip:     jo.Delivery.Zip,
		City:    jo.Delivery.City,
		Address: jo.Delivery.Address,
		Region:  jo.Delivery.Region,
		Email:   jo.Delivery.Email,
	})
	if err != nil {
		panic(err)
	}

	_, err = db.CreatePayment(ctx, pg.CreatePaymentParams{
		OrderID:      orderID,
		Transaction:  jo.Payment.Transaction,
		RequestID:    jo.Payment.RequestID,
		Currency:     jo.Payment.Currency,
		Provider:     jo.Payment.Provider,
		Amount:       jo.Payment.Amount,
		PaymentDt:    jo.Payment.PaymentDt,
		Bank:         jo.Payment.Bank,
		DeliveryCost: jo.Payment.DeliveryCost,
		GoodsTotal:   jo.Payment.GoodsTotal,
		CustomFee:    jo.Payment.CustomFee,
	})
	if err != nil {
		panic(err)
	}

	for _, it := range jo.Items {
		_, err = db.CreateItem(ctx, pg.CreateItemParams{
			OrderID:     orderID,
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
		})
		if err != nil {
			panic(err)
		}
	}

	if err = commit(ctx); err != nil {
		panic(err)
	}
	return nil
}

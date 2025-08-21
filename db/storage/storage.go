package storage

import (
	"context"
	"log/slog"

	"github.com/WhiCu/stgorders/db/cache"
	"github.com/WhiCu/stgorders/db/model"
	"github.com/WhiCu/stgorders/db/pg"
	"github.com/jackc/pgx/v5/pgxpool"
)

type event func(ctx context.Context) error

type Storage struct {
	*pg.Queries
	conn *pgxpool.Pool
	log  *slog.Logger
}

func NewStorage(conn *pgxpool.Pool, log *slog.Logger) *Storage {
	log.Debug("creating storage")
	db := pg.New(conn)
	return &Storage{
		Queries: db,
		conn:    conn,
		log:     log,
	}
}

func (s *Storage) WithTx(ctx context.Context) (storage *Storage, rollback event, commit event, err error) {
	s.log.Debug("starting transaction")
	tx, err := s.conn.Begin(ctx)
	return &Storage{Queries: s.Queries.WithTx(tx), conn: s.conn}, tx.Rollback, tx.Commit, err
}

func (s *Storage) Close() {
	s.conn.Close()
	s.log.Debug("closing connection", slog.String("CONN", s.conn.Config().ConnConfig.Database))
}

func (s *Storage) InitCache(ctx context.Context, cache cache.Cache[string, model.JsonOrder]) (err error) {
	defer func() {
		if err != nil {
			err = WrapErrPreloadCache(err)
		}
	}()
	log := s.log.WithGroup("init_cache")
	orders, err := s.GetLastOrders(ctx, int32(cache.Size()))
	if err != nil {
		log.Error("could not get last orders", slog.String("ERR", err.Error()))
		return err
	}

	for _, o := range orders {
		delivery, err := s.GetDeliveryByOrderUID(ctx, o.OrderUid)
		if err != nil {
			log.Error("could not get delivery", slog.String("ERR", err.Error()))
			return err
		}
		payment, err := s.GetPaymentByTransaction(ctx, o.OrderUid)
		if err != nil {
			log.Error("could not get payment", slog.String("ERR", err.Error()))
			return err
		}
		items, err := s.GetItemsByTrackNumber(ctx, o.TrackNumber)
		if err != nil {
			log.Error("could not get items", slog.String("ERR", err.Error()))
			return err
		}

		jsonOrder := model.ModelJSONOrder(o, delivery, payment, items)

		cache.Set(o.OrderUid, jsonOrder)
	}

	log.Info("cache preloaded", slog.Int("count", cache.Size()))
	return nil
}

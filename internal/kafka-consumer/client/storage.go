package client

import (
	"context"

	"github.com/WhiCu/stgorders/db/pg"
	"github.com/jackc/pgx/v5/pgxpool"
)

type event func(ctx context.Context) error

type Storage struct {
	*pg.Queries
	conn *pgxpool.Pool
}

func NewStorage(conn *pgxpool.Pool) *Storage {
	db := pg.New(conn)
	return &Storage{
		Queries: db,
		conn:    conn,
	}
}

func (s *Storage) WithTx(ctx context.Context) (storage *pg.Queries, rollback event, commit event, err error) {
	tx, err := s.conn.Begin(ctx)
	return s.Queries.WithTx(tx), tx.Rollback, tx.Commit, err

}

func (s *Storage) Close() {
	s.conn.Close()
}

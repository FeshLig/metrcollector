package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBChecker struct {
	pool *pgxpool.Pool
}

func NewDBChecker(db *pgxpool.Pool) *DBChecker {
	return &DBChecker{
		pool: db,
	}
}

func (p *DBChecker) Check(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres wraps PostgreSQL connection pool.
type Postgres struct {
	pool *pgxpool.Pool
}

// NewPostgres creates PostgreSQL connection pool and runs database migrations.
func NewPostgres(ctx context.Context, dsn string) (*Postgres, error) {

	if err := RunMigrations(dsn); err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("db init failed: %w", err)
	}

	return &Postgres{pool: pool}, nil
}

// Close closes PostgreSQL connection pool.
func (p *Postgres) Close() {
	p.pool.Close()
}

package repository

import (
	"context"
	"errors"
	"time"

	"github.com/FeshLig/metrcollector/internal/metric"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	gaugeStr   = "gauge"
	counterStr = "counter"
)

type PostgresStorage struct {
	db *pgxpool.Pool
}

func NewPostgresStorage(db *Postgres) *PostgresStorage {
	return &PostgresStorage{
		db: db.pool,
	}
}

func (p *PostgresStorage) SetGauge(ctx context.Context, name string, value metric.Gauge) error {
	return p.withTxRetry(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			`INSERT INTO metrics (name, type, gauge) 
			VALUES ($1, $2, $3) 
			ON CONFLICT (name, type) 
			DO UPDATE SET gauge = EXCLUDED.gauge`,
			name, gaugeStr, float64(value))

		return err
	})
}

func (p *PostgresStorage) AddCounter(ctx context.Context, name string, value metric.Counter) error {
	return p.withTxRetry(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			`INSERT INTO metrics (name, type, counter) 
			VALUES ($1, $2, $3) 
			ON CONFLICT (name, type) 
			DO UPDATE SET counter = metrics.counter + EXCLUDED.counter`,
			name, counterStr, int64(value))

		return err
	})
}

func (p *PostgresStorage) GetGauge(ctx context.Context, name string) (metric.Gauge, bool) {

	var gauge float64
	found := true

	err := p.withTxRetry(ctx, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `SELECT gauge FROM metrics WHERE type = $1 AND name = $2`, gaugeStr, name).Scan(&gauge)
		if errors.Is(err, pgx.ErrNoRows) {
			found = false
			return nil
		}
		return err
	})
	if err != nil || !found {
		return 0, false
	}
	return metric.Gauge(gauge), true
}

func (p *PostgresStorage) GetCounter(ctx context.Context, name string) (metric.Counter, bool) {

	var counter int64
	found := true

	err := p.withTxRetry(ctx, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `SELECT counter FROM metrics WHERE type = $1 AND name = $2`, counterStr, name).Scan(&counter)
		if errors.Is(err, pgx.ErrNoRows) {
			found = false
			return nil
		}
		return err
	})
	if err != nil || !found {
		return 0, false
	}
	return metric.Counter(counter), true
}

func (p *PostgresStorage) SetCounter(ctx context.Context, name string, value metric.Counter) error {
	return p.withTxRetry(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			`INSERT INTO metrics (name, type, counter) 
			VALUES ($1, $2, $3) 
			ON CONFLICT (name, type) 
			DO UPDATE SET counter = EXCLUDED.counter`,
			name, counterStr, int64(value))
		return err
	})
}

func (p *PostgresStorage) SnapshotGauges(ctx context.Context) map[string]metric.Gauge {

	var result map[string]metric.Gauge

	err := p.withTxRetry(ctx, func(tx pgx.Tx) error {

		gauges := make(map[string]metric.Gauge)
		q, err := tx.Query(ctx, `SELECT name, gauge FROM metrics WHERE type = $1`, gaugeStr)
		if err != nil {
			return err
		}
		defer q.Close()

		for q.Next() {
			var name string
			var value float64

			if err := q.Scan(&name, &value); err != nil {
				return err
			}
			gauges[name] = metric.Gauge(value)
		}

		if err := q.Err(); err != nil {
			return err
		}

		result = gauges

		return nil
	})

	if err != nil {
		return nil
	}

	return result

}

func (p *PostgresStorage) SnapshotCounters(ctx context.Context) map[string]metric.Counter {

	var result map[string]metric.Counter

	err := p.withTxRetry(ctx, func(tx pgx.Tx) error {
		counters := make(map[string]metric.Counter)

		q, err := tx.Query(ctx, `SELECT name, counter FROM metrics WHERE type = $1`, counterStr)
		if err != nil {
			return err
		}
		defer q.Close()

		for q.Next() {
			var name string
			var value int64

			if err := q.Scan(&name, &value); err != nil {
				return err
			}
			counters[name] = metric.Counter(value)
		}
		if err := q.Err(); err != nil {
			return err
		}
		result = counters
		return nil
	})

	if err != nil {
		return nil
	}
	return result

}

func (p *PostgresStorage) SetMetrics(ctx context.Context, gauges map[string]metric.Gauge, counters map[string]metric.Counter) error {

	return p.withTxRetry(ctx, func(tx pgx.Tx) error {

		batch := &pgx.Batch{}

		for name, value := range gauges {
			batch.Queue(`
                INSERT INTO metrics (name, type, gauge) 
				VALUES ($1, $2, $3) 
				ON CONFLICT (name, type) 
				DO UPDATE SET gauge = EXCLUDED.gauge
            `, name, gaugeStr, value)
		}

		for name, value := range counters {
			batch.Queue(`
				INSERT INTO metrics (name, type, counter) 
				VALUES ($1, $2, $3) 
				ON CONFLICT (name, type) 
				DO UPDATE SET counter = metrics.counter + EXCLUDED.counter
			`, name, counterStr, value)
		}

		br := tx.SendBatch(ctx, batch)

		for i := 0; i < batch.Len(); i++ {
			_, err := br.Exec()
			if err != nil {
				br.Close()
				return err
			}
		}

		if err := br.Close(); err != nil {
			return err
		}

		return nil
	})

}

func (p *PostgresStorage) SnapshotMetrics(ctx context.Context) (map[string]metric.Gauge, map[string]metric.Counter) {
	gauges := p.SnapshotGauges(ctx)
	counters := p.SnapshotCounters(ctx)

	return gauges, counters
}

func (p *PostgresStorage) withTxRetry(ctx context.Context, fn func(pgx.Tx) error) error {

	const maxRetries = 3
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {

		tx, err := p.db.Begin(ctx)
		if err != nil {
			if isRetriable(err) {
				if attempt == maxRetries {
					break
				}

				backoff := time.Duration(1+2*attempt) * time.Second

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(backoff):
				}
				lastErr = err
				continue
			}
			return err
		}

		err = fn(tx)
		if err != nil {
			tx.Rollback(ctx)

			if isRetriable(err) {
				if attempt == maxRetries {
					break
				}

				backoff := time.Duration(1+2*attempt) * time.Second

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(backoff):
				}
				lastErr = err
				continue
			}
			return err
		}

		err = tx.Commit(ctx)
		if err != nil {
			if isRetriable(err) {
				if attempt == maxRetries {
					break
				}

				backoff := time.Duration(1+2*attempt) * time.Second

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(backoff):
				}
				lastErr = err
				continue
			}
			return err
		}

		return nil
	}

	return lastErr
}

func isRetriable(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgerrcode.IsConnectionException(pgErr.Code)
	}

	return false
}

func (p *PostgresStorage) Check(ctx context.Context) error {
	return p.db.Ping(ctx)
}

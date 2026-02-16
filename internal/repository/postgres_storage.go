package repository

import (
	"context"
	"fmt"

	"github.com/FeshLig/metrcollector/internal/metric"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO: добавить всем методам ошибки на return

const (
	gaugeStr   = "gauge"
	counterStr = "counter"
)

type PostgresStorage struct {
	db *pgxpool.Pool
}

func NewPostgresStorage(db *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

func (p *PostgresStorage) SetGauge(ctx context.Context, name string, value metric.Gauge) error {
	_, err := p.db.Exec(ctx,
		`INSERT INTO metrics (name, type, gauge) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (name, type) 
		DO UPDATE SET gauge = EXCLUDED.gauge`,
		name, gaugeStr, float64(value))

	return err
}

func (p *PostgresStorage) AddCounter(ctx context.Context, name string, value metric.Counter) error {
	_, err := p.db.Exec(ctx,
		`INSERT INTO metrics (name, type, counter) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (name, type) 
		DO UPDATE SET counter = metrics.counter + EXCLUDED.counter`,
		name, counterStr, int64(value))

	return err
}

func (p *PostgresStorage) GetGauge(ctx context.Context, name string) (metric.Gauge, bool) {

	var gauge float64

	err := p.db.QueryRow(ctx, `SELECT gauge FROM metrics WHERE type = $1 AND name = $2`, gaugeStr, name).Scan(&gauge)
	if err == pgx.ErrNoRows {
		return 0, false
	}
	if err != nil {
		// fmt.Errorf("counter scan error: %w", err)
		return 0, false
	}

	return metric.Gauge(gauge), true
}

func (p *PostgresStorage) GetCounter(ctx context.Context, name string) (metric.Counter, bool) {

	var counter int64

	err := p.db.QueryRow(ctx, `SELECT counter FROM metrics WHERE type = $1 AND name = $2`, counterStr, name).Scan(&counter)
	if err == pgx.ErrNoRows {
		return 0, false
	}
	if err != nil {
		// fmt.Errorf("counter scan error: %w", err)
		return 0, false
	}

	return metric.Counter(counter), true
}

func (p *PostgresStorage) SetCounter(name string, value metric.Counter) {
	p.db.Exec(context.TODO(),
		`INSERT INTO metrics (name, type, counter) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (name, type) 
		DO UPDATE SET counter = EXCLUDED.counter`,
		name, counterStr, int64(value))

}

func (p *PostgresStorage) SnapshotGauges(ctx context.Context) map[string]metric.Gauge {
	gauges := make(map[string]metric.Gauge)

	q, err := p.db.Query(ctx, `SELECT name, gauge FROM metrics WHERE type = $1`, gaugeStr)
	if err != nil {
		return nil
	}
	defer q.Close()

	for q.Next() {
		var name string
		var value float64

		if err := q.Scan(&name, &value); err != nil {
			return nil
		}
		gauges[name] = metric.Gauge(value)
	}

	if err := q.Err(); err != nil {
		return nil
	}

	return gauges

}

func (p *PostgresStorage) SnapshotCounters(ctx context.Context) map[string]metric.Counter {
	counters := make(map[string]metric.Counter)

	q, err := p.db.Query(ctx, `SELECT name, counter FROM metrics WHERE type = $1`, counterStr)
	if err != nil {
		return nil
	}
	defer q.Close()

	for q.Next() {
		var name string
		var value int64

		if err := q.Scan(&name, &value); err != nil {
			return nil
		}
		counters[name] = metric.Counter(value)
	}
	if err := q.Err(); err != nil {
		return nil
	}

	return counters

}

func (p *PostgresStorage) SetMetrics(ctx context.Context, gauges map[string]metric.Gauge, counters map[string]metric.Counter) error {

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

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
	defer br.Close()

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

	if err = tx.Commit(ctx); err != nil {
		fmt.Printf("%s\n", err.Error())
		return err
	} else {
		fmt.Println("commit ok")
	}
	return nil

}

func (p *PostgresStorage) SnapshotMetrics(ctx context.Context) (map[string]metric.Gauge, map[string]metric.Counter) {
	gauges := p.SnapshotGauges(ctx)
	counters := p.SnapshotCounters(ctx)

	return gauges, counters
}

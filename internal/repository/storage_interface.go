package repository

import (
	"context"

	"github.com/FeshLig/metrcollector/internal/metric"
)

// Storage provides metric persistence operations.
type Storage interface {
	SetGauge(ctx context.Context, name string, value metric.Gauge) error
	AddCounter(ctx context.Context, name string, delta metric.Counter) error

	GetGauge(ctx context.Context, name string) (metric.Gauge, bool)
	GetCounter(ctx context.Context, name string) (metric.Counter, bool)

	SnapshotGauges(ctx context.Context) map[string]metric.Gauge
	SnapshotCounters(ctx context.Context) map[string]metric.Counter
	SnapshotMetrics(ctx context.Context) (map[string]metric.Gauge, map[string]metric.Counter)

	SetMetrics(ctx context.Context, gauges map[string]metric.Gauge, counters map[string]metric.Counter) error

	Check(ctx context.Context) error
}

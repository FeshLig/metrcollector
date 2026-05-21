package repository

import (
	"context"
	"sync"

	"github.com/FeshLig/metrcollector/internal/metric"
)

// MemStorage stores metrics in memory.
type MemStorage struct {
	mu       sync.Mutex
	gauges   map[string]metric.Gauge
	counters map[string]metric.Counter
}

// NewMemStorage creates new in-memory storage instance.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]metric.Gauge),
		counters: make(map[string]metric.Counter),
	}
}

// SetGauge stores gauge metric value.
func (m *MemStorage) SetGauge(ctx context.Context, name string, value metric.Gauge) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauges[name] = m.gauges[name].SetGauge(value)
	return nil
}

// GetGauge returns gauge metric value by name.
func (m *MemStorage) GetGauge(ctx context.Context, name string) (metric.Gauge, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	gauge, ok := m.gauges[name]

	return gauge, ok
}

// AddCounter increments counter metric value.
func (m *MemStorage) AddCounter(ctx context.Context, name string, delta metric.Counter) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] = m.counters[name].AddCounter(delta)
	return nil
}

// SetCounter sets counter metric value.
func (m *MemStorage) SetCounter(ctx context.Context, name string, value metric.Counter) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] = m.counters[name].SetCounter(value)
	return nil
}

// GetCounter returns counter metric value by name.
func (m *MemStorage) GetCounter(ctx context.Context, name string) (metric.Counter, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	counter, ok := m.counters[name]

	return counter, ok
}

// SetMetrics stores multiple metrics.
func (m *MemStorage) SetMetrics(ctx context.Context, gauges map[string]metric.Gauge, counters map[string]metric.Counter) error {
	for name, value := range gauges {
		m.SetGauge(ctx, name, value)
	}

	for name, value := range counters {
		m.AddCounter(ctx, name, value)
	}

	return nil
}

// SnapshotGauges returns copy of all gauge metrics.
func (m *MemStorage) SnapshotGauges(ctx context.Context) map[string]metric.Gauge {
	m.mu.Lock()
	defer m.mu.Unlock()
	copy := make(map[string]metric.Gauge, len(m.gauges))
	for k, v := range m.gauges {
		copy[k] = v
	}
	return copy
}

// SnapshotCounters returns copy of all counter metrics.
func (m *MemStorage) SnapshotCounters(ctx context.Context) map[string]metric.Counter {

	m.mu.Lock()
	defer m.mu.Unlock()
	copy := make(map[string]metric.Counter, len(m.counters))
	for k, v := range m.counters {
		copy[k] = v
	}
	return copy
}

// SnapshotMetrics returns copies of all stored metrics.
func (m *MemStorage) SnapshotMetrics(ctx context.Context) (map[string]metric.Gauge, map[string]metric.Counter) {
	gauges := m.SnapshotGauges(ctx)
	counters := m.SnapshotCounters(ctx)

	return gauges, counters
}

// Check verifies storage availability.
func (m *MemStorage) Check(ctx context.Context) error {
	return nil
}

package repository

import (
	"context"
	"sync"

	"github.com/FeshLig/metrcollector/internal/metric"
)

type MemStorage struct {
	mu       sync.Mutex
	gauges   map[string]metric.Gauge
	counters map[string]metric.Counter
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]metric.Gauge),
		counters: make(map[string]metric.Counter),
	}
}

func (m *MemStorage) SetGauge(ctx context.Context, name string, value metric.Gauge) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauges[name] = m.gauges[name].SetGauge(value)
	return nil
}

func (m *MemStorage) GetGauge(ctx context.Context, name string) (metric.Gauge, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	gauge, ok := m.gauges[name]

	return gauge, ok
}

func (m *MemStorage) AddCounter(ctx context.Context, name string, delta metric.Counter) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] = m.counters[name].AddCounter(delta)
	return nil
}

func (m *MemStorage) SetCounter(ctx context.Context, name string, value metric.Counter) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] = m.counters[name].SetCounter(value)
	return nil
}

func (m *MemStorage) GetCounter(ctx context.Context, name string) (metric.Counter, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	counter, ok := m.counters[name]

	return counter, ok
}

func (m *MemStorage) SetMetrics(ctx context.Context, gauges map[string]metric.Gauge, counters map[string]metric.Counter) error {
	for name, value := range gauges {
		m.SetGauge(ctx, name, value)
	}

	for name, value := range counters {
		m.AddCounter(ctx, name, value)
	}

	return nil
}

func (m *MemStorage) SnapshotGauges(ctx context.Context) map[string]metric.Gauge {
	m.mu.Lock()
	defer m.mu.Unlock()
	copy := make(map[string]metric.Gauge, len(m.gauges))
	for k, v := range m.gauges {
		copy[k] = v
	}
	return copy
}

func (m *MemStorage) SnapshotCounters(ctx context.Context) map[string]metric.Counter {

	m.mu.Lock()
	defer m.mu.Unlock()
	copy := make(map[string]metric.Counter, len(m.counters))
	for k, v := range m.counters {
		copy[k] = v
	}
	return copy
}

func (m *MemStorage) SnapshotMetrics(ctx context.Context) (map[string]metric.Gauge, map[string]metric.Counter) {
	gauges := m.SnapshotGauges(ctx)
	counters := m.SnapshotCounters(ctx)

	return gauges, counters
}

func (m *MemStorage) Check(ctx context.Context) error {
	return nil
}

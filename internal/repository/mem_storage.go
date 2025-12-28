package repository

import (
	"sync"

	"github.com/FeshLig/metrcollector/internal/metric"
)

type Storage interface {
	SnapshotGauges() map[string]metric.Gauge
	SnapshotCounters() map[string]metric.Counter
}

// возможно нужно добавить потокобезопасность, т.к. одновременно может происходить и копия и запись
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

func (m *MemStorage) SetGauge(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauges[name] = metric.Gauge(value)
}

func (m *MemStorage) AddCounter(name string, delta int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] += metric.Counter(delta)
}

func (m *MemStorage) SnapshotGauges() map[string]metric.Gauge {
	m.mu.Lock()
	defer m.mu.Unlock()
	copy := make(map[string]metric.Gauge, len(m.gauges))
	for k, v := range m.gauges {
		copy[k] = v
	}
	return copy
}

func (m *MemStorage) SnapshotCounters() map[string]metric.Counter {
	m.mu.Lock()
	defer m.mu.Unlock()
	copy := make(map[string]metric.Counter, len(m.counters))
	for k, v := range m.counters {
		copy[k] = v
	}
	return copy
}

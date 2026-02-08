package repository

import (
	"sync"

	"github.com/FeshLig/metrcollector/internal/metric"
)

type Storage interface {
	SetGauge(name string, value metric.Gauge)
	AddCounter(name string, delta metric.Counter)

	GetGauge(name string) (metric.Gauge, bool)
	GetCounter(name string) (metric.Counter, bool)

	SnapshotGauges() map[string]metric.Gauge
	SnapshotCounters() map[string]metric.Counter
}

type MemStorage struct {
	mu       sync.RWMutex
	gauges   map[string]metric.Gauge
	counters map[string]metric.Counter
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]metric.Gauge),
		counters: make(map[string]metric.Counter),
	}
}

func (m *MemStorage) SetGauge(name string, value metric.Gauge) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauges[name] = m.gauges[name].SetGauge(value)
}

func (m *MemStorage) GetGauge(name string) (metric.Gauge, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	gauge, ok := m.gauges[name]

	return gauge, ok
}

func (m *MemStorage) AddCounter(name string, delta metric.Counter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] = m.counters[name].AddCounter(delta)
}

func (m *MemStorage) SetCounter(name string, value metric.Counter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] = m.counters[name].SetCounter(value)
}

func (m *MemStorage) GetCounter(name string) (metric.Counter, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	counter, ok := m.counters[name]

	return counter, ok
}

func (m *MemStorage) SetMetrics(gauges map[string]metric.Gauge, counters map[string]metric.Counter) {
	for name, value := range gauges {
		m.SetGauge(name, value)
	}

	for name, value := range counters {
		m.SetCounter(name, value)
	}
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

func (m *MemStorage) SnapshotMetrics() (map[string]metric.Gauge, map[string]metric.Counter) {
	gauges := m.SnapshotGauges()
	counters := m.SnapshotCounters()

	return gauges, counters
}

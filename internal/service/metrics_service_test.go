package service_test

import (
	"context"
	"testing"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/metric"
	"github.com/FeshLig/metrcollector/internal/service"
)

type mockStorage struct {
	gauges   map[string]metric.Gauge
	counters map[string]metric.Counter

	setGaugeCalled   bool
	addCounterCalled bool
}

func newMockStorage() *mockStorage {
	return &mockStorage{
		gauges:   make(map[string]metric.Gauge),
		counters: make(map[string]metric.Counter),
	}
}

func (m *mockStorage) SetGauge(ctx context.Context, name string, value metric.Gauge) error {
	m.setGaugeCalled = true
	m.gauges[name] = value
	return nil
}

func (m *mockStorage) AddCounter(ctx context.Context, name string, value metric.Counter) error {
	m.addCounterCalled = true
	m.counters[name] += value
	return nil
}

func (m *mockStorage) GetGauge(ctx context.Context, name string) (metric.Gauge, bool) {
	v, ok := m.gauges[name]
	return v, ok
}

func (m *mockStorage) GetCounter(ctx context.Context, name string) (metric.Counter, bool) {
	v, ok := m.counters[name]
	return v, ok
}

func (m *mockStorage) SnapshotGauges(ctx context.Context) map[string]metric.Gauge {
	return m.gauges
}

func (m *mockStorage) SnapshotCounters(ctx context.Context) map[string]metric.Counter {
	return m.counters
}

func (m *mockStorage) SetMetrics(ctx context.Context, gauges map[string]metric.Gauge, counters map[string]metric.Counter) {
	for name, value := range gauges {
		m.gauges[name] = value
	}
	for name, value := range counters {
		m.counters[name] = value
	}
}

func (m *mockStorage) SnapshotMetrics(ctx context.Context) (map[string]metric.Gauge, map[string]metric.Counter) {
	gauges := m.SnapshotGauges(ctx)
	counters := m.SnapshotCounters(ctx)

	return gauges, counters
}

type mockPersister struct {
	saveCalled bool
}

func (m *mockPersister) SaveNow() {
	m.saveCalled = true
}

func TestMetricService_Update(t *testing.T) {
	value := 42.5
	delta := int64(10)

	tests := []struct {
		name        string
		metric      dto.Metrics
		syncSave    bool
		expectError bool
		errCode     service.ErrorCode
	}{
		{
			name: "valid gauge",
			metric: dto.Metrics{
				ID:    "g1",
				MType: dto.Gauge,
				Value: &value,
			},
		},
		{
			name: "gauge without value",
			metric: dto.Metrics{
				ID:    "g1",
				MType: dto.Gauge,
			},
			expectError: true,
			errCode:     service.ErrInvalidValue,
		},
		{
			name: "valid counter",
			metric: dto.Metrics{
				ID:    "c1",
				MType: dto.Counter,
				Delta: &delta,
			},
		},
		{
			name: "counter without delta",
			metric: dto.Metrics{
				ID:    "c1",
				MType: dto.Counter,
			},
			expectError: true,
			errCode:     service.ErrInvalidValue,
		},
		{
			name: "unknown type",
			metric: dto.Metrics{
				ID:    "x",
				MType: "histogram",
			},
			expectError: true,
			errCode:     service.ErrInvalidType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := newMockStorage()
			persister := &mockPersister{}

			svc := service.NewMetricService(storage, persister, tt.syncSave)

			err := svc.Update(tt.metric)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				svcErr, ok := err.(*service.ServiceError)
				if !ok {
					t.Fatalf("expected ServiceError, got %T", err)
				}
				if svcErr.Code != tt.errCode {
					t.Fatalf("expected error code %v, got %v", tt.errCode, svcErr.Code)
				}
				return
			}

			if tt.metric.MType == dto.Gauge && !storage.setGaugeCalled && !tt.expectError {
				t.Fatal("expected SetGauge to be called")
			}

			if tt.metric.MType == dto.Counter && !storage.addCounterCalled && !tt.expectError {
				t.Fatal("expected AddCounter to be called")
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestMetricService_Get(t *testing.T) {
	storage := newMockStorage()
	persister := &mockPersister{}

	v := metric.Gauge(12.3)
	c := metric.Counter(7)

	storage.gauges["g1"] = v
	storage.counters["c1"] = c

	svc := service.NewMetricService(storage, persister, false)

	t.Run("get gauge", func(t *testing.T) {
		res, err := svc.Get(dto.Metrics{
			ID:    "g1",
			MType: dto.Gauge,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Value == nil || *res.Value != float64(v) {
			t.Fatal("incorrect gauge value")
		}
	})

	t.Run("get counter", func(t *testing.T) {
		res, err := svc.Get(dto.Metrics{
			ID:    "c1",
			MType: dto.Counter,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Delta == nil || *res.Delta != int64(c) {
			t.Fatal("incorrect counter value")
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.Get(dto.Metrics{
			ID:    "missing",
			MType: dto.Gauge,
		})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestMetricService_SnapshotGaugeMetrics(t *testing.T) {
	storage := newMockStorage()
	persister := &mockPersister{}

	storage.gauges["g1"] = metric.Gauge(1.5)
	storage.gauges["g2"] = metric.Gauge(2.5)

	svc := service.NewMetricService(storage, persister, false)

	res := svc.SnapshotGaugeMetrics()

	if len(res) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(res))
	}

	found := map[string]float64{}
	for _, m := range res {
		if m.MType != dto.Gauge {
			t.Fatalf("expected gauge type, got %s", m.MType)
		}
		if m.Value == nil {
			t.Fatal("value is nil")
		}
		found[m.ID] = *m.Value
	}

	if found["g1"] != 1.5 || found["g2"] != 2.5 {
		t.Fatal("incorrect snapshot values")
	}
}

func TestMetricService_SnapshotCounterMetrics(t *testing.T) {
	storage := newMockStorage()
	persister := &mockPersister{}

	storage.counters["c1"] = metric.Counter(10)
	storage.counters["c2"] = metric.Counter(20)

	svc := service.NewMetricService(storage, persister, false)

	res := svc.SnapshotCounterMetrics()

	if len(res) != 2 {
		t.Fatalf("expected 2 metrics, got %d", len(res))
	}

	found := map[string]int64{}
	for _, m := range res {
		if m.MType != dto.Counter {
			t.Fatalf("expected counter type, got %s", m.MType)
		}
		if m.Delta == nil {
			t.Fatal("delta is nil")
		}
		found[m.ID] = *m.Delta
	}

	if found["c1"] != 10 || found["c2"] != 20 {
		t.Fatal("incorrect snapshot values")
	}
}

func TestMetricService_Update_SyncSave(t *testing.T) {
	value := 1.0

	storage := newMockStorage()
	persister := &mockPersister{}

	svc := service.NewMetricService(storage, persister, true)

	err := svc.Update(dto.Metrics{
		ID:    "g1",
		MType: dto.Gauge,
		Value: &value,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !persister.saveCalled {
		t.Fatal("expected SaveNow to be called")
	}
}

func TestMetricService_Update_NoSyncSave(t *testing.T) {
	value := 1.0

	storage := newMockStorage()
	persister := &mockPersister{}

	svc := service.NewMetricService(storage, persister, false)

	_ = svc.Update(dto.Metrics{
		ID:    "g1",
		MType: dto.Gauge,
		Value: &value,
	})

	if persister.saveCalled {
		t.Fatal("SaveNow should not be called")
	}
}

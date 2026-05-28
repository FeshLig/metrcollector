package service

import (
	"context"
	"fmt"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/metric"
	"github.com/FeshLig/metrcollector/internal/repository"
)

// MetricsService provides operations for working with metrics.
type MetricsService interface {
	Update(m dto.Metrics) error
	Updates(m []dto.Metrics) error
	Get(m dto.Metrics) (dto.Metrics, error)
	SnapshotGaugeMetrics() []dto.Metrics
	SnapshotCounterMetrics() []dto.Metrics
	Check(ctx context.Context) error
}

// filePrs provides synchronous metric persistence.
type filePrs interface {
	SaveNow()
}

// MetricServiceImpl implements MetricsService interface.
type MetricServiceImpl struct {
	storage   repository.Storage
	persister filePrs
	syncSave  bool
}

// NewMetricService creates new metric service instance.
func NewMetricService(
	s repository.Storage,
	f filePrs,
	syncSave bool,
) *MetricServiceImpl {
	return &MetricServiceImpl{
		storage:   s,
		persister: f,
		syncSave:  syncSave,
	}
}

// Update updates single metric value.
func (s *MetricServiceImpl) Update(m dto.Metrics) error {

	name := m.ID
	metricType := m.MType

	switch metricType {

	case dto.Gauge:
		if m.Value == nil {
			return &ServiceError{
				Code: ErrInvalidValue,
				Msg:  "empty gauge value",
			}
		}
		s.storage.SetGauge(context.TODO(), name, metric.Gauge(*m.Value))
		if s.syncSave && s.persister != nil {
			s.persister.SaveNow()
		}

	case dto.Counter:
		if m.Delta == nil {
			return &ServiceError{
				Code: ErrInvalidValue,
				Msg:  "empty counter delta",
			}
		}
		s.storage.AddCounter(context.TODO(), name, metric.Counter(*m.Delta))
		if s.syncSave && s.persister != nil {
			s.persister.SaveNow()
		}

	default:
		return &ServiceError{
			Code: ErrInvalidType,
			Msg:  fmt.Sprintf("unknown metric type: %s", metricType),
		}

	}

	return nil

}

// Updates updates multiple metrics atomically.
func (s *MetricServiceImpl) Updates(m []dto.Metrics) error {

	gauges := make(map[string]metric.Gauge, len(m))
	counters := make(map[string]metric.Counter, len(m))

	for _, mVal := range m {
		name := mVal.ID
		metricType := mVal.MType

		switch metricType {

		case dto.Gauge:
			if mVal.Value == nil {
				return &ServiceError{
					Code: ErrInvalidValue,
					Msg:  "empty gauge value",
				}
			}
			gauges[name] = metric.Gauge(*mVal.Value)

		case dto.Counter:
			if mVal.Delta == nil {
				return &ServiceError{
					Code: ErrInvalidValue,
					Msg:  "empty counter delta",
				}
			}
			counters[name] += metric.Counter(*mVal.Delta)

		default:
			return &ServiceError{
				Code: ErrInvalidType,
				Msg:  fmt.Sprintf("unknown metric type: %s", metricType),
			}

		}
	}

	if err := s.storage.SetMetrics(context.TODO(), gauges, counters); err != nil {
		return err
	}
	if s.syncSave && s.persister != nil {
		s.persister.SaveNow()
	}

	return nil

}

// Get returns metric value by name and type.
func (s *MetricServiceImpl) Get(m dto.Metrics) (dto.Metrics, error) {

	name := m.ID
	metricType := m.MType

	result := dto.Metrics{
		ID:    name,
		MType: metricType,
	}

	switch metricType {

	case dto.Gauge:
		value, ok := s.storage.GetGauge(context.TODO(), name)
		if !ok {
			return m, &ServiceError{
				Code: ErrNotFound,
				Msg:  fmt.Sprintf("gauge with name %s doesn't exist", name),
			}
		}
		v := float64(value)
		result.Value = &v

	case dto.Counter:
		value, ok := s.storage.GetCounter(context.TODO(), name)
		if !ok {
			return m, &ServiceError{
				Code: ErrNotFound,
				Msg:  fmt.Sprintf("counter with name %s doesn't exist", name),
			}
		}
		v := int64(value)
		result.Delta = &v

	default:
		return m, &ServiceError{
			Code: ErrInvalidType,
			Msg:  fmt.Sprintf("unknown metric type: %s", metricType),
		}

	}

	return result, nil

}

// SnapshotGaugeMetrics returns snapshot of all gauge metrics.
func (s *MetricServiceImpl) SnapshotGaugeMetrics() []dto.Metrics {

	gauges := s.storage.SnapshotGauges(context.TODO())
	metrics := make([]dto.Metrics, 0, len(gauges))

	for name, value := range gauges {
		v := float64(value)
		metrics = append(metrics, dto.Metrics{
			ID:    name,
			MType: dto.Gauge,
			Value: &v,
		})
	}

	return metrics

}

// SnapshotCounterMetrics returns snapshot of all counter metrics.
func (s *MetricServiceImpl) SnapshotCounterMetrics() []dto.Metrics {

	counters := s.storage.SnapshotCounters(context.TODO())
	metrics := make([]dto.Metrics, 0, len(counters))

	for name, value := range counters {
		v := int64(value)
		metrics = append(metrics, dto.Metrics{
			ID:    name,
			MType: dto.Counter,
			Delta: &v,
		})
	}

	return metrics

}

// Check verifies storage availability.
func (s *MetricServiceImpl) Check(ctx context.Context) error {
	return s.storage.Check(ctx)
}

package service

import (
	"fmt"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/metric"
	"github.com/FeshLig/metrcollector/internal/repository"
)

type MetricsService interface {
	Update(m dto.Metrics) error
	Get(m dto.Metrics) (dto.Metrics, error)
	SnapshotGaugeMetrics() []dto.Metrics
	SnapshotCounterMetrics() []dto.Metrics
}

type filePrs interface {
	SaveNow()
}

type MetricServiceImpl struct {
	storage   repository.Storage
	persister filePrs
	syncSave  bool
}

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
		s.storage.SetGauge(name, metric.Gauge(*m.Value))
		if s.syncSave {
			s.persister.SaveNow()
		}

	case dto.Counter:
		if m.Delta == nil {
			return &ServiceError{
				Code: ErrInvalidValue,
				Msg:  "empty counter delta",
			}
		}
		s.storage.AddCounter(name, metric.Counter(*m.Delta))
		if s.syncSave {
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

func (s *MetricServiceImpl) Get(m dto.Metrics) (dto.Metrics, error) {

	name := m.ID
	metricType := m.MType

	result := dto.Metrics{
		ID:    name,
		MType: metricType,
	}

	switch metricType {

	case dto.Gauge:
		value, ok := s.storage.GetGauge(name)
		if !ok {
			return m, &ServiceError{
				Code: ErrNotFound,
				Msg:  fmt.Sprintf("gauge with name %s doesn't exist", name),
			}
		}
		v := float64(value)
		result.Value = &v

	case dto.Counter:
		value, ok := s.storage.GetCounter(name)
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

func (s *MetricServiceImpl) SnapshotGaugeMetrics() []dto.Metrics {

	var metrics []dto.Metrics
	gauges := s.storage.SnapshotGauges()
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

func (s *MetricServiceImpl) SnapshotCounterMetrics() []dto.Metrics {

	var metrics []dto.Metrics
	counters := s.storage.SnapshotCounters()
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

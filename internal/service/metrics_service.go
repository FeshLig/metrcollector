package service

import (
	"fmt"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/metric"
)

type MetricsService interface {
	Update(m dto.Metrics) error
	Get(m dto.Metrics) (dto.Metrics, error)
}

type Storage interface {
	SetGauge(name string, value metric.Gauge)
	AddCounter(name string, delta metric.Counter)

	GetGauge(name string) (metric.Gauge, bool)
	GetCounter(name string) (metric.Counter, bool)

	// SnapshotGauges() map[string]metric.Gauge
	// SnapshotCounters() map[string]metric.Counter
}

type MetricServiceImpl struct {
	storage Storage
}

func NewMetricService(s Storage) *MetricServiceImpl {
	return &MetricServiceImpl{
		storage: s,
	}
}

func (s *MetricServiceImpl) Update(m dto.Metrics) error {

	name := m.ID
	metricType := m.MType

	switch metricType {

	case "gauge":
		if m.Value == nil {
			return &ServiceError{
				Code: ErrInvalidValue,
				Msg:  "empty gauge value",
			}
		}
		s.storage.SetGauge(name, metric.Gauge(*m.Value))

	case "counter":
		if m.Delta == nil {
			return &ServiceError{
				Code: ErrInvalidValue,
				Msg:  "empty counter delta",
			}
		}
		s.storage.AddCounter(name, metric.Counter(*m.Delta))

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

	case "gauge":
		value, ok := s.storage.GetGauge(name)
		if !ok {
			return m, &ServiceError{
				Code: ErrNotFound,
				Msg:  fmt.Sprintf("gauge with name %s doesn't exist", name),
			}
		}
		v := float64(value)
		result.Value = &v

	case "counter":
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

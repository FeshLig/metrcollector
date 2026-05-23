package service

import (
	"strconv"
	"testing"

	"github.com/FeshLig/metrcollector/internal/dto"
	"github.com/FeshLig/metrcollector/internal/repository"
)

func BenchmarkMetricService_UpdateGauge(b *testing.B) {
	storage := repository.NewMemStorage()

	service := NewMetricService(
		storage,
		nil,
		false,
	)

	value := 123.45

	metric := dto.Metrics{
		ID:    "Alloc",
		MType: dto.Gauge,
		Value: &value,
	}

	b.ReportAllocs()

	for b.Loop() {
		if err := service.Update(metric); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMetricService_UpdateCounter(b *testing.B) {
	storage := repository.NewMemStorage()

	service := NewMetricService(
		storage,
		nil,
		false,
	)

	delta := int64(1)

	metric := dto.Metrics{
		ID:    "PollCount",
		MType: dto.Counter,
		Delta: &delta,
	}

	b.ReportAllocs()

	for b.Loop() {
		if err := service.Update(metric); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMetricService_Updates(b *testing.B) {
	storage := repository.NewMemStorage()

	service := NewMetricService(
		storage,
		nil,
		false,
	)

	metrics := make([]dto.Metrics, 0, 100)

	for i := 0; i < 100; i++ {
		name := strconv.Itoa(i)

		if i%2 == 0 {
			value := float64(i)

			metrics = append(metrics, dto.Metrics{
				ID:    name,
				MType: dto.Gauge,
				Value: &value,
			})
		} else {
			delta := int64(i)

			metrics = append(metrics, dto.Metrics{
				ID:    name,
				MType: dto.Counter,
				Delta: &delta,
			})
		}
	}

	b.ReportAllocs()

	for b.Loop() {
		if err := service.Updates(metrics); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMetricService_GetGauge(b *testing.B) {
	storage := repository.NewMemStorage()

	service := NewMetricService(
		storage,
		nil,
		false,
	)

	value := 123.45

	err := service.Update(dto.Metrics{
		ID:    "Alloc",
		MType: dto.Gauge,
		Value: &value,
	})
	if err != nil {
		b.Fatal(err)
	}

	request := dto.Metrics{
		ID:    "Alloc",
		MType: dto.Gauge,
	}

	b.ReportAllocs()

	for b.Loop() {
		_, err := service.Get(request)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMetricService_GetCounter(b *testing.B) {
	storage := repository.NewMemStorage()

	service := NewMetricService(
		storage,
		nil,
		false,
	)

	delta := int64(100)

	err := service.Update(dto.Metrics{
		ID:    "PollCount",
		MType: dto.Counter,
		Delta: &delta,
	})
	if err != nil {
		b.Fatal(err)
	}

	request := dto.Metrics{
		ID:    "PollCount",
		MType: dto.Counter,
	}

	b.ReportAllocs()

	for b.Loop() {
		_, err := service.Get(request)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMetricService_SnapshotGaugeMetrics(b *testing.B) {
	storage := repository.NewMemStorage()

	service := NewMetricService(
		storage,
		nil,
		false,
	)

	for i := 0; i < 1000; i++ {
		value := float64(i)

		err := service.Update(dto.Metrics{
			ID:    strconv.Itoa(i),
			MType: dto.Gauge,
			Value: &value,
		})
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()

	for b.Loop() {
		_ = service.SnapshotGaugeMetrics()
	}
}

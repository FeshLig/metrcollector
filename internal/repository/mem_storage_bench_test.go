package repository

import (
	"context"
	"strconv"
	"testing"

	"github.com/FeshLig/metrcollector/internal/metric"
)

func BenchmarkMemStorage_SetGauge(b *testing.B) {
	storage := NewMemStorage()
	ctx := context.Background()

	b.ReportAllocs()

	i := 0
	for b.Loop() {
		name := strconv.Itoa(i)

		err := storage.SetGauge(
			ctx,
			name,
			metric.Gauge(i),
		)
		if err != nil {
			b.Fatal(err)
		}
		i++
	}
}

func BenchmarkMemStorage_AddCounter(b *testing.B) {
	storage := NewMemStorage()
	ctx := context.Background()

	b.ReportAllocs()

	i := 0
	for b.Loop() {
		name := strconv.Itoa(i)

		err := storage.AddCounter(
			ctx,
			name,
			metric.Counter(i),
		)
		if err != nil {
			b.Fatal(err)
		}
		i++
	}
}

func BenchmarkMemStorage_GetGauge(b *testing.B) {
	storage := NewMemStorage()
	ctx := context.Background()

	err := storage.SetGauge(
		ctx,
		"test",
		metric.Gauge(123),
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()

	for b.Loop() {
		_, ok := storage.GetGauge(ctx, "test")
		if !ok {
			b.Fatal("gauge not found")
		}
	}
}

func BenchmarkMemStorage_GetCounter(b *testing.B) {
	storage := NewMemStorage()
	ctx := context.Background()

	err := storage.SetCounter(
		ctx,
		"test",
		metric.Counter(123),
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()

	for b.Loop() {
		_, ok := storage.GetCounter(ctx, "test")
		if !ok {
			b.Fatal("counter not found")
		}
	}
}

func BenchmarkMemStorage_SetMetrics(b *testing.B) {
	storage := NewMemStorage()
	ctx := context.Background()

	gauges := make(map[string]metric.Gauge, 100)
	counters := make(map[string]metric.Counter, 100)

	for i := 0; i < 100; i++ {
		name := strconv.Itoa(i)

		gauges[name] = metric.Gauge(i)
		counters[name] = metric.Counter(i)
	}

	b.ReportAllocs()

	for b.Loop() {
		err := storage.SetMetrics(
			ctx,
			gauges,
			counters,
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMemStorage_ParallelAddCounter(b *testing.B) {
	storage := NewMemStorage()
	ctx := context.Background()

	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := storage.AddCounter(
				ctx,
				"test",
				1,
			)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

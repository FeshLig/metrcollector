package agent

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/FeshLig/metrcollector/internal/metric"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// MetricsCollector collects runtime and system metrics.
type MetricsCollector struct {
	storage repository.Storage
}

// NewMetricCollector creates new metrics collector instance.
func NewMetricCollector(memStorage repository.Storage) *MetricsCollector {
	return &MetricsCollector{
		storage: memStorage,
	}
}

// CollectMetrics collects Go runtime metrics and stores them.
func (c *MetricsCollector) CollectMetrics(ctx context.Context) {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memStorage := c.storage

	memStorage.SetGauge(ctx, "Alloc", metric.Gauge(m.Alloc))
	memStorage.SetGauge(ctx, "BuckHashSys", metric.Gauge(m.BuckHashSys))
	memStorage.SetGauge(ctx, "Frees", metric.Gauge(m.Frees))
	memStorage.SetGauge(ctx, "GCCPUFraction", metric.Gauge(m.GCCPUFraction))
	memStorage.SetGauge(ctx, "GCSys", metric.Gauge(m.GCSys))
	memStorage.SetGauge(ctx, "HeapAlloc", metric.Gauge(m.HeapAlloc))
	memStorage.SetGauge(ctx, "HeapIdle", metric.Gauge(m.HeapIdle))
	memStorage.SetGauge(ctx, "HeapInuse", metric.Gauge(m.HeapInuse))
	memStorage.SetGauge(ctx, "HeapObjects", metric.Gauge(m.HeapObjects))
	memStorage.SetGauge(ctx, "HeapReleased", metric.Gauge(m.HeapReleased))
	memStorage.SetGauge(ctx, "HeapSys", metric.Gauge(m.HeapSys))
	memStorage.SetGauge(ctx, "LastGC", metric.Gauge(m.LastGC))
	memStorage.SetGauge(ctx, "Lookups", metric.Gauge(m.Lookups))
	memStorage.SetGauge(ctx, "MCacheInuse", metric.Gauge(m.MCacheInuse))
	memStorage.SetGauge(ctx, "MCacheSys", metric.Gauge(m.MCacheSys))
	memStorage.SetGauge(ctx, "MSpanInuse", metric.Gauge(m.MSpanInuse))
	memStorage.SetGauge(ctx, "MSpanSys", metric.Gauge(m.MSpanSys))
	memStorage.SetGauge(ctx, "Mallocs", metric.Gauge(m.Mallocs))
	memStorage.SetGauge(ctx, "NextGC", metric.Gauge(m.NextGC))
	memStorage.SetGauge(ctx, "NumForcedGC", metric.Gauge(m.NumForcedGC))
	memStorage.SetGauge(ctx, "NumGC", metric.Gauge(m.NumGC))
	memStorage.SetGauge(ctx, "OtherSys", metric.Gauge(m.OtherSys))
	memStorage.SetGauge(ctx, "PauseTotalNs", metric.Gauge(m.PauseTotalNs))
	memStorage.SetGauge(ctx, "StackInuse", metric.Gauge(m.StackInuse))
	memStorage.SetGauge(ctx, "StackSys", metric.Gauge(m.StackSys))
	memStorage.SetGauge(ctx, "Sys", metric.Gauge(m.Sys))
	memStorage.SetGauge(ctx, "TotalAlloc", metric.Gauge(m.TotalAlloc))

	memStorage.SetGauge(ctx, "RandomValue", metric.Gauge(getRandomFloat()))

}

// CollectGopsutil collects system metrics using gopsutil package.
func (c *MetricsCollector) CollectGopsutil(ctx context.Context) error {

	m, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	c.storage.SetGauge(ctx, "TotalMemory", metric.Gauge(m.Total))
	c.storage.SetGauge(ctx, "FreeMemory", metric.Gauge(m.Free))

	percents, err := cpu.Percent(0, true)
	if err != nil {
		return err
	}

	for i, percent := range percents {
		name := fmt.Sprintf("CPUutilization%d", i+1)
		c.storage.SetGauge(ctx, name, metric.Gauge(percent))
	}

	return nil
}

func getRandomFloat() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Float64()
}

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
func (c *MetricsCollector) CollectMetrics() {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memStorage := c.storage

	memStorage.SetGauge(context.Background(), "Alloc", metric.Gauge(m.Alloc))
	memStorage.SetGauge(context.Background(), "BuckHashSys", metric.Gauge(m.BuckHashSys))
	memStorage.SetGauge(context.Background(), "Frees", metric.Gauge(m.Frees))
	memStorage.SetGauge(context.Background(), "GCCPUFraction", metric.Gauge(m.GCCPUFraction))
	memStorage.SetGauge(context.Background(), "GCSys", metric.Gauge(m.GCSys))
	memStorage.SetGauge(context.Background(), "HeapAlloc", metric.Gauge(m.HeapAlloc))
	memStorage.SetGauge(context.Background(), "HeapIdle", metric.Gauge(m.HeapIdle))
	memStorage.SetGauge(context.Background(), "HeapInuse", metric.Gauge(m.HeapInuse))
	memStorage.SetGauge(context.Background(), "HeapObjects", metric.Gauge(m.HeapObjects))
	memStorage.SetGauge(context.Background(), "HeapReleased", metric.Gauge(m.HeapReleased))
	memStorage.SetGauge(context.Background(), "HeapSys", metric.Gauge(m.HeapSys))
	memStorage.SetGauge(context.Background(), "LastGC", metric.Gauge(m.LastGC))
	memStorage.SetGauge(context.Background(), "Lookups", metric.Gauge(m.Lookups))
	memStorage.SetGauge(context.Background(), "MCacheInuse", metric.Gauge(m.MCacheInuse))
	memStorage.SetGauge(context.Background(), "MCacheSys", metric.Gauge(m.MCacheSys))
	memStorage.SetGauge(context.Background(), "MSpanInuse", metric.Gauge(m.MSpanInuse))
	memStorage.SetGauge(context.Background(), "MSpanSys", metric.Gauge(m.MSpanSys))
	memStorage.SetGauge(context.Background(), "Mallocs", metric.Gauge(m.Mallocs))
	memStorage.SetGauge(context.Background(), "NextGC", metric.Gauge(m.NextGC))
	memStorage.SetGauge(context.Background(), "NumForcedGC", metric.Gauge(m.NumForcedGC))
	memStorage.SetGauge(context.Background(), "NumGC", metric.Gauge(m.NumGC))
	memStorage.SetGauge(context.Background(), "OtherSys", metric.Gauge(m.OtherSys))
	memStorage.SetGauge(context.Background(), "PauseTotalNs", metric.Gauge(m.PauseTotalNs))
	memStorage.SetGauge(context.Background(), "StackInuse", metric.Gauge(m.StackInuse))
	memStorage.SetGauge(context.Background(), "StackSys", metric.Gauge(m.StackSys))
	memStorage.SetGauge(context.Background(), "Sys", metric.Gauge(m.Sys))
	memStorage.SetGauge(context.Background(), "TotalAlloc", metric.Gauge(m.TotalAlloc))

	memStorage.SetGauge(context.Background(), "RandomValue", metric.Gauge(getRandomFloat()))

}

// CollectGopsutil collects system metrics using gopsutil package.
func (c *MetricsCollector) CollectGopsutil() error {

	m, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	c.storage.SetGauge(context.Background(), "TotalMemory", metric.Gauge(m.Total))
	c.storage.SetGauge(context.Background(), "FreeMemory", metric.Gauge(m.Free))

	percents, err := cpu.Percent(0, true)
	if err != nil {
		return err
	}

	for i, percent := range percents {
		name := fmt.Sprintf("CPUutilization%d", i+1)
		c.storage.SetGauge(context.Background(), name, metric.Gauge(percent))
	}

	return nil
}

func getRandomFloat() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Float64()
}

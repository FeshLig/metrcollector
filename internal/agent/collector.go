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

type MetricsCollector struct {
	storage repository.Storage
}

func NewMetricCollector(memStorage repository.Storage) *MetricsCollector {
	return &MetricsCollector{
		storage: memStorage,
	}
}

func (c *MetricsCollector) CollectMetrics() {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memStorage := c.storage

	memStorage.SetGauge(context.TODO(), "Alloc", metric.Gauge(m.Alloc))
	memStorage.SetGauge(context.TODO(), "BuckHashSys", metric.Gauge(m.BuckHashSys))
	memStorage.SetGauge(context.TODO(), "Frees", metric.Gauge(m.Frees))
	memStorage.SetGauge(context.TODO(), "GCCPUFraction", metric.Gauge(m.GCCPUFraction))
	memStorage.SetGauge(context.TODO(), "GCSys", metric.Gauge(m.GCSys))
	memStorage.SetGauge(context.TODO(), "HeapAlloc", metric.Gauge(m.HeapAlloc))
	memStorage.SetGauge(context.TODO(), "HeapIdle", metric.Gauge(m.HeapIdle))
	memStorage.SetGauge(context.TODO(), "HeapInuse", metric.Gauge(m.HeapInuse))
	memStorage.SetGauge(context.TODO(), "HeapObjects", metric.Gauge(m.HeapObjects))
	memStorage.SetGauge(context.TODO(), "HeapReleased", metric.Gauge(m.HeapReleased))
	memStorage.SetGauge(context.TODO(), "HeapSys", metric.Gauge(m.HeapSys))
	memStorage.SetGauge(context.TODO(), "LastGC", metric.Gauge(m.LastGC))
	memStorage.SetGauge(context.TODO(), "Lookups", metric.Gauge(m.Lookups))
	memStorage.SetGauge(context.TODO(), "MCacheInuse", metric.Gauge(m.MCacheInuse))
	memStorage.SetGauge(context.TODO(), "MCacheSys", metric.Gauge(m.MCacheSys))
	memStorage.SetGauge(context.TODO(), "MSpanInuse", metric.Gauge(m.MSpanInuse))
	memStorage.SetGauge(context.TODO(), "MSpanSys", metric.Gauge(m.MSpanSys))
	memStorage.SetGauge(context.TODO(), "Mallocs", metric.Gauge(m.Mallocs))
	memStorage.SetGauge(context.TODO(), "NextGC", metric.Gauge(m.NextGC))
	memStorage.SetGauge(context.TODO(), "NumForcedGC", metric.Gauge(m.NumForcedGC))
	memStorage.SetGauge(context.TODO(), "NumGC", metric.Gauge(m.NumGC))
	memStorage.SetGauge(context.TODO(), "OtherSys", metric.Gauge(m.OtherSys))
	memStorage.SetGauge(context.TODO(), "PauseTotalNs", metric.Gauge(m.PauseTotalNs))
	memStorage.SetGauge(context.TODO(), "StackInuse", metric.Gauge(m.StackInuse))
	memStorage.SetGauge(context.TODO(), "StackSys", metric.Gauge(m.StackSys))
	memStorage.SetGauge(context.TODO(), "Sys", metric.Gauge(m.Sys))
	memStorage.SetGauge(context.TODO(), "TotalAlloc", metric.Gauge(m.TotalAlloc))

	memStorage.SetGauge(context.TODO(), "RandomValue", metric.Gauge(getRandomFloat()))

	memStorage.AddCounter(context.TODO(), "PollCount", 1)
}

func (c *MetricsCollector) CollectGopsutil() error {

	m, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	c.storage.SetGauge(context.TODO(), "TotalMemory", metric.Gauge(m.Total))
	c.storage.SetGauge(context.TODO(), "FreeMemory", metric.Gauge(m.Free))

	percents, err := cpu.Percent(0, true)
	if err != nil {
		return err
	}

	for i, percent := range percents {
		name := fmt.Sprintf("CPUutilization%d", i+1)
		c.storage.SetGauge(context.TODO(), name, metric.Gauge(percent))
	}

	return nil
}

func getRandomFloat() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Float64()
}

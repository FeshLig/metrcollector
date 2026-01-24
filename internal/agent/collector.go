package agent

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/FeshLig/metrcollector/internal/metric"
)

type UpdateMetrics interface {
	SetGauge(name string, value metric.Gauge)
	AddCounter(name string, delta metric.Counter)
}

type MetricsCollector struct {
	storage UpdateMetrics
}

func NewMetricCollector(memStorage UpdateMetrics) *MetricsCollector {
	return &MetricsCollector{
		storage: memStorage,
	}
}

func (c *MetricsCollector) CollectMetrics() {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memStorage := c.storage

	memStorage.SetGauge("Alloc", metric.Gauge(m.Alloc))
	memStorage.SetGauge("BuckHashSys", metric.Gauge(m.BuckHashSys))
	memStorage.SetGauge("Frees", metric.Gauge(m.Frees))
	memStorage.SetGauge("GCCPUFraction", metric.Gauge(m.GCCPUFraction))
	memStorage.SetGauge("GCSys", metric.Gauge(m.GCSys))
	memStorage.SetGauge("HeapAlloc", metric.Gauge(m.HeapAlloc))
	memStorage.SetGauge("HeapIdle", metric.Gauge(m.HeapIdle))
	memStorage.SetGauge("HeapInuse", metric.Gauge(m.HeapInuse))
	memStorage.SetGauge("HeapObjects", metric.Gauge(m.HeapObjects))
	memStorage.SetGauge("HeapReleased", metric.Gauge(m.HeapReleased))
	memStorage.SetGauge("HeapSys", metric.Gauge(m.HeapSys))
	memStorage.SetGauge("LastGC", metric.Gauge(m.LastGC))
	memStorage.SetGauge("Lookups", metric.Gauge(m.Lookups))
	memStorage.SetGauge("MCacheInuse", metric.Gauge(m.MCacheInuse))
	memStorage.SetGauge("MCacheSys", metric.Gauge(m.MCacheSys))
	memStorage.SetGauge("MSpanInuse", metric.Gauge(m.MSpanInuse))
	memStorage.SetGauge("MSpanSys", metric.Gauge(m.MSpanSys))
	memStorage.SetGauge("Mallocs", metric.Gauge(m.Mallocs))
	memStorage.SetGauge("NextGC", metric.Gauge(m.NextGC))
	memStorage.SetGauge("NumForcedGC", metric.Gauge(m.NumForcedGC))
	memStorage.SetGauge("NumGC", metric.Gauge(m.NumGC))
	memStorage.SetGauge("OtherSys", metric.Gauge(m.OtherSys))
	memStorage.SetGauge("PauseTotalNs", metric.Gauge(m.PauseTotalNs))
	memStorage.SetGauge("StackInuse", metric.Gauge(m.StackInuse))
	memStorage.SetGauge("StackSys", metric.Gauge(m.StackSys))
	memStorage.SetGauge("Sys", metric.Gauge(m.Sys))
	memStorage.SetGauge("TotalAlloc", metric.Gauge(m.TotalAlloc))

	memStorage.SetGauge("RandomValue", metric.Gauge(getRandomFloat()))

	memStorage.AddCounter("PollCount", 1)
}

func getRandomFloat() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Float64()
}

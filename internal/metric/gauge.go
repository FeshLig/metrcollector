package metric

type Gauge float64

// type Gauges map[string]Gauge

func (s Gauge) SetGauge(value float64) Gauge {
	s = Gauge(value)
	return s
}

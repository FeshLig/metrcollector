package metric

type Gauge float64

type Gauges map[string]Gauge

func (s Gauges) SetGauge(name string, value float64) {
	s[name] = Gauge(value)
}

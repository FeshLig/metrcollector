package metric

type Gauge float64

func (s Gauge) SetGauge(value float64) Gauge {
	s = Gauge(value)
	return s
}

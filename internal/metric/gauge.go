package metric

type Gauge float64

func (s Gauge) SetGauge(value Gauge) Gauge {
	s = value
	return s
}
